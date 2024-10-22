package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/btc"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
)

func (p *Processor) ProcessSendBitcoinWithdrawals(reqs ...bridgeTypes.WithdrawalRequest) (reprocessable bool, err error) {
	if len(reqs) == 0 {
		return false, nil
	}

	var (
		params        = make(map[string]*big.Int, len(reqs))
		withdrawalTxs = make([]data.WithdrawalTx, len(reqs))
		depositIds    = make([]int64, len(reqs))
	)
	for i, req := range reqs {
		params[req.Data.DestinationAddress] = req.Data.WithdrawalAmount
		depositIds[i] = req.DepositDbId
	}

	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, depositIds...) }()

	proxy, err := p.proxies.Proxy(reqs[0].Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return false, bridgeTypes.ErrChainNotSupported
		}
		return true, errors.Wrap(err, "failed to get proxy")
	}
	if proxy.Type() != bridgeTypes.ChainTypeBitcoin {
		return false, bridgeTypes.ErrChainNotSupported
	}
	btcProxy, ok := proxy.(btc.BridgeProxy)
	if !ok {
		return false, bridgeTypes.ErrChainNotSupported
	}

	hash, err := btcProxy.SendBitcoins(params)
	if err != nil {
		return true, errors.Wrap(err, "failed to send withdrawals")
	}

	for i, req := range reqs {
		withdrawalTxs[i] = data.WithdrawalTx{
			DepositId: req.DepositDbId,
			ChainId:   req.Data.DestinationChainId,
			TxHash:    hash,
		}
	}

	if err = p.db.New().SetWithdrawalTxs(withdrawalTxs...); err != nil {
		return false, errors.Wrap(err, "failed to set withdrawals")
	}

	return false, nil
}
