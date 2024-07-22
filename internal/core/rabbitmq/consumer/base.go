package consumer

import (
	"context"
	"fmt"

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

type Consumer struct {
	channel           *amqp.Channel
	name              string
	logger            *logan.Entry
	deliveryProcessor rabbitTypes.DeliveryProcessor
	deliveryResender  rabbitTypes.DeliveryResender
}

func New(
	channel *amqp.Channel,
	name string,
	logger *logan.Entry,
	deliveryProcessor rabbitTypes.DeliveryProcessor,
	deliveryResender rabbitTypes.DeliveryResender,
) rabbitTypes.Consumer {
	return &Consumer{
		channel:           channel,
		name:              name,
		logger:            logger,
		deliveryProcessor: deliveryProcessor,
		deliveryResender:  deliveryResender,
	}
}

func (c *Consumer) Consume(ctx context.Context, queue string) error {
	deliveries, err := c.channel.Consume(
		queue, c.name, false, false, false, false, nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channel")
	}

	c.logger.Info("consuming started")

	for {
		select {
		case <-ctx.Done():
			return errors.Wrap(c.channel.Close(), "failed to close channel")
		case delivery, ok := <-deliveries:
			if !ok {
				return nil
			}

			logger := c.logger.WithField("delivery_tag", delivery.DeliveryTag)
			logger.Debug("delivery received")

			reprocessable, callback, err := c.deliveryProcessor.ProcessDelivery(delivery)
			if err == nil {
				c.ack(logger, delivery)
				continue
			}

			logger.WithError(err).Error("failed to process delivery")
			if !reprocessable {
				logger.Debug("delivery is not reprocessable")
				c.nack(logger, delivery, false)
				continue
			}

			err = c.deliveryResender.ResendDelivery(queue, delivery)
			if err == nil {
				logger.Debug("delivery resent")
				c.ack(logger, delivery)
				continue
			}

			logger.WithError(err).Error("failed to resend delivery")
			if errors.Is(err, rabbitTypes.ErrorMaxResendReached) {
				if callback != nil {
					if err := callback(); err != nil {
						logger.WithError(err).Error("failed to call reprocess fail callback")
					}
				}

				c.nack(logger, delivery, false)
			} else {
				c.nack(logger, delivery, true)
			}
		}
	}
}

func (c *Consumer) ack(logger *logan.Entry, delivery amqp.Delivery) {
	if err := delivery.Ack(false); err != nil {
		logger.WithError(err).Error("failed to ack delivery")
	} else {
		logger.Debug("delivery acked")
	}
}

func (c *Consumer) nack(logger *logan.Entry, delivery amqp.Delivery, requeue bool) {
	if err := delivery.Nack(false, requeue); err != nil {
		logger.WithError(err).Error("failed to nack delivery")
	} else {
		logger.Debug("delivery nacked")
	}
}

func GetName(prefix string, index uint) string {
	return fmt.Sprintf("%s_%d", prefix, index)
}
