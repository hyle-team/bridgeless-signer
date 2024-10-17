package types

import (
	coretypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	"github.com/pkg/errors"
)

const DefaultNativeTokenAddress = "0x0000000000000000000000000000000000000000"

var (
	ErrPairNotFound                = errors.New("pair not found")
	ErrTokenInfoNotFound           = errors.New("token info not found")
	ErrTransactionAlreadySubmitted = errors.New("transaction already submitted")
)

type TokenPairer interface {
	GetDestinationTokenInfo(
		srcChainId string,
		srcTokenAddr string,
		dstChainId string,
	) (coretypes.TokenInfo, error)
}

type Tokener interface {
	GetTokenInfo(chainId string, addr string) (coretypes.TokenInfo, error)
}

type TxSubmitter interface {
	SubmitDeposits(depositTxs ...coretypes.Transaction) error
}

type Bridger interface {
	Tokener
	TokenPairer
	TxSubmitter
}
