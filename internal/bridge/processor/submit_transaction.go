package processor

import (
	coretypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *Processor) SubmitTransactions(reqs ...bridgeTypes.SubmitTransactionRequest) (reprocessable bool, err error) {
	if len(reqs) == 0 {
		return false, nil
	}

	depositIds := make([]int64, 0, len(reqs))
	for i, req := range reqs {
		depositIds[i] = req.DepositDbId
	}

	var selectSubmitted = false
	deposits, err := p.db.Select(data.DepositsSelector{Ids: depositIds, Submitted: &selectSubmitted})
	if err != nil {
		return true, errors.Wrap(err, "failed to get deposits")
	}

	depositTxs := make([]coretypes.Transaction, 0, len(deposits))
	for i, d := range deposits {
		depositTxs[i] = d.ToTransaction()
	}

	if err = p.coreConnector.SubmitDeposits(depositTxs...); err != nil {
		return true, errors.Wrap(err, "failed to submit deposits")
	}

	return false, nil
}
