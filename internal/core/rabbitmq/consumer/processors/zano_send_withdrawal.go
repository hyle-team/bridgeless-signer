package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type ZanoSendWithdrawalHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewZanoSendWithdrawalHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.RequestProcessor[processor.ZanoSignedWithdrawalRequest] {
	return &ZanoSendWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h ZanoSendWithdrawalHandler) ProcessRequest(request processor.ZanoSignedWithdrawalRequest) (reprocessable bool, err error) {
	submitReq, reprocessable, err := h.processor.ProcessZanoSendWithdrawalRequest(request)
	if err != nil {
		return reprocessable, errors.Wrap(err, "failed to process zano send withdrawal request")
	}

	if err = h.publisher.PublishSubmitTransactionRequest(*submitReq); err != nil {
		return true, errors.Wrap(err, "failed to send submit withdraw request")
	}

	return false, nil
}

func (h ZanoSendWithdrawalHandler) ReprocessFailedCallback(request processor.ZanoSignedWithdrawalRequest) error {
	return errors.Wrap(
		h.processor.SetWithdrawStatusFailed(request.DepositDbId),
		"failed to set withdraw status failed",
	)
}
