package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	HTTPGatewayConfigurer
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	getter kv.Getter
	HTTPGatewayConfigurer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:                getter,
		Databaser:             pgdb.NewDatabaser(getter),
		Listenerer:            comfig.NewListenerer(getter),
		Logger:                comfig.NewLogger(getter, comfig.LoggerOpts{}),
		HTTPGatewayConfigurer: NewHTTPGatewayConfigurer(getter),
	}
}
