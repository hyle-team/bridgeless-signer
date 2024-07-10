package processors

import (
	"encoding/json"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type FormWithdrawalHandler struct {
	processor processor.Processor
	producer  rabbitTypes.Producer
}

func NewFormWithdrawalHandler(
	processor processor.Processor,
	producer rabbitTypes.Producer,
) rabbitTypes.DeliveryProcessor {
	return &FormWithdrawalHandler{
		processor: processor,
		producer:  producer,
	}
}

func (h *FormWithdrawalHandler) ProcessDelivery(delivery amqp.Delivery) (reprocessable bool, rprFailCallback func() error, err error) {
	var request bridgeTypes.FormWithdrawRequest
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

	signReq, reprocessable, err := h.processor.ProcessFormWithdrawRequest(request)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process get deposit request")
	}

	if err = h.producer.SendSignWithdrawRequest(*signReq); err != nil {
		return true, rprFailCallback, errors.Wrap(err, "failed to send form withdraw request")
	}

	return false, nil, nil
}
