package processors

import (
	"encoding/json"
	"fmt"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type GetDepositHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewGetDepositHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.DeliveryProcessor {
	return &GetDepositHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h *GetDepositHandler) ProcessDelivery(delivery amqp.Delivery) (reprocessable bool, rprFailCallback func() error, err error) {
	var request processor.GetDepositRequest
	if err = json.Unmarshal(delivery.Body, &request); err != nil {
		return false, nil, errors.Wrap(err, "failed to unmarshal delivery body")
	}

	rprFailCallback = func() error {
		return errors.Wrap(
			h.processor.SetWithdrawStatusFailed(request.DepositDbId),
			"failed to set withdraw status failed",
		)
	}

	withdrawReq, reprocessable, err := h.processor.ProcessGetDepositRequest(request)
	if err != nil {
		return reprocessable, rprFailCallback, errors.Wrap(err, "failed to process get deposit request")
	}

	reprocessable = true
	switch withdrawReq.Destination {
	case bridgeTypes.ChainTypeEVM:
		err = h.publisher.PublishEthereumSignWithdrawalRequest(*withdrawReq)
	case bridgeTypes.ChainTypeBitcoin:
		err = h.publisher.PublishBitcoinSendWithdrawalRequest(*withdrawReq)
	case bridgeTypes.ChainTypeZano:
		err = h.publisher.PublishZanoSignWithdrawalRequest(*withdrawReq)
	default:
		err = errors.New(fmt.Sprintf("invalid destination type: %v", withdrawReq.Destination))
		reprocessable = false
	}

	return reprocessable, rprFailCallback, errors.Wrap(err, "failed to send deposit processing request")
}
