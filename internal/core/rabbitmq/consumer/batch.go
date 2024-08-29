package consumer

import (
	"context"
	"encoding/json"
	"github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/config"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/logan/v3"
	"time"
)

const (
	SubmitTransactionConsumerPrefix = "submit_transaction_consumer"
)

type amqpParsedEntry[T rabbitTypes.Identifiable] struct {
	Delivery amqp.Delivery
	Entry    T
}

type BatchConsumer[T rabbitTypes.Identifiable] struct {
	channel *amqp.Channel
	name    string
	logger  *logan.Entry

	deliveryResender rabbitTypes.DeliveryResender
	batch            []amqpParsedEntry[T]
	batchProcessor   rabbitTypes.BatchProcessor[T]

	opts config.BatchConsumingOpts
}

func NewBatch[T rabbitTypes.Identifiable](
	channel *amqp.Channel,
	name string,
	logger *logan.Entry,
	batchProcessor rabbitTypes.BatchProcessor[T],
	deliveryResender rabbitTypes.DeliveryResender,
	opts config.BatchConsumingOpts,
) rabbitTypes.Consumer {
	return &BatchConsumer[T]{
		channel: channel,
		name:    name,
		logger:  logger,

		batchProcessor:   batchProcessor,
		deliveryResender: deliveryResender,

		batch: make([]amqpParsedEntry[T], 0, opts.MaxSize),
		opts:  opts,
	}
}

func (c *BatchConsumer[T]) Consume(ctx context.Context, queue string) error {
	deliveries, err := c.channel.Consume(queue, c.name, false, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get consumer channel")
	}

	c.logger.Info("consuming started")

	ticker := time.NewTicker(c.opts.Period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if len(c.batch) != 0 {
				// return ack-ed deliveries to the sender in case of context cancellation
				c.logger.Info("resending ack-ed deliveries")
				for _, entry := range c.batch {
					_ = c.deliveryResender.ResendDelivery(queue, entry.Delivery)
				}
			}

			c.logger.Info("consuming stopped")
			c.batch = c.batch[:0]

			return errors.Wrap(c.channel.Close(), "failed to close channel")
		case delivery, ok := <-deliveries:
			if !ok {
				return nil
			}

			logger := c.logger.WithField("delivery_tag", delivery.DeliveryTag)
			logger.Debug("delivery received")

			var msg T
			if err = json.Unmarshal(delivery.Body, &msg); err != nil {
				logger.WithError(err).Error("failed to unmarshal delivery body")
				nack(logger, delivery, false)
				continue
			}

			c.batch = append(c.batch, amqpParsedEntry[T]{Delivery: delivery, Entry: msg})
			ack(logger, delivery)

			if len(c.batch) == c.opts.MaxSize {
				c.processBatch(queue)
			}
		case <-ticker.C:
			c.processBatch(queue)
		}
	}
}

func (c *BatchConsumer[T]) processBatch(queue string) {
	if len(c.batch) == 0 {
		return
	}

	// emptying the batch
	defer func() { c.batch = c.batch[:0] }()

	logger := c.logger.WithField("batch_size", len(c.batch))
	logger.Debug("processing batch")

	entryBatch := make([]T, len(c.batch))
	for i, entry := range c.batch {
		entryBatch[i] = entry.Entry
	}

	reprocessable, callback, err := c.batchProcessor.ProcessBatch(entryBatch)
	if err == nil {
		logger.Debug("batch processed")
		return
	}

	logger.WithError(err).Error("failed to process batch")
	if !reprocessable {
		logger.Debug("batch is not reprocessable")
		return
	}

	var callbackIds []int64
	for _, entry := range c.batch {
		// shadowing original logger
		logger := logger.WithField("delivery_tag", entry.Delivery.DeliveryTag)

		err = c.deliveryResender.ResendDelivery(queue, entry.Delivery)
		if err == nil {
			logger.Debug("delivery resent")
			continue
		}

		logger.WithError(err).Error("failed to resend delivery")
		if errors.Is(err, rabbitTypes.ErrorMaxResendReached) {
			callbackIds = append(callbackIds, entry.Entry.Id())
		}
	}

	if callback != nil && len(callbackIds) > 0 {
		if err = callback(callbackIds...); err != nil {
			logger.WithField("callback_ids", callbackIds).WithError(err).Error("failed to set batch status failed")
		}
	}
}