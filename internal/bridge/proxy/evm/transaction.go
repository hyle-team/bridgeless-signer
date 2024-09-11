package evm

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *proxy) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	tx, pending, err := p.chain.Rpc.TransactionByHash(ctx, txHash)
	if err != nil {
		if err.Error() == "not found" {
			return nil, bridgeTypes.ErrTxNotFound
		}

		return nil, errors.Wrap(err, "failed to get transaction by hash")
	}
	if pending {
		return nil, bridgeTypes.ErrTxPending
	}

	receipt, err := p.chain.Rpc.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx receipt")
	}
	if receipt == nil {
		return nil, errors.New("receipt is nil")
	}

	return receipt, nil
}

func (p *proxy) GetTransactionStatus(txHash string) (bridgeTypes.TransactionStatus, error) {
	receipt, err := p.GetTransactionReceipt(common.HexToHash(txHash))
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrTxPending) {
			return bridgeTypes.TransactionStatusPending, nil
		}
		if errors.Is(err, bridgeTypes.ErrTxNotFound) {
			return bridgeTypes.TransactionStatusNotFound, nil
		}

		return bridgeTypes.TransactionStatusUnknown, errors.Wrap(err, "failed to get transaction receipt")
	}

	switch receipt.Status {
	case ethTypes.ReceiptStatusSuccessful:
		return bridgeTypes.TransactionStatusSuccessful, nil
	case ethTypes.ReceiptStatusFailed:
		return bridgeTypes.TransactionStatusFailed, nil
	default:
		return bridgeTypes.TransactionStatusUnknown, nil
	}
}
