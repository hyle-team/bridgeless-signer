package chain

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"reflect"
)

type Bitcoin struct {
	Rpc           *rpcclient.Client
	Receivers     []btcutil.Address
	Wallet        string
	Confirmations uint64
	Params        *chaincfg.Params
}

func (c Chain) Bitcoin() (Bitcoin, error) {
	if c.Type != types.ChainTypeBitcoin {
		return Bitcoin{}, errors.New("invalid chain type")
	}

	chain := Bitcoin{Wallet: c.Wallet, Confirmations: c.Confirmations}

	if chain.Wallet == "" {
		return chain, errors.New("wallet is not set")
	}
	if err := figure.Out(&chain.Rpc).FromInterface(c.Rpc).With(bitcoinHooks).Please(); err != nil {
		return chain, errors.Wrap(err, "failed to init bitcoin chain rpc")
	}

	var receivers []string
	if err := figure.Out(&receivers).FromInterface(c.BridgeAddresses).With(figure.BaseHooks).Please(); err != nil {
		return chain, errors.Wrap(err, "failed to decode bitcoin receivers")
	}
	if len(receivers) == 0 {
		return chain, errors.New("receivers list is empty")
	}

	if c.Network == NetworkMainnet {
		chain.Params = &chaincfg.MainNetParams
	}
	if c.Network == NetworkTestnet {
		chain.Params = &chaincfg.TestNet3Params
	}

	chain.Receivers = make([]btcutil.Address, len(receivers))
	for i, raw := range receivers {
		addr, err := btcutil.DecodeAddress(raw, chain.Params)
		if err != nil {
			return Bitcoin{}, errors.Wrap(err, fmt.Sprintf("failed to decode bitcoin receiver %s", raw))
		}

		chain.Receivers[i] = addr
	}

	return chain, nil
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
				Host:         clientConfig.Host,
				User:         clientConfig.User,
				Pass:         clientConfig.Pass,
				HTTPPostMode: true,
				DisableTLS:   clientConfig.DisableTLS,
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
