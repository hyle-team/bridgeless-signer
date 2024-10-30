package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type ZanoSignWithdrawalHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewZanoSignWithdrawalHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.RequestProcessor[processor.WithdrawalRequest] {
	return &ZanoSignWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h ZanoSignWithdrawalHandler) ProcessRequest(request processor.WithdrawalRequest) (reprocessable bool, err error) {
	signedWithdrawReq, reprocessable, err := h.processor.ProcessZanoSignWithdrawalRequest(request)
	if err != nil {
		return reprocessable, errors.Wrap(err, "failed to process zano sign withdrawal request")
	}

	if err = h.publisher.PublishZanoSendWithdrawalRequest(*signedWithdrawReq); err != nil {
		return true, errors.Wrap(err, "failed to send zano send withdraw request")
	}

	return false, nil
}

func (h ZanoSignWithdrawalHandler) ReprocessFailedCallback(request processor.WithdrawalRequest) error {
	return errors.Wrap(
		h.processor.SetWithdrawStatusFailed(request.DepositDbId),
		"failed to set withdraw status failed",
	)
}
