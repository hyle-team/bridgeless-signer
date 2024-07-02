package data

import (
	"fmt"

	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var ErrAlreadySubmitted = errors.New("transaction already submitted")

type DepositsQ interface {
	New() DepositsQ
	Insert(Deposit) (id int64, err error)
	Get(identifier DepositIdentifier) (*Deposit, error)
	UpdateStatus(id int64, status types.WithdrawStatus) error
}

type DepositIdentifier struct {
	TxHash    string `structs:"tx_hash" db:"tx_hash"`
	TxEventId int    `structs:"tx_event_id" db:"tx_event_id"`
	ChainId   string `structs:"chain_id" db:"chain_id"`
}

func (d DepositIdentifier) String() string {
	return fmt.Sprintf("%s-%d-%s", d.TxHash, d.TxEventId, d.ChainId)
}

type Deposit struct {
	Id int64 `structs:"-" db:"id"`
	DepositIdentifier
	Status            types.WithdrawStatus `structs:"status" db:"status"`
	WithdrawalTxHash  *string              `structs:"withdrawal_tx_hash" db:"withdrawal_tx_hash"`
	WithdrawalChainId *int64               `structs:"withdrawal_chain_id" db:"withdrawal_chain_id"`
}

func (d Deposit) Reprocessable() bool {
	return d.Status == types.WithdrawStatus_FAILED
}

func (d Deposit) ToStatusResponse() *types.CheckWithdrawResponse {
	result := &types.CheckWithdrawResponse{
		Status: d.Status,
	}

	if d.WithdrawalTxHash != nil && d.WithdrawalChainId != nil {
		result.ResultTransaction = &types.Transaction{
			Hash:    *d.WithdrawalTxHash,
			ChainId: *d.WithdrawalChainId,
		}
	}

	return result
}
