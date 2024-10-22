package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
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
) rabbitTypes.BatchProcessor[bridgeTypes.WithdrawalRequest] {
	return &BitcoinSendWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h *BitcoinSendWithdrawalHandler) ProcessBatch(batch []bridgeTypes.WithdrawalRequest) (reprocessable bool, rprFailCallback func(ids ...int64) error, err error) {
	if len(batch) == 0 {
		return false, nil, nil
	}

	defer func() {
		if reprocessable {
			rprFailCallback = func(ids ...int64) error {
				return errors.Wrap(
					h.processor.SetWithdrawStatusFailed(ids...),
					"failed to set withdraw status failed",
				)
			}
		}

	}()

	reprocessable, err = h.processor.ProcessSendBitcoinWithdrawals(batch...)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process send bitcoin withdrawal request")
	}

	for _, entry := range batch {
		submitTxReq := bridgeTypes.SubmitTransactionRequest{DepositDbId: entry.DepositDbId}
		if err = h.publisher.PublishSubmitTransactionRequest(submitTxReq); err != nil {
			return true, rprFailCallback, errors.Wrap(err, "failed to send submit transaction request")
		}
	}

	return false, nil, nil
}
