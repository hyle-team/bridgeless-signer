package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
)

func (p *bridgeProxy) FormWithdrawalTransaction(data data.DepositData) (*types.Transaction, error) {
	if data.DestinationChainId == nil || data.DestinationChainId.String() != p.chain.Id.String() {
		return nil, errors.New("invalid destination chain id")
	}

	if !common.IsHexAddress(data.DestinationAddress) {
		return nil, bridgeTypes.ErrInvalidReceiverAddress
	}

	if data.DestinationTokenAddress == nil {
		return nil, bridgeTypes.ErrDestinationTokenAddressRequired
	}

	// transact opts prevent the transaction from being sent to
	// the network, returning the transaction object only
	return p.bridgeContract.BridgeOut(
		bridgeOutTransactOpts(p.getTransactionNonce()),
		data.TokenAddress,
		common.HexToAddress(data.DestinationAddress),
		data.Amount,
		data.OriginTxId(),
	)
}

func (p *bridgeProxy) SendWithdrawalTransaction(signedTx *types.Transaction) error {
	return errors.Wrap(
		p.chain.Rpc.SendTransaction(context.Background(), signedTx),
		"failed to send withdrawal transaction",
	)
}

func (p *bridgeProxy) getTransactionNonce() *big.Int {
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
