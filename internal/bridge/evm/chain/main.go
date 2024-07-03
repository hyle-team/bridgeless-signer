package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Chain struct {
	Id            *big.Int          `fig:"id,required"`
	Rpc           *ethclient.Client `fig:"rpc,required"`
	BridgeAddress common.Address    `fig:"bridge_address,required"`
	Confirmations int64             `fig:"confirmations,required"`
}
