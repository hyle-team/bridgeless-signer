package core

import (
	"context"
	"fmt"

	bridgeProcessor "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/core/api"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/handler"
	"github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/consumer"
	consumerProcessors "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/consumer/processors"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/pkg/errors"
)

type consumerConfig struct {
	deliveryProcessor rabbitTypes.DeliveryProcessor
	prefix            string
}

// RunConsumers runs consumers for all queues.
func RunConsumers(
	ctx context.Context,
	cfg config.Config,
	producer rabbitTypes.Producer,
	processor *bridgeProcessor.Processor,
) {
	var (
		logger       = cfg.Log()
		rabbitCfg    = cfg.RabbitMQConfig()
		consumersMap = map[string]consumerConfig{
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
		go func(index uint) {
			idx := index + 1

			for queue, consumerCfg := range consumersMap {
				go func(queue string, consumerCfg consumerConfig) {
					consumerName := consumer.GetName(consumerCfg.prefix, idx)

					logger.Info(fmt.Sprintf("starting consumer %s...", consumerName))
					err := consumer.New(
						rabbitCfg.NewChannel(),
						consumerName,
						logger.WithField("consumer", consumerName),
						consumerCfg.deliveryProcessor,
						producer,
					).Consume(ctx, queue)

					if err != nil {
						panic(errors.Wrap(err, fmt.Sprintf("failed to consume for %s", consumerName)))
					}
				}(queue, consumerCfg)
			}

		}(i)
	}
}

func RunServer(
	ctx context.Context,
	cfg config.Config,
	proxiesRepo bridgeTypes.ProxiesRepository,
	producer rabbitTypes.Producer,
) {
	logger := cfg.Log()

	server := api.NewServer(
		cfg.Listener(),
		cfg.RESTGatewayConfig(),
		handler.NewServiceHandler(
			pg.NewDepositsQ(cfg.DB()),
			logger.WithField("component", "grpc-handler"),
			proxiesRepo,
			producer,
		),
	)

	go func() {
		logger.Info("starting rest gateway...")
		if err := server.RunRESTGateway(ctx); err != nil {
			panic(errors.Wrap(err, "rest gateway error occurred"))
		}
	}()

	go func() {
		logger.Info("starting grpc server...")
		if err := server.RunGRPC(ctx); err != nil {
			panic(errors.Wrap(err, "grpc server error occurred"))
		}
	}()
}
