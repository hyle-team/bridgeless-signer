package server

import (
	"context"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/ctx"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/requests"
	apiTypes "github.com/hyle-team/bridgeless-signer/internal/core/api/types"
	"github.com/hyle-team/bridgeless-signer/resources"

	"github.com/hyle-team/bridgeless-signer/internal/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (grpcImplementation) CheckWithdrawal(ctxt context.Context, request *resources.CheckWithdrawalRequest) (*resources.CheckWithdrawalResponse, error) {
	var (
		proxies = ctx.Proxies(ctxt)
		db      = ctx.DB(ctxt)
		logger  = ctx.Logger(ctxt)
	)

	wr, err := requests.CheckWithdrawalRequest(request, proxies)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	depositIdentifier := data.DepositIdentifier{
		TxHash:    wr.Deposit.TxHash,
		TxEventId: int(wr.Deposit.TxEventIndex),
		ChainId:   wr.Deposit.ChainId,
	}

	tx, err := db.Get(depositIdentifier)
	if err != nil {
		logger.WithError(err).Error("failed to get deposit")
		return nil, apiTypes.ErrInternal
	}
	if tx == nil {
		return nil, status.Error(codes.NotFound, "deposit not found")
	}

	if tx.Status == resources.WithdrawalStatus_TX_PENDING && tx.WithdrawalTxHash != nil {
		proxy, err := proxies.Proxy(*tx.WithdrawalChainId)
		if err != nil {
			logger.WithError(err).Error("failed to get proxy")
			return nil, apiTypes.ErrInternal // should not happen if the chain is supported, but just in case
		}

		st, err := proxy.GetTransactionStatus(*tx.WithdrawalTxHash)
		if err != nil {
			logger.WithError(err).Error("failed to get tx receipt")
			return nil, apiTypes.ErrInternal
		}

		if st != bridgeTypes.TransactionStatusPending {
			switch st {
			case bridgeTypes.TransactionStatusFailed:
				tx.Status = resources.WithdrawalStatus_TX_FAILED
			case bridgeTypes.TransactionStatusSuccessful:
				tx.Status = resources.WithdrawalStatus_TX_SUCCESSFUL
			}
			// updating in the db
			if err = db.UpdateWithdrawalStatus(tx.Status, tx.Id); err != nil {
				logger.WithError(err).Error("failed to update transaction status")
				return nil, apiTypes.ErrInternal
			}
		}
	}

	return tx.ToStatusResponse(), nil
}
