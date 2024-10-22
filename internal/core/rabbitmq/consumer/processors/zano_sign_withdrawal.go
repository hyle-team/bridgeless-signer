package processors

import (
	"encoding/json"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ZanoSignWithdrawalHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewZanoSignWithdrawalHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.DeliveryProcessor {
	return &ZanoSignWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h *ZanoSignWithdrawalHandler) ProcessDelivery(delivery amqp.Delivery) (reprocessable bool, rprFailCallback func() error, err error) {
	var request bridgeTypes.WithdrawalRequest
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

	signedWithdrawReq, reprocessable, err := h.processor.ProcessZanoSignWithdrawalRequest(request)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process zano sign withdrawal request")
	}

	if err = h.publisher.PublishZanoSendWithdrawalRequest(*signedWithdrawReq); err != nil {
		return true, rprFailCallback, errors.Wrap(err, "failed to send zano send withdraw request")
	}

	return false, nil, nil
}
