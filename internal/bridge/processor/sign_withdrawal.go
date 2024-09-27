package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}

	signHash, err := proxy.GetSignHash(req.Data)
	if err != nil {
		return nil, true, errors.Wrap(err, "failed to form withdrawal transaction")
	}

	signature, err := p.signer.SignMessage(signHash)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to sign message")
	}
	req.Data.Signature = signature

	return &bridgeTypes.WithdrawalRequest{
		Data:        req.Data,
		DepositDbId: req.DepositDbId,
	}, false, nil
}
