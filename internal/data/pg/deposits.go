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
			depositsTxHash:    deposit.TxHash,
			depositsTxEventId: deposit.TxEventId,
			depositsChainId:   deposit.ChainId,
			depositsStatus:    deposit.Status,
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

func (d *depositsQ) UpdateStatus(id int64, status types.WithdrawalStatus) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsStatus, status).
		Where(squirrel.Eq{depositsId: id})

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
