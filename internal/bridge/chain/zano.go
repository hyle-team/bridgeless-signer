package chain

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/pkg/zano"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"reflect"
)

type Zano struct {
	Client        *zano.Sdk
	Confirmations uint64
	Receivers     []string
}

func (c Chain) Zano() (Zano, error) {
	if c.Type != types.ChainTypeZano {
		return Zano{}, errors.New("invalid chain type")
	}

	chain := Zano{Confirmations: c.Confirmations}
	if err := figure.Out(&chain.Receivers).FromInterface(c.BridgeAddresses).With(figure.BaseHooks).Please(); err != nil {
		return chain, errors.Wrap(err, "failed to decode zano receivers")
	}
	if err := figure.Out(&chain.Client).FromInterface(c.Rpc).With(zanoHooks).Please(); err != nil {
		return chain, errors.Wrap(err, "failed to decode zano client")
	}

	return chain, nil
}

var zanoHooks = figure.Hooks{
	"*zano.Sdk": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case map[string]interface{}:
			var clientConfig struct {
				DaemonRpc string `fig:"daemon_rpc,required"`
				WalletRpc string `fig:"wallet_rpc,required"`
			}

			if err := figure.Out(&clientConfig).With(figure.BaseHooks).From(v).Please(); err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to figure out zano rpc config")
			}

			sdk := zano.NewSDK(clientConfig.WalletRpc, clientConfig.DaemonRpc)
			return reflect.ValueOf(sdk), nil
		default:
			return reflect.Value{}, errors.Errorf("unsupported conversion from %T", value)
		}
	},
}
