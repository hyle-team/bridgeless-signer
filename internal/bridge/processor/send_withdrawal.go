package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessSendWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	// ensure that withdrawal request was not already processed
	deposit, err := p.db.Get(req.Data.DepositIdentifier)
	if err != nil {
		return true, errors.Wrap(err, "failed to check if deposit already processed")
	}
	if deposit == nil {
		return true, errors.New("deposit was not found in the database")
	}
	if !deposit.WithdrawalAllowed() {
		return false, errors.New("withdrawal transaction was already sent")
	}

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId.String())
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return false, bridgeTypes.ErrChainNotSupported
		}
		return true, errors.Wrap(err, "failed to get proxy")
	}

	// rollback if transaction failed to be sent
	txConn := p.db.New()
	err = txConn.Transaction(func() error {
		if tempErr := txConn.SetWithdrawalTx(
			req.DepositDbId, req.Transaction.Hash().Hex(), req.Data.DestinationChainId.String(),
		); tempErr != nil {
			return errors.Wrap(tempErr, "failed to set withdrawal tx")
		}

		return errors.Wrap(proxy.SendWithdrawalTransaction(req.Transaction), "failed to send withdrawal transaction")
	})
	return err != nil, err
}
