package data

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	bridgetypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"math/big"
)

const OriginTxIdPattern = "%s-%d-%s"

var ErrAlreadySubmitted = errors.New("transaction already submitted")

type DepositsQ interface {
	New() DepositsQ
	Insert(Deposit) (id int64, err error)
	Select(selector DepositsSelector) ([]Deposit, error)
	Get(identifier DepositIdentifier) (*Deposit, error)
	SetDepositData(data DepositData) error
	UpdateWithdrawalStatus(status types.WithdrawalStatus, ids ...int64) error
	UpdateSubmitStatus(status types.SubmitWithdrawalStatus, ids ...int64) error
	SetWithdrawalTxs(txs ...WithdrawalTx) error
	Transaction(f func() error) error
}

type WithdrawalTx struct {
	DepositId int64
	TxHash    string
	ChainId   string
}

type DepositIdentifier struct {
	TxHash    string `structs:"tx_hash" db:"tx_hash"`
	TxEventId int    `structs:"tx_event_id" db:"tx_event_id"`
	ChainId   string `structs:"chain_id" db:"chain_id"`
}

type DepositsSelector struct {
	Ids       []int64
	Submitted *bool
}

func (d DepositIdentifier) String() string {
	return fmt.Sprintf(OriginTxIdPattern, d.TxHash, d.TxEventId, d.ChainId)
}

type Deposit struct {
	Id int64 `structs:"-" db:"id"`
	DepositIdentifier
	Status types.WithdrawalStatus `structs:"status" db:"status"`

	Depositor       *string `structs:"depositor" db:"depositor"`
	Amount          *string `structs:"amount" db:"amount"`
	DepositToken    *string `structs:"deposit_token" db:"deposit_token"`
	Receiver        *string `structs:"receiver" db:"receiver"`
	WithdrawalToken *string `structs:"withdrawal_token" db:"withdrawal_token"`
	DepositBlock    *int64  `structs:"deposit_block" db:"deposit_block"`

	WithdrawalTxHash  *string `structs:"withdrawal_tx_hash" db:"withdrawal_tx_hash"`
	WithdrawalChainId *string `structs:"withdrawal_chain_id" db:"withdrawal_chain_id"`

	IsWrappedToken *bool `structs:"is_wrapped_token" db:"is_wrapped_token"`

	SubmitStatus types.SubmitWithdrawalStatus `structs:"submit_status" db:"submit_status"`
}

func (d Deposit) Reprocessable() bool {
	return d.Status == types.WithdrawalStatus_FAILED ||
		d.Status == types.WithdrawalStatus_TX_FAILED
}

func (d Deposit) WithdrawalAllowed() bool {
	if d.WithdrawalTxHash == nil {
		return true
	}

	return d.Status == types.WithdrawalStatus_REPROCESSING
}

func (d Deposit) ToStatusResponse() *types.CheckWithdrawalResponse {
	result := &types.CheckWithdrawalResponse{
		Status: d.Status,
		DepositData: &types.DepositData{
			EventIndex:      int64(d.TxEventId),
			Depositor:       d.Depositor,
			Amount:          d.Amount,
			DepositToken:    d.DepositToken,
			WithdrawalToken: d.WithdrawalToken,
			Receiver:        d.Receiver,
			BlockNumber:     d.DepositBlock,
			IsWrapped:       d.IsWrappedToken,
		},
		DepositTransaction: &types.Transaction{
			Hash:    d.TxHash,
			ChainId: d.ChainId,
		},
		SubmitStatus: d.SubmitStatus,
	}

	if d.WithdrawalTxHash != nil && d.WithdrawalChainId != nil {
		result.WithdrawalTransaction = &types.Transaction{
			Hash:    *d.WithdrawalTxHash,
			ChainId: *d.WithdrawalChainId,
		}
	}

	return result
}

func (d Deposit) ToTransaction() bridgetypes.Transaction {
	tx := bridgetypes.Transaction{
		DepositTxHash:    d.TxHash,
		DepositTxIndex:   uint64(d.TxEventId),
		DepositChainId:   d.ChainId,
		WithdrawalTxHash: stringOrEmpty(d.WithdrawalTxHash),
		Depositor:        stringOrEmpty(d.Depositor),
		Amount:           stringOrEmpty(d.Amount),
		DepositToken:     stringOrEmpty(d.DepositToken),
		Receiver:         stringOrEmpty(d.Receiver),
		WithdrawalToken:  stringOrEmpty(d.WithdrawalToken),

		WithdrawalChainId: stringOrEmpty(d.WithdrawalChainId),

		IsWrapped: boolOrFalse(d.IsWrappedToken),
	}

	if d.DepositBlock != nil {
		tx.DepositBlock = uint64(*d.DepositBlock)
	}

	return tx
}

type DepositData struct {
	DepositIdentifier
	DestinationChainId string

	SourceAddress      common.Address
	DestinationAddress string

	Amount *big.Int

	TokenAddress            common.Address
	DestinationTokenAddress common.Address

	IsWrappedToken bool

	Block int64
}

func (d DepositData) OriginTxId() string {
	return d.DepositIdentifier.String()
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func boolOrFalse(b *bool) bool {
	if b == nil {
		return false
	}

	return true
}
