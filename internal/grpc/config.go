package grpc

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type RESTGatewayConfigurer interface {
	RESTGatewayConfig() RESTGatewayConfig
}

type RESTGatewayConfig struct {
	Address string `fig:"addr,required"`
}

func NewRESTGatewayConfigurer(getter kv.Getter) RESTGatewayConfigurer {
	return &gatewayConfigurer{
		getter: getter,
	}
}

type gatewayConfigurer struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *gatewayConfigurer) RESTGatewayConfig() RESTGatewayConfig {
	return c.once.Do(func() interface{} {
		const yamlKey = "rest_gateway"
		var conf RESTGatewayConfig

		if err := figure.
			Out(&conf).
			From(kv.MustGetStringMap(c.getter, yamlKey)).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to configure REST gateway"))
		}

		return conf
	}).(RESTGatewayConfig)
}
