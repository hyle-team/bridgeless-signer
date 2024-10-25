package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/evm"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessEthSignWithdrawalRequest(req WithdrawalRequest) (res *SubmitTransactionRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}
	evmProxy, ok := proxy.(evm.BridgeProxy)
	if !ok {
		return nil, false, bridgeTypes.ErrInvalidProxyType
	}

	signHash, err := evmProxy.GetSignHash(req.Data)
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

	return &SubmitTransactionRequest{
		DepositDbId: req.DepositDbId,
	}, false, nil
}
