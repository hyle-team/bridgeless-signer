package core

import (
	"context"
	"fmt"
	bridgeProcessor "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/server"
	"github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/consumer"
	consumerProcessors "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/consumer/processors"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"sync"
)

const (
	serviceComponent = "component"
	componentPart    = "part"

	componentServer   = "server"
	componentConsumer = "consumer"
)

type baseConsumer struct {
	deliveryProcessor rabbitTypes.DeliveryProcessor
	prefix            string
}

// RunConsumers runs consumers for all queues.
func RunConsumers(
	ctx context.Context,
	wg *sync.WaitGroup,
	cfg config.Config,
	producer rabbitTypes.Producer,
	processor *bridgeProcessor.Processor,
) {
	var (
		logger       = cfg.Log()
		rabbitCfg    = cfg.RabbitMQConfig()
		consumersMap = map[string]baseConsumer{
			rabbitTypes.GetDepositQueue: {
				deliveryProcessor: consumerProcessors.NewGetDepositHandler(processor, producer),
				prefix:            consumer.GetDepositConsumerPrefix,
			},
			rabbitTypes.FormWithdrawalQueue: {
				deliveryProcessor: consumerProcessors.NewFormWithdrawalHandler(processor, producer),
				prefix:            consumer.FormWithdrawalConsumerPrefix,
			},
			rabbitTypes.SignWithdrawalQueue: {
				deliveryProcessor: consumerProcessors.NewSignWithdrawalHandler(processor, producer),
				prefix:            consumer.SignWithdrawalConsumerPrefix,
			},
			rabbitTypes.SubmitWithdrawalQueue: {
				deliveryProcessor: consumerProcessors.NewSubmitWithdrawalHandler(processor, producer),
				prefix:            consumer.SubmitWithdrawalConsumerPrefix,
			},
		}
	)

	for i := uint(0); i < rabbitCfg.ConsumerInstances; i++ {
		idx := i + 1
		for queue, consumerCfg := range consumersMap {
			wg.Add(1)
			go func(id uint, queue string, consumerCfg baseConsumer) {
				defer wg.Done()

				consumerName := consumer.GetName(consumerCfg.prefix, id)
				cns := consumer.NewBase(
					rabbitCfg.NewChannel(),
					consumerName,
					logger.
						WithField(serviceComponent, componentConsumer).
						WithField(componentPart, consumerName),
					consumerCfg.deliveryProcessor,
					producer,
				)

				if err := cns.Consume(ctx, queue); err != nil {
					logger.WithError(err).Error(fmt.Sprintf("failed to consume for %s", consumerName))
				}
			}(idx, queue, consumerCfg)
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		cns := consumer.NewBatch[bridgeTypes.SubmitTransactionRequest](
			rabbitCfg.NewChannel(),
			consumer.SubmitTransactionConsumerPrefix,
			logger.
				WithField(serviceComponent, componentConsumer).
				WithField(componentPart, consumer.SubmitTransactionConsumerPrefix),
			consumerProcessors.NewSubmitTransactionHandler(processor),
			producer,
			rabbitCfg.TxSubmitterOpts,
		)
		if err := cns.Consume(ctx, rabbitTypes.SubmitTransactionQueue); err != nil {
			logger.WithError(err).Error(fmt.Sprintf("failed to consume for %s", consumer.SubmitTransactionConsumerPrefix))
		}
	}()
	go func() {
		defer wg.Done()
		cns := consumer.NewBatch[bridgeTypes.BitcoinWithdrawalRequest](
			rabbitCfg.NewChannel(),
			consumer.SubmitBitcoinWithdrawalConsumerPrefix,
			logger.
				WithField(serviceComponent, componentConsumer).
				WithField(componentPart, consumer.SubmitBitcoinWithdrawalConsumerPrefix),
			consumerProcessors.NewSubmitBitcoinWithdrawalHandler(processor, producer),
			producer,
			rabbitCfg.BitcoinSubmitterOpts,
		)
		if err := cns.Consume(ctx, rabbitTypes.SubmitBitcoinWithdrawalQueue); err != nil {
			logger.WithError(err).Error(fmt.Sprintf("failed to consume for %s", consumer.SubmitBitcoinWithdrawalConsumerPrefix))
		}
	}()

}

func RunServer(
	ctx context.Context,
	wg *sync.WaitGroup,
	cfg config.Config,
	proxiesRepo bridgeTypes.ProxiesRepository,
	producer rabbitTypes.Producer,
) {
	logger := cfg.Log()
	srv := server.NewServer(
		cfg.GRPCListener(),
		cfg.HTTPListener(),
		pg.NewDepositsQ(cfg.DB()),
		proxiesRepo,
		producer,
		logger.WithField(serviceComponent, componentServer),
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := srv.RunHTTP(ctx); err != nil {
			logger.WithError(err).Error("rest gateway error occurred")
		}
	}()

	go func() {
		defer wg.Done()
		if err := srv.RunGRPC(ctx); err != nil {
			logger.WithError(err).Error("grpc server error occurred")
		}
	}()
}
