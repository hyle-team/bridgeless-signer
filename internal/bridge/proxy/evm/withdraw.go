package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"math/big"
)

func (p *proxy) FormWithdrawalTransaction(data data.DepositData) (*types.Transaction, error) {
	// transact opts prevent the transaction from being sent to
	// the network, returning the transaction object only

	if IsAddressEmpty(data.DestinationTokenAddress) {
		return p.bridgeContract.BridgeOutNative(
			bridgeOutTransactOpts(p.getTransactionNonce()),
			common.HexToAddress(data.DestinationAddress),
			data.Amount,
			data.OriginTxId(),
		)
	}

	return p.bridgeContract.BridgeOut(
		bridgeOutTransactOpts(p.getTransactionNonce()),
		data.DestinationTokenAddress,
		common.HexToAddress(data.DestinationAddress),
		data.Amount,
		data.OriginTxId(),
		data.IsWrappedToken,
	)
}

func (p *proxy) SendWithdrawalTransaction(signedTx *types.Transaction) error {
	return p.chain.Rpc.SendTransaction(context.Background(), signedTx)
}

func (p *proxy) getTransactionNonce() *big.Int {
	p.nonceM.Lock()
	defer p.nonceM.Unlock()

	nonce := big.NewInt(0).SetUint64(p.signerNonce)
	p.signerNonce++

	return nonce
}

func bridgeOutTransactOpts(nonce *big.Int) *bind.TransactOpts {
	const gasLimit = 300000

	return &bind.TransactOpts{
		GasLimit: gasLimit,
		Nonce:    nonce,
		// prevent the transaction from being sent to the network
		NoSend: true,
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			// skip signing
			return transaction, nil
		},
	}
}
