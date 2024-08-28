package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/connectors/core"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/tokens"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
)

type Processor struct {
	proxies       bridgeTypes.ProxiesRepository
	db            data.DepositsQ
	signer        *signer.Signer
	tokenPairer   tokens.TokenPairer
	coreConnector *core.Connector
}

func New(
	proxies bridgeTypes.ProxiesRepository,
	db data.DepositsQ,
	signer *signer.Signer,
	tokenPairer tokens.TokenPairer,
	coreConnector *core.Connector,

) *Processor {
	return &Processor{proxies: proxies, db: db, signer: signer, tokenPairer: tokenPairer, coreConnector: coreConnector}
}

func (p *Processor) SetWithdrawStatusFailed(id int64) error {
	return errors.Wrap(p.db.UpdateWithdrawalStatus(id, types.WithdrawalStatus_FAILED), "failed to update deposit status")
}

func (p *Processor) SetSubmitStatusFailed(ids ...int64) error {
	return errors.Wrap(p.db.UpdateSubmitStatus(types.SubmitWithdrawalStatus_SUBMIT_FAILED, ids...), "failed to update submit status")
}

func (p *Processor) updateInvalidDepositStatus(id int64, err error, reprocessable bool) error {
	if err == nil || reprocessable {
		return err
	}

	if tempErr := p.db.UpdateWithdrawalStatus(id, types.WithdrawalStatus_INVALID); tempErr != nil {
		return errors.Wrap(tempErr, "failed to update deposit status")
	}

	return err
}
