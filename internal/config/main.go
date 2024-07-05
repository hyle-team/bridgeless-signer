package config

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm/chain"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	"github.com/hyle-team/bridgeless-signer/internal/grpc"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	grpc.RESTGatewayConfigurer
	signer.Signerer
	chain.Chainer
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	getter kv.Getter
	grpc.RESTGatewayConfigurer
	signer.Signerer
	chain.Chainer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:                getter,
		Databaser:             pgdb.NewDatabaser(getter),
		Listenerer:            comfig.NewListenerer(getter),
		Logger:                comfig.NewLogger(getter, comfig.LoggerOpts{}),
		RESTGatewayConfigurer: grpc.NewRESTGatewayConfigurer(getter),
		Signerer:              signer.NewSignerer(getter),
		Chainer:               chain.NewChainer(getter),
	}
}
