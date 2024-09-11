package types

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
)

type WithdrawalRequest struct {
	DepositDbId int64
	Data        data.DepositData
	Transaction *types.Transaction
}

type GetDepositRequest struct {
	DepositDbId       int64
	DepositIdentifier data.DepositIdentifier
}

type FormWithdrawalRequest struct {
	DepositDbId int64
	Destination ChainType
	Data        data.DepositData
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
