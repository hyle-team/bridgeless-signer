package cli

import (
	"context"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm"
	bridgeProcessor "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/core"
	rabbitProducer "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/producer"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/pkg/errors"
)

func RunService(ctx context.Context, cfg config.Config) error {
	var (
		serviceSigner = cfg.Signer()
		rabbitCfg     = cfg.RabbitMQConfig()
	)

	proxiesRepo, err := evm.NewProxiesRepository(cfg.Chains(), serviceSigner.Address())
	if err != nil {
		return errors.Wrap(err, "failed to create proxiesRepo repository")
	}

	producer, err := rabbitProducer.New(rabbitCfg.NewChannel(), rabbitCfg.ResendParams)
	if err != nil {
		return errors.Wrap(err, "failed to create publisher")
	}

	processor := bridgeProcessor.New(proxiesRepo, pg.NewDepositsQ(cfg.DB()), serviceSigner, cfg.TokenPairer())

	go core.RunServer(ctx, cfg, proxiesRepo, producer)
	go core.RunConsumers(ctx, cfg, producer, processor)

	<-ctx.Done()
	return ctx.Err()
}
