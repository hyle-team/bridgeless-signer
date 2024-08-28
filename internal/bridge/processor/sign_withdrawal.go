package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	tx, err := p.signer.SignTx(req.Transaction, req.Data.DestinationChainId)
	if err != nil {
		// TODO: should be reprocessable or not?
		return res, true, errors.Wrap(err, "failed to sign withdrawal transaction")
	}

	return &bridgeTypes.WithdrawalRequest{
		Data:        req.Data,
		DepositDbId: req.DepositDbId,
		Transaction: tx,
	}, false, nil
}
