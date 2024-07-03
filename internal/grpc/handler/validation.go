package handler

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var txHashPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)

func (h *ServiceHandler) ValidateWithdrawRequest(request *types.WithdrawRequest) error {
	if request == nil {
		return errors.New("request is nil")
	}

	deposit := request.Deposit
	if deposit == nil {
		return errors.New("deposit is nil")
	}

	err := validation.Errors{
		"tx_hash":     validation.Validate(deposit.TxHash, validation.Required, validation.Match(txHashPattern)),
		"tx_event_id": validation.Validate(deposit.TxEventIndex, validation.Required, validation.Min(0)),
		"chain_id":    validation.Validate(deposit.ChainId, validation.Required),
	}.Filter()

	if err == nil {
		if !h.proxyRepo.SupportsChain(deposit.ChainId) {
			return ErrChainNotSupported
		}
	}

	return err
}
