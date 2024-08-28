package tokens

import (
	"github.com/pkg/errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrPairNotFound                 = errors.New("pair not found")
	ErrSourceTokenNotSupported      = errors.New("source token not supported")
	ErrDestinationTokenNotSupported = errors.New("destination token not supported")
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
