package types

import (
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	"github.com/pkg/errors"
)

var (
	ErrPairNotFound      = errors.New("pair not found")
	ErrTokenInfoNotFound = errors.New("token info not found")
)

type TokenPairer interface {
	GetDestinationTokenInfo(
		srcChainId string,
		srcTokenAddr common.Address,
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