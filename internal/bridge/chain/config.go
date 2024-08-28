package chain

import (
	"github.com/btcsuite/btcd/rpcclient"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"reflect"
)

type Chainer interface {
	Chains() []Chain
}

type chainer struct {
	once   comfig.Once
	getter kv.Getter
}

func NewChainer(getter kv.Getter) Chainer {
	return &chainer{
		getter: getter,
	}
}

func (c *chainer) Chains() []Chain {
	return c.once.Do(func() interface{} {
		var cfg struct {
			Chains []Chain `fig:"list,required"`
		}

		if err := figure.
			Out(&cfg).
			With(figure.BaseHooks, figure.EthereumHooks, bitcoinHooks).
			From(kv.MustGetStringMap(c.getter, "chains")).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out chains"))
		}

		if len(cfg.Chains) == 0 {
			panic(errors.New("no chains were configured"))
		}

		return cfg.Chains
	}).([]Chain)
}

var bitcoinHooks = figure.Hooks{
	"*rpcclient.Client": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case map[string]interface{}:
			var clientConfig struct {
				Host       string `fig:"host,required"`
				User       string `fig:"user,required"`
				Pass       string `fig:"pass,required"`
				DisableTLS bool   `fig:"disable_tls"`
			}

			if err := figure.Out(&clientConfig).With(figure.BaseHooks).From(v).Please(); err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to figure out bitcoin rpc client config")
			}

			client, err := rpcclient.New(&rpcclient.ConnConfig{
				Host:       clientConfig.Host,
				User:       clientConfig.User,
				Pass:       clientConfig.Pass,
				DisableTLS: clientConfig.DisableTLS,
			}, nil)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to create bitcoin rpc client")
			}

			return reflect.ValueOf(client), nil
		default:
			return reflect.Value{}, errors.Errorf("unsupported conversion from %T", value)
		}
	},
}
