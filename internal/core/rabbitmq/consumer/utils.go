package consumer

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/logan/v3"
)

func ack(logger *logan.Entry, delivery amqp.Delivery) {
	if err := delivery.Ack(false); err != nil {
		logger.WithError(err).Error("failed to ack delivery")
	} else {
		logger.Debug("delivery acked")
	}
}

func nack(logger *logan.Entry, delivery amqp.Delivery, requeue bool) {
	if err := delivery.Nack(false, requeue); err != nil {
		logger.WithError(err).Error("failed to nack delivery")
	} else {
		logger.Debug("delivery nacked")
	}
}

func GetName(prefix string, index uint) string {
	return fmt.Sprintf("%s_%d", prefix, index)
}
