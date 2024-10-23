package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/zano"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
)

type WithdrawalRequest struct {
	DepositDbId int64
	Destination types.ChainType
	Data        data.DepositData
}

type GetDepositRequest struct {
	DepositDbId       int64
	DepositIdentifier data.DepositIdentifier
}

type ZanoSignedWithdrawalRequest struct {
	DepositDbId int64
	Data        data.DepositData
	Transaction zano.SignedTransaction
}

func (r WithdrawalRequest) Id() int64 {
	return r.DepositDbId
}

type SubmitTransactionRequest struct {
	DepositDbId int64
}

func (r SubmitTransactionRequest) Id() int64 {
	return r.DepositDbId
}
