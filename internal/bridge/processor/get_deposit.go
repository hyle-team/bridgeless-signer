package processor

import (
	"fmt"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/btc"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/pkg/tokens"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessGetDepositRequest(req bridgeTypes.GetDepositRequest) (data *bridgeTypes.FormWithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.DepositIdentifier.ChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return data, false, errors.Wrap(err, fmt.Sprintf("chain id: %v", req.DepositIdentifier.ChainId))
		}
		return data, true, errors.Wrap(err, "failed to get source proxy")
	}

	depositData, err := proxy.GetDepositData(req.DepositIdentifier)
	if err != nil {
		reprocessable = true
		if errors.Is(err, bridgeTypes.ErrTxFailed) ||
			errors.Is(err, bridgeTypes.ErrDepositNotFound) ||
			errors.Is(err, bridgeTypes.ErrInvalidDepositedAmount) ||
			errors.Is(err, bridgeTypes.ErrInvalidScriptPubKey) {
			reprocessable = false
		}

		return nil, reprocessable, errors.Wrap(err, "failed to get deposit data")
	}

	dstProxy, err := p.proxies.Proxy(depositData.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return data, false, errors.Wrap(err, fmt.Sprintf("chain id: %v", req.DepositIdentifier.ChainId))
		}
		return data, true, errors.Wrap(err, "failed to get destination proxy")
	}
	if !dstProxy.AddressValid(depositData.DestinationAddress) {
		return data, false, errors.Wrap(bridgeTypes.ErrInvalidReceiverAddress, depositData.DestinationAddress)
	}

	switch dstProxy.Type() {
	case bridgeTypes.ChainTypeBitcoin:
		if depositData.Amount.Int64() < btc.MinSatoshisPerOutput {
			return data, false, bridgeTypes.ErrInvalidDepositedAmount
		}
	case bridgeTypes.ChainTypeEVM:
		depositData.DestinationTokenAddress, err = p.tokenPairer.GetDestinationTokenAddress(
			depositData.ChainId,
			depositData.TokenAddress,
			depositData.DestinationChainId,
		)
		if err != nil {
			reprocessable = true
			if errors.Is(err, tokens.ErrSourceTokenNotSupported) ||
				errors.Is(err, tokens.ErrDestinationTokenNotSupported) ||
				errors.Is(err, tokens.ErrPairNotFound) {
				reprocessable = false
			}

			return nil, reprocessable, errors.Wrap(err, "failed to get destination token address")
		}
	default:
		return data, false, errors.Wrap(err, fmt.Sprintf("invalid chain type: %v", dstProxy.Type()))
	}

	if err = p.db.New().SetDepositData(*depositData); err != nil {
		return nil, true, errors.Wrap(err, "failed to save deposit data")
	}

	return &bridgeTypes.FormWithdrawalRequest{
		DepositDbId: req.DepositDbId,
		Data:        *depositData,
		Destination: dstProxy.Type(),
	}, false, nil
}
