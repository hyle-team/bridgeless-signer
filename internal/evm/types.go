package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyle-team/bridgeless-signer/internal/data"
)

const DepositEvent = "BridgeIn"

type BridgeProxy interface {
	SupportsChain(chainId string) bool
}

type DepositData struct {
	data.DepositIdentifier

	DestinationChainId *big.Int
	SourceAddress      common.Address
	DestinationAddress string
	Amount             *big.Int
	TokenAddress       common.Address
}

func (d DepositData) OriginTxId() string {
	return d.DepositIdentifier.String()
}
