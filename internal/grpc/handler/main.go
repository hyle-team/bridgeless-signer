package handler

import (
	"context"

	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"gitlab.com/distributed_lab/logan/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const internalError = "internal error"

var (
	ErrInternal           = status.Error(codes.Internal, internalError)
	ErrTxAlreadySubmitted = status.Error(codes.AlreadyExists, "transaction already submitted")
)

// ServiceHandler is an implementation of the API interface.
type ServiceHandler struct {
	db     data.DepositsQ
	logger logan.Entry
}

func NewServiceHandler(
	db data.DepositsQ,
	logger logan.Entry,
) *ServiceHandler {
	return &ServiceHandler{
		db:     db,
		logger: logger,
	}
}

func (h *ServiceHandler) SubmitWithdraw(_ context.Context, request *types.WithdrawRequest) error {
	if err := ValidateWithdrawRequest(request); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	dbconn := h.db.New()
	depositIdentifier := data.DepositIdentifier{
		TxHash:    request.Deposit.TxHash,
		TxEventId: int(request.Deposit.TxEventIndex),
		ChainId:   request.Deposit.ChainId,
	}
	deposit, err := dbconn.Get(depositIdentifier)
	if err != nil {
		h.logger.WithError(err).Error("failed to get transaction")
		return ErrInternal
	}

	if deposit != nil {
		if !deposit.Reprocessable() {
			return ErrTxAlreadySubmitted
		}

		deposit.Status = types.WithdrawStatus_PROCESSING
		if err = dbconn.UpdateStatus(deposit.Id, deposit.Status); err != nil {
			h.logger.WithError(err).Error("failed to update transaction status")
			return ErrInternal
		}
	} else {
		deposit = &data.Deposit{
			DepositIdentifier: depositIdentifier,
			Status:            types.WithdrawStatus_PROCESSING,
		}

		if deposit.Id, err = dbconn.Insert(*deposit); err != nil {
			if err == data.ErrAlreadySubmitted {
				return ErrTxAlreadySubmitted
			}
			h.logger.WithError(err).Error("failed to insert transaction")
			return ErrInternal
		}
	}

	//TODO add message to the  AMQP queue

	return nil
}

func (h *ServiceHandler) CheckWithdraw(_ context.Context, request *types.WithdrawRequest) (*types.CheckWithdrawResponse, error) {
	if err := ValidateWithdrawRequest(request); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	dbconn := h.db.New()
	depositIdentifier := data.DepositIdentifier{
		TxHash:    request.Deposit.TxHash,
		TxEventId: int(request.Deposit.TxEventIndex),
		ChainId:   request.Deposit.ChainId,
	}
	tx, err := dbconn.Get(depositIdentifier)
	if err != nil {
		h.logger.WithError(err).Error("failed to get transaction")
		return nil, ErrInternal
	}
	if tx == nil {
		return nil, status.Error(codes.NotFound, "transaction not found")
	}

	if tx.Status == types.WithdrawStatus_TX_PENDING && tx.WithdrawalTxHash != nil {
		// TODO: get tx status from the withdrawal chain

		// TODO: update tx status in the database

	}

	return tx.ToStatusResponse(), nil
}
