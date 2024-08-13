package pg

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	depositsTable     = "deposits"
	depositsTxHash    = "tx_hash"
	depositsTxEventId = "tx_event_id"
	depositsChainId   = "chain_id"
	depositsStatus    = "status"
	depositsId        = "id"

	depositsDepositor       = "depositor"
	depositsAmount          = "amount"
	depositsDepositToken    = "deposit_token"
	depositsReceiver        = "receiver"
	depositsWithdrawalToken = "withdrawal_token"
	depositsDepositBlock    = "deposit_block"

	depositsWithdrawalTxHash  = "withdrawal_tx_hash"
	depositsWithdrawalChainId = "withdrawal_chain_id"

	depositsSubmitStatus = "submit_status"
)

type depositsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func (d *depositsQ) New() data.DepositsQ {
	return NewDepositsQ(d.db.Clone())
}

func (d *depositsQ) SetWithdrawalTx(depositId int64, txHash, chainId string) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsStatus, types.WithdrawalStatus_TX_PENDING).
		Set(depositsWithdrawalTxHash, txHash).
		Set(depositsWithdrawalChainId, chainId).
		Where(squirrel.Eq{depositsId: depositId})

	return d.db.Exec(stmt)
}

func (d *depositsQ) Insert(deposit data.Deposit) (int64, error) {
	stmt := squirrel.
		Insert(depositsTable).
		SetMap(map[string]interface{}{
			depositsTxHash:       deposit.TxHash,
			depositsTxEventId:    deposit.TxEventId,
			depositsChainId:      deposit.ChainId,
			depositsStatus:       deposit.Status,
			depositsSubmitStatus: deposit.SubmitStatus,
		}).
		Suffix("RETURNING id")

	var id int64
	if err := d.db.Get(&id, stmt); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			err = data.ErrAlreadySubmitted
		}

		return id, err
	}

	return id, nil
}

func (d *depositsQ) Get(identifier data.DepositIdentifier) (*data.Deposit, error) {
	var deposit data.Deposit
	err := d.db.Get(&deposit, d.selector.Where(squirrel.Eq{
		depositsTxHash:    identifier.TxHash,
		depositsTxEventId: identifier.TxEventId,
		depositsChainId:   identifier.ChainId,
	}))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &deposit, err
}

func (d *depositsQ) Select(selector data.DepositsSelector) ([]data.Deposit, error) {
	query := d.applySelector(selector, d.selector)
	var deposits []data.Deposit
	if err := d.db.Select(&deposits, query); err != nil {
		return nil, err
	}

	return deposits, nil
}

func (d *depositsQ) UpdateWithdrawalStatus(id int64, status types.WithdrawalStatus) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsStatus, status).
		Where(squirrel.Eq{depositsId: id})

	return d.db.Exec(stmt)
}

func (d *depositsQ) UpdateSubmitStatus(status types.SubmitWithdrawalStatus, ids ...int64) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsSubmitStatus, status).
		Where(squirrel.Eq{depositsId: ids})

	return d.db.Exec(stmt)
}

func (d *depositsQ) SetDepositData(data data.DepositData) error {
	fields := map[string]interface{}{
		depositsDepositor:    strings.ToLower(data.SourceAddress.String()),
		depositsAmount:       data.Amount.String(),
		depositsDepositToken: strings.ToLower(data.TokenAddress.String()),
		depositsReceiver:     strings.ToLower(data.DestinationAddress),
		depositsDepositBlock: data.Block,
	}

	if data.DestinationTokenAddress != nil {
		fields[depositsWithdrawalToken] = strings.ToLower(data.DestinationTokenAddress.String())
	}

	return d.db.Exec(squirrel.Update(depositsTable).SetMap(fields))
}

func NewDepositsQ(db *pgdb.DB) data.DepositsQ {
	return &depositsQ{
		db:       db.Clone(),
		selector: squirrel.Select("*").From(depositsTable),
	}
}

func (d *depositsQ) Transaction(f func() error) error {
	return d.db.Transaction(f)
}

func (d *depositsQ) applySelector(selector data.DepositsSelector, sql squirrel.SelectBuilder) squirrel.SelectBuilder {
	if len(selector.Ids) > 0 {
		sql = sql.Where(squirrel.Eq{depositsId: selector.Ids})
	}

	if selector.Submitted != nil {
		sql = sql.Where(squirrel.Eq{depositsSubmitStatus: types.SubmitWithdrawalStatus_NOT_SUBMITTED})
	}

	return sql
}
