package tokens

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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
	GetDestinationTokenInfo(
		srcChainId string,
		srcTokenAddr common.Address,
		dstChainId string,
	) (common.Address, bool, error)
}
