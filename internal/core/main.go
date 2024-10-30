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

// RunConsumers runs consumers for all queues.
func RunConsumers(
	ctx context.Context,
	wg *sync.WaitGroup,
	cfg config.Config,
	producer rabbitTypes.Producer,
	processor *bridgeProcessor.Processor,
) {
	var (
		logger    = cfg.Log()
		rabbitCfg = cfg.RabbitMQConfig()
	)

	for i := uint(0); i < rabbitCfg.BaseConsumerInstances; i++ {
		idx := i + 1

		// initializing new instances per loop
		baseConsumers := map[string]rabbitTypes.Consumer{
			rabbitTypes.GetDepositQueue: consumer.NewBase[bridgeProcessor.GetDepositRequest](
				rabbitCfg.NewChannel(),
				consumer.GetName(consumer.GetDepositConsumerPrefix, idx),
				logger.
					WithField(serviceComponent, componentConsumer).
					WithField(componentPart, consumer.GetName(consumer.GetDepositConsumerPrefix, idx)),
				consumerProcessors.NewGetDepositHandler(processor, producer),
				producer,
			),
			rabbitTypes.EthSignWithdrawalQueue: consumer.NewBase[bridgeProcessor.WithdrawalRequest](
				rabbitCfg.NewChannel(),
				consumer.GetName(consumer.EthSignWithdrawalConsumerPrefix, idx),
				logger.
					WithField(serviceComponent, componentConsumer).
					WithField(componentPart, consumer.GetName(consumer.EthSignWithdrawalConsumerPrefix, idx)),
				consumerProcessors.NewEthereumSignWithdrawalHandler(processor, producer),
				producer,
			),
			rabbitTypes.ZanoSignWithdrawalQueue: consumer.NewBase[bridgeProcessor.WithdrawalRequest](
				rabbitCfg.NewChannel(),
				consumer.GetName(consumer.ZanoSignWithdrawalConsumerPrefix, idx),
				logger.
					WithField(serviceComponent, componentConsumer).
					WithField(componentPart, consumer.GetName(consumer.ZanoSignWithdrawalConsumerPrefix, idx)),
				consumerProcessors.NewZanoSignWithdrawalHandler(processor, producer),
				producer,
			),
			rabbitTypes.ZanoSendWithdrawalQueue: consumer.NewBase[bridgeProcessor.ZanoSignedWithdrawalRequest](
				rabbitCfg.NewChannel(),
				consumer.GetName(consumer.ZanoSendWithdrawalConsumerPrefix, idx),
				logger.
					WithField(serviceComponent, componentConsumer).
					WithField(componentPart, consumer.GetName(consumer.ZanoSendWithdrawalConsumerPrefix, idx)),
				consumerProcessors.NewZanoSendWithdrawalHandler(processor, producer),
				producer,
			),
		}

		wg.Add(len(baseConsumers))

		for queue, cns := range baseConsumers {
			go func(cns rabbitTypes.Consumer, queue string) {
				defer wg.Done()

				if err := cns.Consume(ctx, queue); err != nil {
					logger.WithError(err).Error(fmt.Sprintf("failed to consume for %s", cns.Name()))
				}

			}(cns, queue)
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		cns := consumer.NewBatch[bridgeProcessor.SubmitTransactionRequest](
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
		cns := consumer.NewBatch[bridgeProcessor.WithdrawalRequest](
			rabbitCfg.NewChannel(),
			consumer.BitcoinSendWithdrawalConsumerPrefix,
			logger.
				WithField(serviceComponent, componentConsumer).
				WithField(componentPart, consumer.BitcoinSendWithdrawalConsumerPrefix),
			consumerProcessors.NewBitcoinSendWithdrawalHandler(processor, producer),
			producer,
			rabbitCfg.BitcoinSubmitterOpts,
		)
		if err := cns.Consume(ctx, rabbitTypes.BtcSendWithdrawalQueue); err != nil {
			logger.WithError(err).Error(fmt.Sprintf("failed to consume for %s", consumer.BitcoinSendWithdrawalConsumerPrefix))
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
