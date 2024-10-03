package handler

import (
	"context"

	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *ServiceHandler) SubmitWithdrawal(_ context.Context, request *types.WithdrawalRequest) error {
	if err := h.ValidateWithdrawalRequest(request); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
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
		deposit.Status = types.WithdrawalStatus_REPROCESSING
		if err = dbconn.UpdateWithdrawalStatus(deposit.Status, deposit.Id); err != nil {
			h.logger.WithError(err).Error("failed to update transaction status")
			return ErrInternal
		}
	} else {
		deposit = &data.Deposit{
			DepositIdentifier: depositIdentifier,
			Status:            types.WithdrawalStatus_PROCESSING,
			SubmitStatus:      types.SubmitWithdrawalStatus_NOT_SUBMITTED,
		}

		if deposit.Id, err = dbconn.Insert(*deposit); err != nil {
			if errors.Is(err, data.ErrAlreadySubmitted) {
				return ErrTxAlreadySubmitted
			}
			h.logger.WithError(err).Error("failed to insert transaction")
			return ErrInternal
		}
	}

	if err = h.publisher.SendGetDepositRequest(bridgeTypes.GetDepositRequest{
		DepositDbId:       deposit.Id,
		DepositIdentifier: depositIdentifier,
	}); err != nil {
		h.logger.WithError(err).Error("failed to publish message")
		return ErrInternal
	}

	return nil
}
