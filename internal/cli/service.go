package cli

import (
	"context"
	coreConnector "github.com/hyle-team/bridgeless-signer/internal/bridge/core"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy"
	"os/signal"
	"sync"
	"syscall"

	bridgeProcessor "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/core"
	rabbitProducer "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/producer"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/pkg/errors"
)

func RunService(cfg config.Config) error {
	var (
		ctx, _    = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		wg        = sync.WaitGroup{}
		coreCfg   = cfg.CoreConnectorConfig()
		coreConn  = coreConnector.NewConnector(coreCfg.Connection, coreCfg.Settings)
		rabbitCfg = cfg.RabbitMQConfig()
	)

	proxiesRepo, err := proxy.NewProxiesRepository(cfg.Chains(), cfg.Log())
	if err != nil {
		return errors.Wrap(err, "failed to create proxiesRepo repository")
	}

	producer, err := rabbitProducer.New(rabbitCfg.NewChannel(), rabbitCfg.ResendParams)
	if err != nil {
		return errors.Wrap(err, "failed to create producer")
	}

	processor := bridgeProcessor.New(proxiesRepo, pg.NewDepositsQ(cfg.DB()), cfg.Signer(), coreConn)

	core.RunServer(ctx, &wg, cfg, proxiesRepo, producer)
	core.RunConsumers(ctx, &wg, cfg, producer, processor)

	wg.Wait()

	return ctx.Err()
}
