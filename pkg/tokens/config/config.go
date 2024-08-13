package config

import (
	"github.com/hyle-team/bridgeless-signer/pkg/tokens"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type configTokenPairerConfiger struct {
	once   comfig.Once
	getter kv.Getter
}

func NewConfigTokenPairerConfiger(getter kv.Getter) tokens.TokenPairerConfiger {
	return &configTokenPairerConfiger{
		getter: getter,
	}
}

func (c *configTokenPairerConfiger) TokenPairer() tokens.TokenPairer {
	return c.once.Do(func() interface{} {
		const tokenPairerConfigurerConfigKey = "tokens"
		var cfg struct {
			TokenPairs []TokenPairs `fig:"list,required"`
		}

		if err := figure.
			Out(&cfg).
			With(figure.BaseHooks, figure.EthereumHooks).
			From(kv.MustGetStringMap(c.getter, tokenPairerConfigurerConfigKey)).
			Please(); err != nil {
			panic(err)
		}

		return NewTokenPairer(cfg.TokenPairs)
	}).(tokens.TokenPairer)

}

type Token struct {
	ChainId *big.Int       `fig:"chain_id,required"`
	Address common.Address `fig:"address,required"`
}

type TokenPairs struct {
	Token Token   `fig:"token,required"`
	Pairs []Token `fig:"pairs,required"`
}
