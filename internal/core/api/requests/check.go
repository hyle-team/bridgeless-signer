package requests

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	apiTypes "github.com/hyle-team/bridgeless-signer/internal/core/api/types"
	"github.com/hyle-team/bridgeless-signer/resources"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func CheckWithdrawalRequest(request *resources.CheckWithdrawalRequest, proxies bridgeTypes.ProxiesRepository) (*resources.WithdrawalRequest, error) {
	if request == nil {
		return nil, errors.New("request is not provided")
	}

	result, err := ToWithdrawRequest(request.OriginTxId)
	if err != nil {
		return nil, err
	}

	if _, err = ValidateWithdrawalRequest(result, proxies); err != nil {
		return nil, err
	}

	return result, nil
}

func ToWithdrawRequest(originTxId string) (*resources.WithdrawalRequest, error) {
	params := strings.Split(originTxId, "-")
	if len(params) != 3 {
		return nil, apiTypes.ErrInvalidOriginTxId
	}

	txEventIndex, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return nil, apiTypes.ErrInvalidOriginTxId
	}

	return &resources.WithdrawalRequest{
		Deposit: &resources.Deposit{
			TxHash:       params[0],
			TxEventIndex: txEventIndex,
			ChainId:      params[2],
		},
	}, nil
}
