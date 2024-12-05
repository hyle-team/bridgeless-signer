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

func (r WithdrawalRequest) Id() int64 {
	return r.DepositDbId
}

type GetDepositRequest struct {
	DepositDbId       int64
	DepositIdentifier data.DepositIdentifier
}

func (r GetDepositRequest) Id() int64 { return r.DepositDbId }

type ZanoSignedWithdrawalRequest struct {
	DepositDbId int64
	Data        data.DepositData
	Transaction zano.SignedTransaction
}

func (r ZanoSignedWithdrawalRequest) Id() int64 { return r.DepositDbId }

type SubmitTransactionRequest struct {
	DepositDbId int64
}

func (r SubmitTransactionRequest) Id() int64 {
	return r.DepositDbId
}
