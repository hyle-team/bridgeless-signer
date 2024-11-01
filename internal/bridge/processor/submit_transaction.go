package processor

import (
	coretypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/resources"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessSubmitTransactions(reqs ...SubmitTransactionRequest) (reprocessable bool, err error) {
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
		if tmperr := txConn.UpdateSubmitStatus(resources.SubmitWithdrawalStatus_SUCCESSFUL, depositIds...); tmperr != nil {
			return errors.Wrap(tmperr, "failed to set deposits submitted")
		}

		err = p.core.SubmitDeposits(depositTxs...)
		// ignoring already submitted transaction
		if errors.Is(err, bridgeTypes.ErrTransactionAlreadySubmitted) {
			err = nil
		}

		return errors.Wrap(err, "failed to submit deposits")
	})
	if errors.Is(err, bridgeTypes.ErrTransactionAlreadySubmitted) {
		return false, err
	}

	return err != nil, err
}
