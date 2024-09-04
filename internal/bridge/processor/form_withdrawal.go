package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessFormWithdrawalRequest(req bridgeTypes.FormWithdrawalRequest) (request *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}

	tx, err := proxy.FormWithdrawalTransaction(req.Data)
	if err == nil {
		return nil, true, errors.Wrap(err, "failed to form withdrawal transaction")
	}

	return &bridgeTypes.WithdrawalRequest{
		Data:        req.Data,
		DepositDbId: req.DepositDbId,
		Transaction: tx,
	}, false, nil
}
