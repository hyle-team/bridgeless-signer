package data

import (
	"fmt"
	"math/big"

	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const OriginTxIdPattern = "%s-%d-%s"

var ErrAlreadySubmitted = errors.New("transaction already submitted")

type DepositsQ interface {
	New() DepositsQ
	Insert(Deposit) (id int64, err error)
	Get(identifier DepositIdentifier) (*Deposit, error)
	UpdateStatus(id int64, status types.WithdrawalStatus) error
	SetWithdrawalTx(depositId int64, txHash, chainId string) error
	Transaction(f func() error) error
}

type DepositIdentifier struct {
	TxHash    string `structs:"tx_hash" db:"tx_hash"`
	TxEventId int    `structs:"tx_event_id" db:"tx_event_id"`
	ChainId   string `structs:"chain_id" db:"chain_id"`
}

func (d DepositIdentifier) GetChainId() *big.Int {
	id, ok := new(big.Int).SetString(d.ChainId, 10)
	if !ok {
		return big.NewInt(0)
	}

	return id
}

func (d DepositIdentifier) String() string {
	return fmt.Sprintf(OriginTxIdPattern, d.TxHash, d.TxEventId, d.ChainId)
}

type Deposit struct {
	Id int64 `structs:"-" db:"id"`
	DepositIdentifier
	Status            types.WithdrawalStatus `structs:"status" db:"status"`
	WithdrawalTxHash  *string                `structs:"withdrawal_tx_hash" db:"withdrawal_tx_hash"`
	WithdrawalChainId *string                `structs:"withdrawal_chain_id" db:"withdrawal_chain_id"`
}

func (d Deposit) Reprocessable() bool {
	return d.Status == types.WithdrawalStatus_FAILED || d.Status == types.WithdrawalStatus_TX_FAILED
}

func (d Deposit) ToStatusResponse() *types.CheckWithdrawalResponse {
	result := &types.CheckWithdrawalResponse{
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
