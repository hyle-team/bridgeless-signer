package chain

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

type Network string

const (
	NetworkMainnet Network = "mainnet"
	NetworkTestnet Network = "testnet"
)

func (n Network) Validate() error {
	switch n {
	case NetworkMainnet, NetworkTestnet:
		return nil
	default:
		return errors.New("invalid network")
	}
}

type Chain struct {
	Id            string          `fig:"id,required"`
	Type          types.ChainType `fig:"type,required"`
	Confirmations uint64          `fig:"confirmations,required"`

	// EVM configuration
	EvmRpc        *ethclient.Client `fig:"evm_rpc"`
	BridgeAddress common.Address    `fig:"bridge_address"`

	// Bitcoin configuration
	BitcoinRpc *rpcclient.Client `fig:"bitcoin_rpc"`
	// BitcoinReceivers is a list of allowed addresses (of different types) to receive deposits
	BitcoinReceivers []string `fig:"bitcoin_receivers"`
	Network          Network  `fig:"network"`
}

type Bitcoin struct {
	Rpc           *rpcclient.Client
	Receivers     []btcutil.Address
	Confirmations uint64
	Params        *chaincfg.Params
}

type Evm struct {
	Rpc           *ethclient.Client
	BridgeAddress common.Address
	Confirmations uint64
}

func (c Chain) Bitcoin() (Bitcoin, error) {
	if c.Type != types.ChainTypeBitcoin {
		return Bitcoin{}, errors.New("invalid chain type")
	}
	if c.BitcoinRpc == nil {
		return Bitcoin{}, errors.New("rpc client is nil")
	}
	if len(c.BitcoinReceivers) == 0 {
		return Bitcoin{}, errors.New("receivers list is empty")
	}

	var params *chaincfg.Params
	if c.Network == NetworkMainnet {
		params = &chaincfg.MainNetParams
	}
	if c.Network == NetworkTestnet {
		params = &chaincfg.TestNet3Params
	}

	receivers := make([]btcutil.Address, len(c.BitcoinReceivers))
	for i, raw := range c.BitcoinReceivers {
		addr, err := btcutil.DecodeAddress(raw, params)
		if err != nil {
			return Bitcoin{}, errors.Wrap(err, fmt.Sprintf("failed to decode bitcoin receiver %s", raw))
		}

		receivers[i] = addr
	}

	return Bitcoin{
		Rpc:           c.BitcoinRpc,
		Receivers:     receivers,
		Confirmations: c.Confirmations,
		Params:        params,
	}, nil
}

func (c Chain) Evm() (Evm, error) {
	if c.Type != types.ChainTypeEVM {
		return Evm{}, errors.New("invalid chain type")
	}
	if c.EvmRpc == nil {
		return Evm{}, errors.New("rpc client is nil")
	}
	if c.BridgeAddress == (common.Address{}) {
		return Evm{}, errors.New("bridge address is empty")
	}

	return Evm{
		Rpc:           c.EvmRpc,
		BridgeAddress: c.BridgeAddress,
		Confirmations: c.Confirmations,
	}, nil
}
