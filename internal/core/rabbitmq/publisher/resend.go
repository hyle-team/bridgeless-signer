package publisher

import (
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (p *Publisher) ResendDelivery(queue string, msg amqp.Delivery) error {
	retryCount := p.getCurrentRetryNumber(msg)
	if retryCount >= int32(p.maxRetryCount) {
		return rabbitTypes.ErrorMaxResendReached
	}

	retryCount++
	delay := p.getDelay(retryCount)

	return p.channel.Publish(rabbitTypes.DelayExchange, queue, false, false, p.formResendMsg(msg, retryCount, delay))
}

func (p *Publisher) getCurrentRetryNumber(msg amqp.Delivery) int32 {
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

func (p *Publisher) getDelay(retry int32) int32 {
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

func (p *Publisher) formResendMsg(msg amqp.Delivery, retryCount int32, delay int32) amqp.Publishing {
	return persistentPublishing(msg.Body,
		amqp.Table{
			rabbitTypes.HeaderRetryCountKey: retryCount,
			rabbitTypes.HeaderDelayKey:      delay,
		},
	)
}
