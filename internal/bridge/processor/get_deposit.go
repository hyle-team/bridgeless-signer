package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessGetDepositRequest(req bridgeTypes.GetDepositRequest) (data *bridgeTypes.FormWithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	proxy, err := p.proxies.Proxy(req.DepositIdentifier.ChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return data, false, bridgeTypes.ErrChainNotSupported
		}
		return data, true, errors.Wrap(err, "failed to get proxy")
	}

	depositData, err := proxy.GetDepositData(req.DepositIdentifier)
	if err == nil {
		return &bridgeTypes.FormWithdrawalRequest{
			DepositDbId: req.DepositDbId,
			Data:        *depositData,
		}, false, nil
	}

	reprocessable = true
	if errors.Is(err, bridgeTypes.ErrTxFailed) ||
		errors.Is(err, bridgeTypes.ErrDepositNotFound) {
		reprocessable = false
	}

	return nil, reprocessable, errors.Wrap(err, "failed to get deposit data")
}
