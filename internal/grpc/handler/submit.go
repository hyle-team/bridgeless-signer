package handler

import (
	"context"

	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
)

func (h *ServiceHandler) SubmitWithdraw(_ context.Context, request *types.WithdrawRequest) error {
	if err := h.ValidateWithdrawRequest(request); err != nil {
		return err
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
