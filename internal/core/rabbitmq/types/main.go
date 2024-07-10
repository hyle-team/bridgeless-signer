package types

import (
	"context"

	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DelayExchange           = "delay-exchange"
	DelayQueuePrefix        = "delay-queue"
	DeadLetterExchangeParam = "x-dead-letter-exchange"

	HeadersMatchParam = "x-match"
	HeadersMatchAll   = "all"
	HeadersMatchAny   = "any"

	HeaderDelayKey      = "x-delay"
	HeaderRetryCountKey = "x-retry-count"

	GetDepositQueue       = "get-deposit-queue"
	FormWithdrawalQueue   = "form-withdrawal-queue"
	SignWithdrawalQueue   = "sign-withdrawal-queue"
	SubmitWithdrawalQueue = "submit-withdrawal-queue"
)

var ErrorMaxResendReached = errors.New("max resend count reached")

type Producer interface {
	SendGetDepositRequest(request bridgeTypes.GetDepositRequest) error
	SendFormWithdrawalRequest(request bridgeTypes.FormWithdrawalRequest) error
	SendSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error
	SendSubmitWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error
	DeliveryResender
}

type DeliveryResender interface {
	ResendDelivery(queue string, msg amqp.Delivery) error
}

type Consumer interface {
	Consume(ctx context.Context, queue string) error
}

type DeliveryProcessor interface {
	// ProcessDelivery processes the delivery and returns whether the delivery should be reprocessed,
	// a callback to be called if the reprocessing fails, and an error.
	ProcessDelivery(delivery amqp.Delivery) (reprocessable bool, rprFailCallback func() error, err error)
}
