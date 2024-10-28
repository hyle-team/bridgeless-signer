package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/resources"
	"github.com/pkg/errors"
)

func ValidateWithdrawalRequest(request *resources.WithdrawalRequest, proxies bridgeTypes.ProxiesRepository) error {
	if request == nil {
		return errors.New("request is not provided")
	}

	deposit := request.Deposit
	if deposit == nil {
		return errors.New("deposit is not provided")
	}

	err := validation.Errors{
		"tx_hash":        validation.Validate(deposit.TxHash, validation.Required),
		"tx_event_index": validation.Validate(deposit.TxEventIndex, validation.Min(0)),
		"chain_id":       validation.Validate(deposit.ChainId, validation.Required),
	}.Filter()

	proxy, err := proxies.Proxy(deposit.ChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return err
		}

		return errors.Wrap(err, "failed to get proxy")
	}

	if !proxy.TransactionHashValid(deposit.TxHash) {
		return validation.Errors{"tx_hash": errors.New("invalid transaction hash")}
	}

	return nil
}
