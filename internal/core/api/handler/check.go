package handler

import (
	"context"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"

	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
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

		st, err := proxy.GetTransactionStatus(*tx.WithdrawalTxHash)
		if err != nil {
			h.logger.WithError(err).Error("failed to get tx receipt")
			return nil, ErrInternal
		}

		if st != bridgeTypes.TransactionStatusPending {
			switch st {
			case bridgeTypes.TransactionStatusFailed:
				tx.Status = types.WithdrawalStatus_TX_FAILED
			case bridgeTypes.TransactionStatusSuccessful:
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
