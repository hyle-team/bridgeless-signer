package processor

import (
	coretypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
)

func (p *Processor) SubmitTransactions(reqs ...bridgeTypes.SubmitTransactionRequest) (reprocessable bool, err error) {
	if len(reqs) == 0 {
		return false, nil
	}

	depositIds := make([]int64, len(reqs))
	for i, req := range reqs {
		depositIds[i] = req.DepositDbId
	}

	var selectSubmitted = false
	deposits, err := p.db.Select(data.DepositsSelector{Ids: depositIds, Submitted: &selectSubmitted})
	if err != nil {
		return true, errors.Wrap(err, "failed to get deposits")
	}

	depositTxs := make([]coretypes.Transaction, len(deposits))
	for i, d := range deposits {
		depositTxs[i] = d.ToTransaction()
	}

	// rollback if transaction failed to be sent
	txConn := p.db.New()
	err = txConn.Transaction(func() error {
		if tmperr := txConn.UpdateSubmitStatus(types.SubmitWithdrawalStatus_SUCCESSFUL, depositIds...); tmperr != nil {
			return errors.Wrap(tmperr, "failed to set deposits submitted")
		}

		return errors.Wrap(p.coreConnector.SubmitDeposits(depositTxs...), "failed to submit deposits")
	})

	return err != nil, err
}
