package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/zano"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessZanoSendWithdrawalRequest(req ZanoSignedWithdrawalRequest) (res *SubmitTransactionRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}
	zanoProxy, ok := proxy.(zano.BridgeProxy)
	if !ok {
		return nil, false, bridgeTypes.ErrInvalidProxyType
	}

	hash, err := zanoProxy.EmitAssetSigned(req.Transaction)
	if err != nil {
		return nil, true, errors.Wrap(err, "failed to broadcast withdrawal")
	}

	withdrawalTx := data.WithdrawalTx{
		DepositId: req.DepositDbId,
		ChainId:   req.Data.DestinationChainId,
		TxHash:    hash,
	}
	if err = p.db.New().SetWithdrawalTxs(withdrawalTx); err != nil {
		return nil, false, errors.Wrap(err, "failed to set withdrawal")
	}

	return &SubmitTransactionRequest{DepositDbId: req.DepositDbId}, false, nil
}
