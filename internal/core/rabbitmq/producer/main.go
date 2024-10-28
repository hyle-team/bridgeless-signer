package producer

import (
	"fmt"
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
		rabbitTypes.EthSignWithdrawalQueue,
		rabbitTypes.BtcSendWithdrawalQueue,
		rabbitTypes.SubmitTransactionQueue,
		rabbitTypes.ZanoSignWithdrawalQueue,
		rabbitTypes.ZanoSendWithdrawalQueue,
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
