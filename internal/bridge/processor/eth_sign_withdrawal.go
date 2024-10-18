package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessEthSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.SubmitTransactionRequest, reprocessable bool, err error) {
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
		return nil, true, errors.Wrap(err, "failed to sign message")
	}
	req.Data.Signature = signature

	if err = p.db.New().SetDepositSignature(req.Data); err != nil {
		return nil, true, errors.Wrap(err, "failed to save signature data")
	}

	return &bridgeTypes.SubmitTransactionRequest{
		DepositDbId: req.DepositDbId,
	}, false, nil
}
