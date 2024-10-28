package server

import (
	"context"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/ctx"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/requests"
	apiTypes "github.com/hyle-team/bridgeless-signer/internal/core/api/types"
	"github.com/hyle-team/bridgeless-signer/resources"

	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (grpcImplementation) SubmitWithdrawal(ctxt context.Context, request *resources.WithdrawalRequest) (*resources.Empty, error) {
	var (
		proxies  = ctx.Proxies(ctxt)
		db       = ctx.DB(ctxt)
		logger   = ctx.Logger(ctxt)
		producer = ctx.Producer(ctxt)
	)

	if err := requests.ValidateWithdrawalRequest(request, proxies); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	depositIdentifier := data.DepositIdentifier{
		TxHash:    request.Deposit.TxHash,
		TxEventId: int(request.Deposit.TxEventIndex),
		ChainId:   request.Deposit.ChainId,
	}
	deposit, err := db.Get(depositIdentifier)
	if err != nil {
		logger.WithError(err).Error("failed to get transaction")
		return nil, apiTypes.ErrInternal
	}

	if deposit != nil {
		if !deposit.Reprocessable() {
			return nil, apiTypes.ErrTxAlreadySubmitted
		}

		deposit.Status = resources.WithdrawalStatus_REPROCESSING
		if err = db.UpdateWithdrawalStatus(deposit.Status, deposit.Id); err != nil {
			logger.WithError(err).Error("failed to update transaction status")
			return nil, apiTypes.ErrInternal
		}
	} else {
		deposit = &data.Deposit{
			DepositIdentifier: depositIdentifier,
			Status:            resources.WithdrawalStatus_PROCESSING,
			SubmitStatus:      resources.SubmitWithdrawalStatus_NOT_SUBMITTED,
		}

		if deposit.Id, err = db.Insert(*deposit); err != nil {
			if errors.Is(err, data.ErrAlreadySubmitted) {
				return nil, apiTypes.ErrTxAlreadySubmitted
			}
			logger.WithError(err).Error("failed to insert transaction")
			return nil, apiTypes.ErrInternal
		}
	}

	if err = producer.PublishGetDepositRequest(bridgeTypes.GetDepositRequest{
		DepositDbId:       deposit.Id,
		DepositIdentifier: depositIdentifier,
	}); err != nil {
		logger.WithError(err).Error("failed to publish message")
		return nil, apiTypes.ErrInternal
	}

	return nil, nil
}
