package processor

import (
	coretypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/tokens"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
)

type TxSubmitter interface {
	SubmitDeposits(depositTxs ...coretypes.Transaction) error
}

type Processor struct {
	proxies     bridgeTypes.ProxiesRepository
	db          data.DepositsQ
	signer      *signer.Signer
	tokenPairer tokens.TokenPairer
	submitter   TxSubmitter
}

func New(
	proxies bridgeTypes.ProxiesRepository,
	db data.DepositsQ,
	signer *signer.Signer,
	tokenPairer tokens.TokenPairer,
	submitter TxSubmitter,

) *Processor {
	return &Processor{proxies: proxies, db: db, signer: signer, tokenPairer: tokenPairer, submitter: submitter}
}

func (p *Processor) SetWithdrawStatusFailed(ids ...int64) error {
	return errors.Wrap(p.db.UpdateWithdrawalStatus(types.WithdrawalStatus_FAILED, ids...), "failed to update deposit status")
}

func (p *Processor) SetSubmitStatusFailed(ids ...int64) error {
	return errors.Wrap(p.db.UpdateSubmitStatus(types.SubmitWithdrawalStatus_SUBMIT_FAILED, ids...), "failed to update submit status")
}

func (p *Processor) updateInvalidDepositStatus(err error, reprocessable bool, ids ...int64) error {
	if err == nil || reprocessable {
		return err
	}

	if tempErr := p.db.UpdateWithdrawalStatus(types.WithdrawalStatus_INVALID, ids...); tempErr != nil {
		return errors.Wrap(tempErr, "failed to update deposit status")
	}

	return err
}
