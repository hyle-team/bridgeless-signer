package handler

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"

	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *ServiceHandler) CheckWithdrawal(_ context.Context, request *types.CheckWithdrawalRequest) (*types.CheckWithdrawalResponse, error) {
	wr, err := h.CheckWithdrawalRequest(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dbconn := h.db.New()
	depositIdentifier := data.DepositIdentifier{
		TxHash:    wr.Deposit.TxHash,
		TxEventId: int(wr.Deposit.TxEventIndex),
		ChainId:   wr.Deposit.ChainId,
	}
	tx, err := dbconn.Get(depositIdentifier)
	if err != nil {
		h.logger.WithError(err).Error("failed to get deposit")
		return nil, ErrInternal
	}
	if tx == nil {
		return nil, status.Error(codes.NotFound, "deposit not found")
	}

	if tx.Status == types.WithdrawalStatus_TX_PENDING && tx.WithdrawalTxHash != nil {
		proxy, err := h.proxies.Proxy(*tx.WithdrawalChainId)
		if err != nil {
			h.logger.WithError(err).Error("failed to get proxy")
			return nil, ErrInternal // should not happen if the chain is supported, but just in case
		}

		receipt, err := proxy.GetTransactionReceipt(common.HexToHash(*tx.WithdrawalTxHash))
		if err != nil {
			// omitting only pending txs
			if !errors.Is(err, bridgeTypes.ErrTxPending) {
				// if the tx is still pending, we return the same status
				// otherwise, render error
				h.logger.WithError(err).Error("failed to get tx receipt")
				return nil, ErrInternal
			}
		} else {
			switch receipt.Status {
			case ethTypes.ReceiptStatusFailed:
				tx.Status = types.WithdrawalStatus_TX_FAILED
			case ethTypes.ReceiptStatusSuccessful:
				tx.Status = types.WithdrawalStatus_TX_SUCCESSFUL
			}

			// updating in the db
			if err = dbconn.UpdateWithdrawalStatus(tx.Id, tx.Status); err != nil {
				h.logger.WithError(err).Error("failed to update transaction status")
				return nil, ErrInternal
			}
		}
	}

	return tx.ToStatusResponse(), nil
}
