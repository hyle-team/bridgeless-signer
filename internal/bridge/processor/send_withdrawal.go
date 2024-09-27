package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessSendWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	// ensure that withdrawal request was not already processed
	deposit, err := p.db.Get(req.Data.DepositIdentifier)
	if err != nil {
		return true, errors.Wrap(err, "failed to check if deposit already processed")
	}
	if deposit == nil {
		return true, errors.New("deposit was not found in the database")
	}

	// rollback if transaction failed to be sent
	txConn := p.db.New()
	err = txConn.Transaction(func() error {
		if tempErr := txConn.SetWithdrawalTxs(data.WithdrawalTx{
			DepositId: req.DepositDbId,
			ChainId:   req.Data.DestinationChainId,
		}); tempErr != nil {
			return errors.Wrap(tempErr, "failed to set withdrawal tx")
		}

		return nil
	})
	return err != nil, err
}
