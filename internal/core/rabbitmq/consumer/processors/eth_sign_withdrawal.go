package processors

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/processor"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
)

type EthereumSignWithdrawalHandler struct {
	processor *processor.Processor
	publisher rabbitTypes.Producer
}

func NewEthereumSignWithdrawalHandler(
	processor *processor.Processor,
	publisher rabbitTypes.Producer,
) rabbitTypes.RequestProcessor[processor.WithdrawalRequest] {
	return &EthereumSignWithdrawalHandler{
		processor: processor,
		publisher: publisher,
	}
}

func (h EthereumSignWithdrawalHandler) ProcessRequest(request processor.WithdrawalRequest) (reprocessable bool, err error) {
	submitReq, reprocessable, err := h.processor.ProcessEthSignWithdrawalRequest(request)
	if err != nil {
		return reprocessable, errors.Wrap(err, "failed to process eth sign withdrawal request")
	}

	if err = h.publisher.PublishSubmitTransactionRequest(*submitReq); err != nil {
		return true, errors.Wrap(err, "failed to send submit withdraw request")
	}

	return false, nil
}

func (h EthereumSignWithdrawalHandler) ReprocessFailedCallback(request processor.WithdrawalRequest) error {
	return errors.Wrap(
		h.processor.SetWithdrawStatusFailed(request.DepositDbId),
		"failed to set withdraw status failed",
	)
}
