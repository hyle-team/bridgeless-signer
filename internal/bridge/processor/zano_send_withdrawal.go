package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *Processor) SendZanoWithdraw(req bridgeTypes.ZanoSignedWithdrawalRequest) (reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return false, bridgeTypes.ErrChainNotSupported
		}
		return true, errors.Wrap(err, "failed to get proxy")
	}

	hash, err := proxy.EmitAssetSigned(req.Transaction)
	if err != nil {
		return true, errors.Wrap(err, "failed to broadcast withdrawal")
	}

	withdrawalTx := data.WithdrawalTx{
		DepositId: req.DepositDbId,
		ChainId:   req.Data.DestinationChainId,
		TxHash:    hash,
	}
	if err = p.db.New().SetWithdrawalTxs(withdrawalTx); err != nil {
		return false, errors.Wrap(err, "failed to set withdrawal")
	}

	return false, nil
}
