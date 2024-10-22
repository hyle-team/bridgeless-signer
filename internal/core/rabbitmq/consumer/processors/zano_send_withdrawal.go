package processors

import (
	"encoding/json"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ZanoSendWithdrawalHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewZanoSendWithdrawalHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.DeliveryProcessor {
	return &ZanoSendWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h *ZanoSendWithdrawalHandler) ProcessDelivery(delivery amqp.Delivery) (reprocessable bool, rprFailCallback func() error, err error) {
	var request bridgeTypes.ZanoSignedWithdrawalRequest
	if err = json.Unmarshal(delivery.Body, &request); err != nil {
		return false, nil, errors.Wrap(err, "failed to unmarshal delivery body")
	}

	defer func() {
		if reprocessable {
			rprFailCallback = func() error {
				return errors.Wrap(
					h.processor.SetWithdrawStatusFailed(request.DepositDbId),
					"failed to set withdraw status failed",
				)
			}
		}

	}()

	submitReq, reprocessable, err := h.processor.ProcessZanoSendWithdrawalRequest(request)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process zano send withdrawal request")
	}

	if err = h.publisher.PublishSubmitTransactionRequest(*submitReq); err != nil {
		return true, rprFailCallback, errors.Wrap(err, "failed to send submit withdraw request")
	}

	return false, nil, nil
}
