package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type HTTPGatewayConfigurer interface {
	HTTPGatewayConfig() HTTPGatewayConfig
}

type HTTPGatewayConfig struct {
	Address string `fig:"addr,required"`
}

func NewHTTPGatewayConfigurer(getter kv.Getter) HTTPGatewayConfigurer {
	return &gatewayConfigurer{
		getter: getter,
	}
}

type gatewayConfigurer struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *gatewayConfigurer) HTTPGatewayConfig() HTTPGatewayConfig {
	return c.once.Do(func() interface{} {
		const yamlKey = "http_gateway"
		var conf HTTPGatewayConfig

		conf.Address = ":9000"

		//if err := figure.
		//	Out(conf).
		//	From(kv.MustGetStringMap(c.getter, yamlKey)).
		//	Please(); err != nil {
		//	panic(errors.Wrap(err, "failed to configure http gateway"))
		//}

		return conf
	}).(HTTPGatewayConfig)
}
