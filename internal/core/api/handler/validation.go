package handler

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
)

func (h *ServiceHandler) ValidateWithdrawalRequest(request *types.WithdrawalRequest) error {
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

	proxy, err := h.proxies.Proxy(deposit.ChainId)
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

func (h *ServiceHandler) CheckWithdrawalRequest(request *types.CheckWithdrawalRequest) (*types.WithdrawalRequest, error) {
	if request == nil {
		return nil, errors.New("request is not provided")
	}

	result, err := toWithdrawRequest(request.OriginTxId)
	if err != nil {
		return nil, err
	}

	return result, h.ValidateWithdrawalRequest(result)
}

func toWithdrawRequest(originTxId string) (*types.WithdrawalRequest, error) {
	params := strings.Split(originTxId, "-")
	if len(params) != 3 {
		return nil, ErrInvalidOriginTxId
	}

	txEventIndex, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return nil, ErrInvalidOriginTxId
	}

	return &types.WithdrawalRequest{
		Deposit: &types.Deposit{
			TxHash:       params[0],
			TxEventIndex: txEventIndex,
			ChainId:      params[2],
		},
	}, nil
}
