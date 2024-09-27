package evm

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types/operations"
	"github.com/hyle-team/bridgeless-signer/internal/data"
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
		operation := operations.WithdrawNativeContent{
			Amount:   data.DepositAmount.Bytes(),
			Receiver: operations.To32Bytes(hexutil.MustDecode(data.DestinationAddress)),
			Txhash:   operations.To32Bytes(hexutil.MustDecode(data.TxHash)),
			Txnonce:  operations.IntTo32Bytes(data.TxEventId),
			ChainID:  operations.To32Bytes(hexutil.MustDecode(data.ChainId)),
		}
		return operation.CalculateHash(), nil
	}

	operation := operations.WithdrawERC20Content{
		Amount:       data.DepositAmount.Bytes(),
		Receiver:     operations.To32Bytes(hexutil.MustDecode(data.DestinationAddress)),
		TxHash:       operations.To32Bytes(hexutil.MustDecode(data.TxHash)),
		TxNonce:      operations.IntTo32Bytes(data.TxEventId),
		ChainID:      operations.To32Bytes(hexutil.MustDecode(data.ChainId)),
		TokenAddress: operations.To32Bytes(data.TokenAddress.Bytes()),
		IsWrapped:    operations.BoolToBytes(data.IsWrappedToken),
	}
	return operation.CalculateHash(), nil

}

func (p *proxy) getTransactionNonce() *big.Int {
	p.nonceM.Lock()
	defer p.nonceM.Unlock()

	nonce := big.NewInt(0).SetUint64(p.signerNonce)
	p.signerNonce++

	return nonce
}
