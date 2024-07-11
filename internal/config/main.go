package config

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm/chain"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/tokens"
	tokensConfig "github.com/hyle-team/bridgeless-signer/internal/bridge/tokens/config"
	api "github.com/hyle-team/bridgeless-signer/internal/core/api/config"
	rabbit "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/config"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	api.RESTGatewayConfigurer
	signer.Signerer
	chain.Chainer
	rabbit.Rabbitter
	tokens.TokenPairerConfiger
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	getter kv.Getter
	api.RESTGatewayConfigurer
	signer.Signerer
	chain.Chainer
	tokens.TokenPairerConfiger
	rabbit.Rabbitter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:                getter,
		Databaser:             pgdb.NewDatabaser(getter),
		Listenerer:            comfig.NewListenerer(getter),
		Logger:                comfig.NewLogger(getter, comfig.LoggerOpts{}),
		RESTGatewayConfigurer: api.NewRESTGatewayConfigurer(getter),
		Signerer:              signer.NewSignerer(getter),
		Chainer:               chain.NewChainer(getter),
		Rabbitter:             rabbit.NewRabbitter(getter),
		TokenPairerConfiger:   tokensConfig.NewConfigTokenPairerConfiger(getter),
	}
}
