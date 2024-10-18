package processors

import (
	"encoding/json"

	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SignEthWithdrawalHandler struct {
	processor *processor.Processor
	producer  rabbitTypes.Producer
}

func NewSignEthWithdrawalHandler(
	processor *processor.Processor,
	producer rabbitTypes.Producer,
) rabbitTypes.DeliveryProcessor {
	return &SignEthWithdrawalHandler{
		processor: processor,
		producer:  producer,
	}
}

func (h *SignEthWithdrawalHandler) ProcessDelivery(delivery amqp.Delivery) (reprocessable bool, rprFailCallback func() error, err error) {
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

	submitReq, reprocessable, err := h.processor.ProcessSignEthWithdrawalRequest(request)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process sign withdrawal request")
	}

	if err = h.producer.SendSubmitTransactionRequest(*submitReq); err != nil {
		return true, rprFailCallback, errors.Wrap(err, "failed to send submit withdraw request")
	}

	return false, nil, nil
}
