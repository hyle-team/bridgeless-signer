package chain

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Chain struct {
	Id            *big.Int          `fig:"id,required"`
	Type          types.ChainType   `fig:"type,required"`
	EvmRpc        *ethclient.Client `fig:"evm_rpc"`
	BitcoinRpc    *rpcclient.Client `fig:"bitcoin_rpc"`
	BridgeAddress common.Address    `fig:"bridge_address"`
	Confirmations int64             `fig:"confirmations,required"`
}

func (c Chain) Validate() {
	if c.Id == nil {
		panic("chain id is nil")
	}

	if c.Type == types.ChainTypeEVM && c.EvmRpc == nil {
		panic("EVM rpc client is nil")
	}

	if c.Type == types.ChainTypeBitcoin && c.BitcoinRpc == nil {
		panic("Bitcoin rpc client is nil")
	}
}
