package processors

import (
	"fmt"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type GetDepositHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewGetDepositHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.RequestProcessor[processor.GetDepositRequest] {
	return &GetDepositHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h GetDepositHandler) ProcessRequest(request processor.GetDepositRequest) (reprocessable bool, err error) {
	withdrawReq, reprocessable, err := h.processor.ProcessGetDepositRequest(request)
	if err != nil {
		return reprocessable, errors.Wrap(err, "failed to process get deposit request")
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

	return reprocessable, errors.Wrap(err, "failed to send deposit processing request")
}

func (h GetDepositHandler) ReprocessFailedCallback(request processor.GetDepositRequest) error {
	return errors.Wrap(
		h.processor.SetWithdrawStatusFailed(request.DepositDbId),
		"failed to set withdraw status failed",
	)
}
