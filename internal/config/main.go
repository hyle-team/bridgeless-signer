package config

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm/chain"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm/signature"
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
	signature.Signerer
	chain.Chainer
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	getter kv.Getter
	grpc.RESTGatewayConfigurer
	signature.Signerer
	chain.Chainer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:                getter,
		Databaser:             pgdb.NewDatabaser(getter),
		Listenerer:            comfig.NewListenerer(getter),
		Logger:                comfig.NewLogger(getter, comfig.LoggerOpts{}),
		RESTGatewayConfigurer: grpc.NewRESTGatewayConfigurer(getter),
		Signerer:              signature.NewSignerer(getter),
		Chainer:               chain.NewChainer(getter),
	}
}
