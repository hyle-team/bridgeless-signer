package types

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
)

type WithdrawRequest struct {
	DepositDbId int64
	Data        DepositData
	Transaction *types.Transaction
}

type GetDepositRequest struct {
	DepositDbId       int64
	DepositIdentifier data.DepositIdentifier
}

type FormWithdrawRequest struct {
	DepositDbId int64
	Data        DepositData
}
