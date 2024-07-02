package signature

import (
	"crypto/ecdsa"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Signerer interface {
	Signer() Signer
}

type config struct {
	PrivKey *ecdsa.PrivateKey `fig:"eth_signer,required"`
}

type signerer struct {
	signerOnce comfig.Once
	getter     kv.Getter
}

func NewSignerer(getter kv.Getter) Signerer {
	return &signerer{
		getter: getter,
	}
}

func (s *signerer) Signer() Signer {
	return s.signerOnce.Do(func() interface{} {
		var cfg config

		err := figure.
			Out(&cfg).
			With(figure.EthereumHooks).
			From(kv.MustGetStringMap(s.getter, "signer")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out signer"))
		}

		return NewSigner(cfg.PrivKey)
	}).(Signer)
}
