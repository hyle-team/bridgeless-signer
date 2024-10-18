package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type SubmitTransactionHandler struct {
	processor *processor.Processor
}

func NewSubmitTransactionHandler(
	processor *processor.Processor,
) rabbitTypes.BatchProcessor[bridgeTypes.SubmitTransactionRequest] {
	return &SubmitTransactionHandler{
		processor: processor,
	}
}

func (s SubmitTransactionHandler) ProcessBatch(batch []bridgeTypes.SubmitTransactionRequest) (reprocessable bool, rprFailCallback func(ids ...int64) error, err error) {
	if len(batch) == 0 {
		return false, nil, nil
	}

	defer func() {
		if reprocessable {
			rprFailCallback = func(ids ...int64) error {
				return errors.Wrap(
					s.processor.SetSubmitStatusFailed(ids...),
					"failed to set submit status failed",
				)
			}
		}
	}()

	reprocessable, err = s.processor.ProcessSubmitTransactions(batch...)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process submit transaction request")
	}

	return false, nil, nil
}
