package config

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	core "github.com/hyle-team/bridgeless-signer/internal/bridge/core/config"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	api "github.com/hyle-team/bridgeless-signer/internal/core/api/config"
	rabbit "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/config"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	signer.Signerer
	chain.Chainer
	rabbit.Rabbitter
	core.ConnectorConfigurer
	api.Listenerer
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	getter kv.Getter
	signer.Signerer
	chain.Chainer
	rabbit.Rabbitter
	core.ConnectorConfigurer
	api.Listenerer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:              getter,
		Databaser:           pgdb.NewDatabaser(getter),
		Logger:              comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Signerer:            signer.NewSignerer(getter),
		Chainer:             chain.NewChainer(getter),
		Rabbitter:           rabbit.NewRabbitter(getter),
		ConnectorConfigurer: core.NewConnectorConfigurer(getter),
		Listenerer:          api.NewListenerer(getter),
	}
}
