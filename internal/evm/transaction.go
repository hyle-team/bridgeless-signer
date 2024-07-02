package evm

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (p *bridgeProxy) GetTransactionReceipt(txHash common.Hash, chainId string) (*types.Receipt, error) {
	chain, ok := p.chains[chainId]
	if !ok {
		return nil, ErrChainNotSupported
	}

	ctx := context.Background()
	tx, pending, err := chain.Rpc.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction by hash")
	}
	if pending {
		return nil, ErrTxPending
	}

	receipt, err := chain.Rpc.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt")
	}
	if receipt == nil {
		return nil, errors.New("receipt is nil")
	}

	return receipt, nil
}
