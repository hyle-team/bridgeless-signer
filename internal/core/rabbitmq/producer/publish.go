package producer

import (
	"encoding/json"
	"fmt"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *Producer) PublishGetDepositRequest(request bridgeTypes.GetDepositRequest) error {
	return p.publish(rabbitTypes.GetDepositQueue, request)
}

func (p *Producer) PublishEthereumSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error {
	return p.publish(rabbitTypes.EthSignWithdrawalQueue, request)
}

func (p *Producer) PublishBitcoinSendWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error {
	return p.publish(rabbitTypes.BtcSendWithdrawalQueue, request)
}

func (p *Producer) PublishZanoSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error {
	return p.publish(rabbitTypes.ZanoSignWithdrawalQueue, request)
}

func (p *Producer) PublishZanoSendWithdrawalRequest(request bridgeTypes.ZanoSignedWithdrawalRequest) error {
	return p.publish(rabbitTypes.ZanoSendWithdrawalQueue, request)
}

func (p *Producer) PublishSubmitTransactionRequest(request bridgeTypes.SubmitTransactionRequest) error {
	return p.publish(rabbitTypes.SubmitTransactionQueue, request)
}

func (p *Producer) publish(queue string, marshable any) error {
	raw, err := json.Marshal(marshable)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to marshal message %T", marshable))
	}

	return p.channel.Publish("", queue, false, false, persistentPublishing(raw, nil))
}

func persistentPublishing(body []byte, headers amqp.Table) amqp.Publishing {
	return amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Body:         body,
		Headers:      headers,
	}
}
