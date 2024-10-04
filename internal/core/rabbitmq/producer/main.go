package producer

import (
	"encoding/json"
	"fmt"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/config"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp.Channel

	maxRetryCount uint
	delays        []int32
}

// New creates a new Producer instance.
// It ensures that the exchange and queues are created.
func New(ch *amqp.Channel, resendParams config.ResendParams) (rabbitTypes.Producer, error) {

	// Queues is bound to the default exchange
	var consumerQueues = []string{
		rabbitTypes.GetDepositQueue,
		rabbitTypes.SignWithdrawalQueue,
		rabbitTypes.SubmitBitcoinWithdrawalQueue,
		rabbitTypes.SubmitTransactionQueue,
	}

	for _, queue := range consumerQueues {
		_, err := ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to declare queue %s", queue))
		}
	}

	// Delay exchange is used to route messages to the delay queues
	if err := ch.ExchangeDeclare(
		rabbitTypes.DelayExchange, amqp.ExchangeHeaders,
		true, false, false, false, nil,
	); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to declare exchange %s", rabbitTypes.DelayExchange))
	}

	// Declaring delay queues and bind them to the delay exchange
	for _, delay := range resendParams.Delays {
		qName := getDelayQueueName(rabbitTypes.DelayQueuePrefix, delay)
		_, err := ch.QueueDeclare(qName, true, false, false, false, delayQueueArgs(delay))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to declare delay queue %s", qName))
		}

		// Bind the delay queue to the delay exchange
		if err = ch.QueueBind(qName, "", rabbitTypes.DelayExchange, false, delayQueueBindArgs(delay)); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to bind queue %s to exchange %s", qName, rabbitTypes.DelayExchange))
		}
	}

	return &Producer{
		channel:       ch,
		maxRetryCount: resendParams.MaxRetryCount,
		delays:        resendParams.Delays,
	}, nil
}

func (p *Producer) publish(queue string, marshable any) error {
	raw, err := json.Marshal(marshable)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to marshal %T", marshable))
	}

	return p.channel.Publish("", queue, false, false, persistentPublishing(raw, nil))
}

func (p *Producer) SendGetDepositRequest(request bridgeTypes.GetDepositRequest) error {
	return p.publish(rabbitTypes.GetDepositQueue, request)
}

func (p *Producer) SendSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error {
	return p.publish(rabbitTypes.SignWithdrawalQueue, request)
}

func (p *Producer) SendSubmitBitcoinWithdrawalRequest(request bridgeTypes.BitcoinWithdrawalRequest) error {
	return p.publish(rabbitTypes.SubmitBitcoinWithdrawalQueue, request)
}

func (p *Producer) SendSubmitTransactionRequest(request bridgeTypes.SubmitTransactionRequest) error {
	return p.publish(rabbitTypes.SubmitTransactionQueue, request)
}

func (p *Producer) ResendDelivery(queue string, msg amqp.Delivery) error {
	retryCount := p.getCurrentRetryNumber(msg)
	if retryCount >= int32(p.maxRetryCount) {
		return rabbitTypes.ErrorMaxResendReached
	}

	retryCount++
	delay := p.getDelay(retryCount)

	return p.channel.Publish(rabbitTypes.DelayExchange, queue, false, false, p.formResendMsg(msg, retryCount, delay))
}

func (p *Producer) getCurrentRetryNumber(msg amqp.Delivery) int32 {
	retryRaw, found := msg.Headers[rabbitTypes.HeaderRetryCountKey]
	if !found {
		return 0
	}

	retry, ok := retryRaw.(int32)
	if !ok {
		return 0
	}

	return retry
}

func (p *Producer) getDelay(retry int32) int32 {
	if retry != 0 {
		// Decrement the retry count to get the delay index
		retry--
	}

	if retry >= int32(p.maxRetryCount) {
		return 0
	}

	if int(retry) >= len(p.delays) {
		return p.delays[len(p.delays)-1]
	}

	return p.delays[retry]
}

func (p *Producer) formResendMsg(msg amqp.Delivery, retryCount int32, delay int32) amqp.Publishing {
	return persistentPublishing(msg.Body,
		amqp.Table{
			rabbitTypes.HeaderRetryCountKey: retryCount,
			rabbitTypes.HeaderDelayKey:      delay,
		},
	)
}

func persistentPublishing(body []byte, headers amqp.Table) amqp.Publishing {
	return amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Body:         body,
		Headers:      headers,
	}
}

func getDelayQueueName(qPrefix string, delay int32) string {
	return fmt.Sprintf("%s-%d", qPrefix, delay)
}

func delayQueueArgs(delay int32) amqp.Table {
	return map[string]interface{}{
		// Set the time in milliseconds for which the message will be stored in the queue.
		// After this time, the message will be routed to the dead-letter exchange.
		amqp.QueueMessageTTLArg: delay,
		// The exchange to which the message will be routed after the TTL expires.
		// Using an empty string means that the message will be routed to the default exchange.
		rabbitTypes.DeadLetterExchangeParam: "",
	}
}

func delayQueueBindArgs(delay int32) amqp.Table {
	return map[string]interface{}{
		rabbitTypes.HeadersMatchParam: rabbitTypes.HeadersMatchAll,
		rabbitTypes.HeaderDelayKey:    delay,
	}
}
