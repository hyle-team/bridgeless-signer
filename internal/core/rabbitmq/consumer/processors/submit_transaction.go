package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type SubmitTransactionHandler struct {
	processor *processor.Processor
}

func NewSubmitTransactionHandler(
	processor *processor.Processor,
) rabbitTypes.BatchProcessor[processor.SubmitTransactionRequest] {
	return &SubmitTransactionHandler{
		processor: processor,
	}
}

func (s SubmitTransactionHandler) ProcessBatch(batch []processor.SubmitTransactionRequest) (reprocessable bool, rprFailCallback func(ids ...int64) error, err error) {
	if len(batch) == 0 {
		return false, nil, nil
	}

	rprFailCallback = func(ids ...int64) error {
		return errors.Wrap(
			s.processor.SetSubmitStatusFailed(ids...),
			"failed to set submit status failed",
		)
	}

	reprocessable, err = s.processor.ProcessSubmitTransactions(batch...)

	return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process submit transaction request")
}
