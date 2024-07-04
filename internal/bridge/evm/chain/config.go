package chain

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
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
			With(figure.BaseHooks, figure.EthereumHooks).
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
