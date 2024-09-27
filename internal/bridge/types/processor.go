package types

import (
	"github.com/hyle-team/bridgeless-signer/internal/data"
)

type WithdrawalRequest struct {
	DepositDbId int64
	Destination ChainType
	Data        data.DepositData
}

type GetDepositRequest struct {
	DepositDbId       int64
	DepositIdentifier data.DepositIdentifier
}

type BitcoinWithdrawalRequest struct {
	DepositDbId int64
	Data        data.DepositData
}

func (b BitcoinWithdrawalRequest) Id() int64 {
	return b.DepositDbId
}

type SubmitTransactionRequest struct {
	DepositDbId int64
}

func (r SubmitTransactionRequest) Id() int64 {
	return r.DepositDbId
}
