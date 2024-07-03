package evm

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *bridgeProxy) FormWithdrawalTransaction(data bridgeTypes.DepositData) (*types.Transaction, error) {
	if data.DestinationChainId == nil || data.DestinationChainId.String() != p.chain.Id.String() {
		return nil, errors.New("invalid destination chain id")
	}

	if !common.IsHexAddress(data.DestinationAddress) {
		return nil, bridgeTypes.ErrInvalidReceiverAddress
	}

	// transact opts prevent the transaction from being sent to
	// the network, returning the transaction object only
	return p.bridgeContract.BridgeOut(
		bridgeOutTransactOpts(),
		data.TokenAddress,
		common.HexToAddress(data.DestinationAddress),
		data.Amount,
		data.String(),
	)
}

func (p *bridgeProxy) SendWithdrawalTransaction(signedTx *types.Transaction) error {
	return errors.Wrap(
		p.chain.Rpc.SendTransaction(context.Background(), signedTx),
		"failed to send withdrawal transaction",
	)
}

func bridgeOutTransactOpts() *bind.TransactOpts {
	const gasLimit = 300000

	return &bind.TransactOpts{
		GasLimit: gasLimit,
		// prevent the transaction from being sent to the network
		NoSend: true,
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			// skip signing
			return transaction, nil
		},
	}
}
