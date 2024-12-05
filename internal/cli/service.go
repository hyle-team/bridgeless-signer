package cli

import (
	"context"
	coreConnector "github.com/hyle-team/bridgeless-signer/internal/bridge/core"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
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
		wg        = sync.WaitGroup{}
		coreCfg   = cfg.CoreConnectorConfig()
		coreConn  = coreConnector.NewConnector(coreCfg.Connection, coreCfg.Settings)
		rabbitCfg = cfg.RabbitMQConfig()

		rabbitConnChan = rabbitCfg.Connection.NotifyClose(make(chan *amqp.Error, 1))
		ctx            = appContext(rabbitConnChan)
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

	return context.Cause(ctx)
}

func appContext(rabbit <-chan *amqp.Error) context.Context {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		select {
		case <-sig:
			cancel(nil)
		case err, ok := <-rabbit:
			if ok {
				cancel(rabbitTypes.ErrConnectionClosed)
			} else {
				cancel(err)
			}
		}
	}()

	return ctx
}
