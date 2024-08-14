package cli

import (
	"context"
	"sync"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm"
	bridgeprocessor "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	coreconnector "github.com/hyle-team/bridgeless-signer/internal/connectors/core"
	"github.com/hyle-team/bridgeless-signer/internal/core"
	rabbitproducer "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/producer"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/pkg/errors"
)

func RunService(ctx context.Context, cfg config.Config) error {
	var (
		wg            = sync.WaitGroup{}
		serviceSigner = cfg.Signer()
		coreCfg       = cfg.CoreConnectorConfig()
		coreConnector = coreconnector.NewConnector(coreCfg.Connection, coreCfg.Settings)
		rabbitCfg     = cfg.RabbitMQConfig()
	)

	proxiesRepo, err := evm.NewProxiesRepository(cfg.Chains(), serviceSigner.Address())
	if err != nil {
		return errors.Wrap(err, "failed to create proxiesRepo repository")
	}

	producer, err := rabbitproducer.New(rabbitCfg.NewChannel(), rabbitCfg.ResendParams)
	if err != nil {
		return errors.Wrap(err, "failed to create publisher")
	}

	processor := bridgeprocessor.New(proxiesRepo, pg.NewDepositsQ(cfg.DB()), serviceSigner, coreConnector, coreConnector)

	core.RunServer(ctx, &wg, cfg, proxiesRepo, producer)
	core.RunConsumers(ctx, &wg, cfg, producer, processor)
	wg.Wait()

	return ctx.Err()
}
