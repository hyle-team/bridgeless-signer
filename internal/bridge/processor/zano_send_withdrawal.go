package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessZanoSendWithdrawalRequest(req bridgeTypes.ZanoSignedWithdrawalRequest) (res *bridgeTypes.SubmitTransactionRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}

	hash, err := proxy.EmitAssetSigned(req.Transaction)
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

	return &bridgeTypes.SubmitTransactionRequest{
		DepositDbId: req.DepositDbId,
	}, false, nil
}
