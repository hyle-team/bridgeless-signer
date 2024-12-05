package consumer

import (
	"context"
	"encoding/json"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	GetDepositConsumerPrefix         = "get_deposit_consumer"
	EthSignWithdrawalConsumerPrefix  = "eth_sign_withdrawal_consumer"
	ZanoSignWithdrawalConsumerPrefix = "zano_sign_withdrawal_consumer"
	ZanoSendWithdrawalConsumerPrefix = "zano_send_withdrawal_consumer"
)

type BaseConsumer[T rabbitTypes.Identifiable] struct {
	channel           *amqp.Channel
	name              string
	logger            *logan.Entry
	deliveryProcessor rabbitTypes.RequestProcessor[T]
	deliveryResender  rabbitTypes.DeliveryResender
}

func NewBase[T rabbitTypes.Identifiable](
	channel *amqp.Channel,
	name string,
	logger *logan.Entry,
	deliveryProcessor rabbitTypes.RequestProcessor[T],
	deliveryResender rabbitTypes.DeliveryResender,
) rabbitTypes.Consumer {
	return &BaseConsumer[T]{
		channel:           channel,
		name:              name,
		logger:            logger,
		deliveryProcessor: deliveryProcessor,
		deliveryResender:  deliveryResender,
	}
}

func (c *BaseConsumer[T]) Consume(ctx context.Context, queue string) error {
	deliveries, err := c.channel.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channel")
	}

	c.logger.Info("consuming started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consuming stopped: context canceled")

			return errors.Wrap(c.channel.Close(), "failed to close channel")
		case delivery, ok := <-deliveries:
			if !ok {
				c.logger.Info("consuming stopped: delivery channel closed")

				return nil
			}

			logger := c.logger.WithField("delivery_tag", delivery.DeliveryTag)
			logger.Debug("delivery received")

			var request T
			if err = json.Unmarshal(delivery.Body, &request); err != nil {
				logger.WithError(err).Error("failed to unmarshal delivery body")
				nack(logger, delivery, false)
				continue
			}

			reprocessable, err := c.deliveryProcessor.ProcessRequest(request)
			if err == nil {
				ack(logger, delivery)
				continue
			}

			nack(logger, delivery, false)
			logger.WithError(err).Error("failed to process request")
			if !reprocessable {
				logger.Debug("request is not reprocessable")
				continue
			}

			if err = c.deliveryResender.ResendDelivery(queue, delivery); err == nil {
				logger.Debug("delivery resent")
				continue
			}
			if errors.Is(err, rabbitTypes.ErrMaxResendReached) {
				logger.Debug(rabbitTypes.ErrMaxResendReached)
			} else {
				logger.WithError(err).Error("failed to resend delivery")
			}

			if err = c.deliveryProcessor.ReprocessFailedCallback(request); err != nil {
				logger.WithError(err).Error("failed to execute failed reprocessing callback")
			}
		}
	}
}

func (c *BaseConsumer[T]) Name() string {
	return c.name
}
