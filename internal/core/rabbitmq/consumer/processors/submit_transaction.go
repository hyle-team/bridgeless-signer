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

func (s SubmitTransactionHandler) ProcessBatch(batch []processor.SubmitTransactionRequest) (reprocessable bool, err error) {
	if len(batch) == 0 {
		return false, nil
	}

	reprocessable, err = s.processor.ProcessSubmitTransactions(batch...)

	return reprocessable, errors.Wrap(err, "failed to process submit transaction request")
}

func (s SubmitTransactionHandler) ReprocessFailedCallback(batch []processor.SubmitTransactionRequest) error {
	ids := make([]int64, len(batch))
	for i, req := range batch {
		ids[i] = req.DepositDbId
	}

	return errors.Wrap(
		s.processor.SetSubmitStatusFailed(ids...),
		"failed to set submit status failed",
	)
}
