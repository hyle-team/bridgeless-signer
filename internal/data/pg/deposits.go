package pg

import (
	"database/sql"
	"encoding/hex"
	"github.com/hyle-team/bridgeless-signer/resources"
	"github.com/lib/pq"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/hyle-team/bridgeless-signer/internal/data"
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

	depositsDepositor        = "depositor"
	depositsDepositAmount    = "deposit_amount"
	depositsWithdrawalAmount = "withdrawal_amount"
	depositsDepositToken     = "deposit_token"
	depositsReceiver         = "receiver"
	depositsWithdrawalToken  = "withdrawal_token"
	depositsDepositBlock     = "deposit_block"

	depositsWithdrawalTxHash  = "withdrawal_tx_hash"
	depositsWithdrawalChainId = "withdrawal_chain_id"

	depositsSubmitStatus = "submit_status"

	depositIsWrappedToken = "is_wrapped_token"

	depositSignature = "signature"
)

type depositsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func (d *depositsQ) New() data.DepositsQ {
	return NewDepositsQ(d.db.Clone())
}

func (d *depositsQ) SetWithdrawalTxs(txs ...data.WithdrawalTx) error {
	if len(txs) == 0 {
		return nil
	}

	var (
		hashes = make(pq.StringArray, len(txs))
		chains = make(pq.StringArray, len(txs))
		ids    = make(pq.Int64Array, len(txs))
	)
	for i, tx := range txs {
		hashes[i] = strings.ToLower(tx.TxHash)
		chains[i] = tx.ChainId
		ids[i] = tx.DepositId
	}

	const query string = `
UPDATE deposits
SET
    status = $1,
    withdrawal_tx_hash = unnested_data.tx_hash,
    withdrawal_chain_id = unnested_data.chain_id
FROM (
	SELECT unnest($2::text[]) as tx_hash,
    	   unnest($3::text[]) as chain_id,
    	   unnest($4::bigint[]) as deposit_id
) as unnested_data
WHERE deposits.id = unnested_data.deposit_id
`

	return d.db.ExecRaw(query, resources.WithdrawalStatus_TX_PENDING, hashes, chains, ids)
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

func (d *depositsQ) UpdateWithdrawalStatus(status resources.WithdrawalStatus, ids ...int64) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsStatus, status).
		Where(squirrel.Eq{depositsId: ids})

	return d.db.Exec(stmt)
}

func (d *depositsQ) UpdateSubmitStatus(status resources.SubmitWithdrawalStatus, ids ...int64) error {
	stmt := squirrel.Update(depositsTable).
		Set(depositsSubmitStatus, status).
		Where(squirrel.Eq{depositsId: ids})

	return d.db.Exec(stmt)
}

func (d *depositsQ) SetDepositData(data data.DepositData) error {
	fields := map[string]interface{}{
		depositsDepositAmount:    data.DepositAmount.String(),
		depositsWithdrawalAmount: data.WithdrawalAmount.String(),
		depositsReceiver:         strings.ToLower(data.DestinationAddress),
		depositsDepositBlock:     data.Block,
		depositIsWrappedToken:    data.IsWrappedToken,
		// can be 0x00... in case of native ones
		depositsDepositToken: strings.ToLower(data.TokenAddress),
		depositsDepositor:    strings.ToLower(data.SourceAddress),
		// can be 0x00... in case of native ones
		depositsWithdrawalToken:   strings.ToLower(data.DestinationTokenAddress),
		depositsWithdrawalChainId: data.DestinationChainId,
	}

	return d.db.Exec(squirrel.Update(depositsTable).Where(
		squirrel.Eq{
			depositsTxHash:    data.TxHash,
			depositsTxEventId: data.TxEventId,
			depositsChainId:   data.ChainId,
		},
	).SetMap(fields))
}

func (d *depositsQ) SetDepositSignature(data data.DepositData) error {
	fields := map[string]interface{}{
		depositSignature: strings.ToLower(hex.EncodeToString(data.Signature)),
		depositsStatus:   resources.WithdrawalStatus_WITHDRAWAL_SIGNED,
	}

	return d.db.Exec(squirrel.Update(depositsTable).Where(
		squirrel.Eq{
			depositsTxHash:    data.TxHash,
			depositsTxEventId: data.TxEventId,
			depositsChainId:   data.ChainId,
		},
	).SetMap(fields))
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
		sql = sql.Where(squirrel.Eq{depositsSubmitStatus: resources.SubmitWithdrawalStatus_NOT_SUBMITTED})
	}

	return sql
}
