package cli

import (
	"context"
	coreConnector "github.com/hyle-team/bridgeless-signer/internal/bridge/core"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy"
	"sync"

	bridgeProcessor "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/core"
	rabbitPublisher "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/publisher"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/pkg/errors"
)

func RunService(ctx context.Context, cfg config.Config) error {
	var (
		wg        = sync.WaitGroup{}
		coreCfg   = cfg.CoreConnectorConfig()
		coreConn  = coreConnector.NewConnector(coreCfg.Connection, coreCfg.Settings)
		rabbitCfg = cfg.RabbitMQConfig()
	)

	proxiesRepo, err := proxy.NewProxiesRepository(cfg.Chains(), cfg.Log())
	if err != nil {
		return errors.Wrap(err, "failed to create proxiesRepo repository")
	}

	publisher, err := rabbitPublisher.New(rabbitCfg.NewChannel(), rabbitCfg.ResendParams)
	if err != nil {
		return errors.Wrap(err, "failed to create publisher")
	}

	processor := bridgeProcessor.New(proxiesRepo, pg.NewDepositsQ(cfg.DB()), cfg.Signer(), coreConn)

	core.RunServer(ctx, &wg, cfg, proxiesRepo, publisher)
	core.RunConsumers(ctx, &wg, cfg, publisher, processor)

	wg.Wait()

	return ctx.Err()
}
