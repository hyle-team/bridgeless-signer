package evm

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/evm/operations"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
)

func (p *proxy) WithdrawalAmountValid(amount *big.Int) bool {
	if amount.Cmp(big.NewInt(0)) != 1 {
		return false
	}

	return true
}

func (p *proxy) GetSignHash(data data.DepositData) ([]byte, error) {
	if IsAddressEmpty(data.DestinationTokenAddress) {
		operation, err := operations.NewWithdrawNativeContent(data)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create WithdrawNativeContent operation")
		}

		return operation.CalculateHash(), nil
	}

	operation, err := operations.NewWithdrawERC20Content(data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create WithdrawERC20Content operation")
	}

	return operation.CalculateHash(), nil
}
