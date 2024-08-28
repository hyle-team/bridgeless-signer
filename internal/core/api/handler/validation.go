package handler

import (
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var txHashPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)

func (h *ServiceHandler) ValidateWithdrawalRequest(request *types.WithdrawalRequest) error {
	if request == nil {
		return errors.New("request is not provided")
	}

	deposit := request.Deposit
	if deposit == nil {
		return errors.New("deposit is not provided")
	}

	err := validation.Errors{
		"tx_hash":        validation.Validate(deposit.TxHash, validation.Required, validation.Match(txHashPattern)),
		"tx_event_index": validation.Validate(deposit.TxEventIndex, validation.Min(0)),
		"chain_id":       validation.Validate(deposit.ChainId, validation.Required),
	}.Filter()

	if err == nil {
		if !h.proxies.SupportsChain(deposit.ChainId) {
			return ErrChainNotSupported
		}
	}

	return err
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
