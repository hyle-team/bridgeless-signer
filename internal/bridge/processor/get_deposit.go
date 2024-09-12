package processor

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/btc"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessGetDepositRequest(req bridgeTypes.GetDepositRequest) (data *bridgeTypes.FormWithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.DepositIdentifier.ChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return data, false, errors.Wrap(err, fmt.Sprintf("source chain id: %v", req.DepositIdentifier.ChainId))
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
			return data, false, errors.Wrap(err, fmt.Sprintf("destination chain id: %v", depositData.DestinationChainId))
		}
		return data, true, errors.Wrap(err, "failed to get destination proxy")
	}
	if !dstProxy.AddressValid(depositData.DestinationAddress) {
		return data, false, errors.Wrap(bridgeTypes.ErrInvalidReceiverAddress, depositData.DestinationAddress)
	}

	srcTokenInfo, err := p.core.GetTokenInfo(depositData.ChainId, depositData.TokenAddress.String())
	if err != nil {
		reprocessable = true
		if errors.Is(err, bridgeTypes.ErrTokenInfoNotFound) {
			reprocessable = false
		}
		return nil, reprocessable, errors.Wrap(err, "failed to get source token info")
	}
	dstTokenInfo, err := p.core.GetDestinationTokenInfo(
		depositData.ChainId,
		depositData.TokenAddress,
		depositData.DestinationChainId,
	)
	if err != nil {
		reprocessable = true
		if errors.Is(err, bridgeTypes.ErrPairNotFound) {
			reprocessable = false
		}
		return nil, reprocessable, errors.Wrap(err, "failed to get destination token info")
	}

	transformAmount(depositData.Amount, srcTokenInfo.Decimals, dstTokenInfo.Decimals)

	switch dstProxy.Type() {
	case bridgeTypes.ChainTypeBitcoin:
		if depositData.Amount.IsInt64() && depositData.Amount.Int64() < btc.MinSatoshisPerOutput {
			return nil, false, bridgeTypes.ErrInvalidDepositedAmount
		}
	case bridgeTypes.ChainTypeEVM:
		if depositData.Amount.IsInt64() && depositData.Amount.Int64() == 0 {
			return nil, false, bridgeTypes.ErrInvalidDepositedAmount
		}

		depositData.DestinationTokenAddress = common.HexToAddress(dstTokenInfo.Address)
		depositData.IsWrappedToken = dstTokenInfo.IsWrapped
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
