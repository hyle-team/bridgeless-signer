package pg

import (
	"database/sql"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	depositsTable     = "deposits"
	depositsTxHash    = "tx_hash"
	depositsTxEventId = "tx_event_id"
	depositsChainId   = "chain_id"
	depositsStatus    = "status"
	depositsId        = "id"
)

type depositsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func (d *depositsQ) New() data.DepositsQ {
	return d.NewDepositsQ(d.db.Clone())
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
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &deposit, err
}

func (d *depositsQ) UpdateStatus(id int64, status types.WithdrawStatus) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsStatus, status).
		Where(squirrel.Eq{depositsId: id})

	return d.db.Exec(stmt)
}

func (d *depositsQ) NewDepositsQ(db *pgdb.DB) data.DepositsQ {
	return &depositsQ{
		db:       db,
		selector: squirrel.Select("*").From(depositsTable),
	}
}
