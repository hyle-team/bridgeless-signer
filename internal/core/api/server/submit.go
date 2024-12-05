package server

import (
	"context"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
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

	chainType, err := requests.ValidateWithdrawalRequest(request, proxies)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	depositIdentifier := formDepositIdentifier(request, chainType)
	deposit, err := db.Get(depositIdentifier)
	if err != nil {
		logger.WithError(err).Error("failed to get transaction")
		return nil, apiTypes.ErrInternal
	}

	err = db.Transaction(func() error {
		if deposit != nil {
			if err = checkExistingDeposit(deposit, request, chainType); err != nil {
				return err
			}

			deposit.Status = resources.WithdrawalStatus_REPROCESSING
			if err = db.UpdateWithdrawalStatus(deposit.Status, deposit.Id); err != nil {
				logger.WithError(err).Error("failed to update transaction status")
				return apiTypes.ErrInternal
			}
		} else {
			deposit = &data.Deposit{
				DepositIdentifier: depositIdentifier,
				Status:            resources.WithdrawalStatus_PROCESSING,
				SubmitStatus:      resources.SubmitWithdrawalStatus_NOT_SUBMITTED,
			}

			if deposit.Id, err = db.Insert(*deposit); err != nil {
				if errors.Is(err, data.ErrAlreadySubmitted) {
					return apiTypes.ErrTxAlreadySubmitted
				}

				logger.WithError(err).Error("failed to insert transaction")
				return apiTypes.ErrInternal
			}
		}

		if err = producer.PublishGetDepositRequest(bridgeTypes.GetDepositRequest{
			DepositDbId:       deposit.Id,
			DepositIdentifier: depositIdentifier,
		}); err != nil {
			logger.WithError(err).Error("failed to publish message")
			return apiTypes.ErrInternal
		}

		return nil
	})

	return nil, err
}

func formDepositIdentifier(request *resources.WithdrawalRequest, chainType types.ChainType) data.DepositIdentifier {
	if chainType == types.ChainTypeZano {
		// only one deposit per transaction
		return data.DepositIdentifier{
			TxHash:  request.Deposit.TxHash,
			ChainId: request.Deposit.ChainId,
		}
	}

	// can be multiple deposits per transaction
	return data.DepositIdentifier{
		TxHash:    request.Deposit.TxHash,
		TxEventId: int(request.Deposit.TxEventIndex),
		ChainId:   request.Deposit.ChainId,
	}
}

func checkExistingDeposit(deposit *data.Deposit, request *resources.WithdrawalRequest, chainType types.ChainType) error {
	if deposit == nil || request == nil {
		return nil
	}

	if !deposit.Reprocessable() {
		return apiTypes.ErrTxAlreadySubmitted
	}
	if chainType != types.ChainTypeZano {
		return nil
	}

	// zano deposits are should be unique by tx hash so event indexes should match
	if deposit.TxEventId != int(request.Deposit.TxEventIndex) {
		return apiTypes.ErrTxAlreadySubmitted
	}

	return nil
}
