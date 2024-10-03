package evm

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types/operations"
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
	chainID, ok := new(big.Int).SetString(data.DestinationChainId, 10)
	if !ok {
		return nil, errors.New("invalid chain id")
	}

	if IsAddressEmpty(data.DestinationTokenAddress) {
		operation := operations.WithdrawNativeContent{
			Amount:   operations.To32Bytes(data.DepositAmount.Bytes()),
			Receiver: hexutil.MustDecode(data.DestinationAddress),
			TxHash:   hexutil.MustDecode(data.TxHash),
			TxNonce:  operations.IntTo32Bytes(data.TxEventId),
			ChainID:  operations.To32Bytes(chainID.Bytes()),
		}
		return operation.CalculateHash(), nil
	}

	operation := operations.WithdrawERC20Content{
		Amount:                  operations.To32Bytes(data.DepositAmount.Bytes()),
		Receiver:                hexutil.MustDecode(data.DestinationAddress),
		TxHash:                  hexutil.MustDecode(data.TxHash),
		TxNonce:                 operations.IntTo32Bytes(data.TxEventId),
		ChainID:                 operations.To32Bytes(chainID.Bytes()),
		DestinationTokenAddress: data.DestinationTokenAddress.Bytes(),
		IsWrapped:               operations.BoolToBytes(data.IsWrappedToken),
	}
	return operation.CalculateHash(), nil
}
