package evm

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/evm/operations"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
)

type Operation interface {
	CalculateHash() []byte
}

func (p *proxy) WithdrawalAmountValid(amount *big.Int) bool {
	if amount.Cmp(bridge.ZeroAmount) != 1 {
		return false
	}

	return true
}

func (p *proxy) GetSignHash(data data.DepositData) ([]byte, error) {
	var operation Operation
	var err error

	if data.DestinationTokenAddress == bridge.DefaultNativeTokenAddress {
		operation, err = operations.NewWithdrawNativeContent(data)
	} else {
		operation, err = operations.NewWithdrawERC20Content(data)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to create operation")
	}

	hash := operation.CalculateHash()
	prefixedHash := operations.SetSignaturePrefix(hash)

	return prefixedHash, nil
}
