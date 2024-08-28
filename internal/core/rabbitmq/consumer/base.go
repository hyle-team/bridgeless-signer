package consumer

import (
	"context"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	GetDepositConsumerPrefix       = "get_deposit_consumer"
	FormWithdrawalConsumerPrefix   = "form_withdrawal_consumer"
	SignWithdrawalConsumerPrefix   = "sign_withdrawal_consumer"
	SubmitWithdrawalConsumerPrefix = "submit_withdrawal_consumer"
)

type BaseConsumer struct {
	channel           *amqp.Channel
	name              string
	logger            *logan.Entry
	deliveryProcessor rabbitTypes.DeliveryProcessor
	deliveryResender  rabbitTypes.DeliveryResender
}

func NewBase(
	channel *amqp.Channel,
	name string,
	logger *logan.Entry,
	deliveryProcessor rabbitTypes.DeliveryProcessor,
	deliveryResender rabbitTypes.DeliveryResender,
) rabbitTypes.Consumer {
	return &BaseConsumer{
		channel:           channel,
		name:              name,
		logger:            logger,
		deliveryProcessor: deliveryProcessor,
		deliveryResender:  deliveryResender,
	}
}

func (c *BaseConsumer) Consume(ctx context.Context, queue string) error {
	deliveries, err := c.channel.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channel")
	}

	c.logger.Info("consuming started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("consuming stopped")
			return errors.Wrap(c.channel.Close(), "failed to close channel")
		case delivery, ok := <-deliveries:
			if !ok {
				return nil
			}

			logger := c.logger.WithField("delivery_tag", delivery.DeliveryTag)
			logger.Debug("delivery received")

			reprocessable, callback, err := c.deliveryProcessor.ProcessDelivery(delivery)
			if err == nil {
				ack(logger, delivery)
				continue
			}

			logger.WithError(err).Error("failed to process delivery")
			if !reprocessable {
				logger.Debug("delivery is not reprocessable")
				nack(logger, delivery, false)
				continue
			}

			if err = c.deliveryResender.ResendDelivery(queue, delivery); err == nil {
				logger.Debug("delivery resent")
				ack(logger, delivery)
				continue
			}

			if errors.Is(err, rabbitTypes.ErrorMaxResendReached) {
				logger.Debug(err.Error())
				if callback != nil {
					if err := callback(); err != nil {
						logger.WithError(err).Error("failed to call reprocess fail callback")
					}
				}

				nack(logger, delivery, false)
			} else {
				logger.WithError(err).Error("failed to resend delivery")
				nack(logger, delivery, true)
			}
		}
	}
}
