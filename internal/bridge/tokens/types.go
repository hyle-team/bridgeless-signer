package tokens

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrSourceTokenNotSupported      = fmt.Errorf("source token not supported")
	ErrDestinationTokenNotSupported = fmt.Errorf("destination token not supported")
)

type TokenPairerConfiger interface {
	TokenPairer() TokenPairer
}

type TokenPairer interface {
	GetDestinationTokenAddress(
		srcChainId *big.Int,
		srcTokenAddr common.Address,
		dstChainId *big.Int,
	) (common.Address, error)
}
