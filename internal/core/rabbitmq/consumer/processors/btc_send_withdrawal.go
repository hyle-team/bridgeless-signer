package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type BitcoinSendWithdrawalHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewBitcoinSendWithdrawalHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.BatchProcessor[processor.WithdrawalRequest] {
	return &BitcoinSendWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h BitcoinSendWithdrawalHandler) ProcessBatch(batch []processor.WithdrawalRequest) (reprocessable bool, err error) {
	if len(batch) == 0 {
		return false, nil
	}

	reprocessable, err = h.processor.ProcessSendBitcoinWithdrawals(batch...)
	if err != nil {
		return reprocessable, errors.Wrap(err, "failed to process send bitcoin withdrawal request")
	}

	for _, entry := range batch {
		submitTxReq := processor.SubmitTransactionRequest{DepositDbId: entry.DepositDbId}
		if err = h.publisher.PublishSubmitTransactionRequest(submitTxReq); err != nil {
			return true, errors.Wrap(err, "failed to send submit transaction request")
		}
	}

	return false, nil
}

func (h BitcoinSendWithdrawalHandler) ReprocessFailedCallback(batch []processor.WithdrawalRequest) error {
	ids := make([]int64, len(batch))
	for i, req := range batch {
		ids[i] = req.DepositDbId
	}

	return errors.Wrap(
		h.processor.SetWithdrawStatusFailed(ids...),
		"failed to set withdraw status failed",
	)
}
