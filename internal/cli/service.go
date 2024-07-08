package cli

import (
	"context"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/internal/data/pg"
	"github.com/hyle-team/bridgeless-signer/internal/grpc"
	"github.com/hyle-team/bridgeless-signer/internal/grpc/handler"
	"github.com/pkg/errors"
)

func RunService(cfg config.Config) error {
	// TODO: add proper ctx configuration
	ctx := context.Background()
	signer := cfg.Signer()

	proxies, err := evm.NewProxiesRepository(cfg.Chains(), signer.Address())
	if err != nil {
		panic(errors.Wrap(err, "failed to create proxies repository"))
	}

	srv := grpc.NewServer(
		cfg.Listener(),
		cfg.RESTGatewayConfig(),
		handler.NewServiceHandler(
			pg.NewDepositsQ(cfg.DB()),
			cfg.Log().WithField("service", "REST handler"),
			proxies,
		),
	)

	go func() {
		if err = srv.RunRESTGateway(ctx); err != nil {
			cfg.Log().WithError(err).Fatal("rest gateway error occurred")
		}
	}()

	cfg.Log().Info("service started")

	if err = srv.RunGRPC(context.Background()); err != nil {
		return err
	}

	return nil
}
