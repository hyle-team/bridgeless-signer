package processor

import (
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/pkg/tokens"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessFormWithdrawalRequest(req bridgeTypes.FormWithdrawalRequest) (request *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId.String())
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}

	dstTokenAddress, err := p.tokenPairer.GetDestinationTokenAddress(
		req.Data.DepositIdentifier.GetChainId(),
		req.Data.TokenAddress,
		req.Data.DestinationChainId,
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
	req.Data.DestinationTokenAddress = &dstTokenAddress

	var tx *ethTypes.Transaction
	txConn := p.db.New()
	err = txConn.Transaction(func() error {
		tmpErr := txConn.SetDepositData(req.Data)
		if tmpErr != nil {
			return errors.Wrap(tmpErr, "failed to save deposit data")
		}

		tx, tmpErr = proxy.FormWithdrawalTransaction(req.Data)
		return errors.Wrap(tmpErr, "failed to form withdrawal transaction")
	})
	if err == nil {
		return &bridgeTypes.WithdrawalRequest{
			Data:        req.Data,
			DepositDbId: req.DepositDbId,
			Transaction: tx,
		}, false, nil
	}

	reprocessable = true
	if errors.Is(err, bridgeTypes.ErrInvalidReceiverAddress) {
		reprocessable = false
	}

	return nil, reprocessable, errors.Wrap(err, "failed to form withdrawal transaction")
}
