package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/resources"
	"github.com/pkg/errors"
)

type Processor struct {
	proxies bridgeTypes.ProxiesRepository
	db      data.DepositsQ
	signer  *signer.Signer
	core    bridgeTypes.Bridger
}

func New(
	proxies bridgeTypes.ProxiesRepository,
	db data.DepositsQ,
	signer *signer.Signer,
	core bridgeTypes.Bridger,

) *Processor {
	return &Processor{proxies: proxies, db: db, signer: signer, core: core}
}

func (p *Processor) SetWithdrawStatusFailed(ids ...int64) error {
	return errors.Wrap(p.db.UpdateWithdrawalStatus(resources.WithdrawalStatus_FAILED, ids...), "failed to update deposit status")
}

func (p *Processor) SetSubmitStatusFailed(ids ...int64) error {
	return errors.Wrap(p.db.UpdateSubmitStatus(resources.SubmitWithdrawalStatus_SUBMIT_FAILED, ids...), "failed to update submit status")
}

func (p *Processor) updateInvalidDepositStatus(err error, reprocessable bool, ids ...int64) error {
	if err == nil || reprocessable {
		return err
	}

	if tempErr := p.db.UpdateWithdrawalStatus(resources.WithdrawalStatus_INVALID, ids...); tempErr != nil {
		return errors.Wrap(tempErr, "failed to update deposit status")
	}

	return err
}
