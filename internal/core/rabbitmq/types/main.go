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

	// Make sure to use values without "x-" prefix for the headers
	// that should be comparable by the headers exchange.
	// In this case, whe use "delay" to route messages to the specific delay queue
	// and "x-retry-count" just to count the number of retries. It's not used for routing.

	HeaderDelayKey      = "delay"
	HeaderRetryCountKey = "x-retry-count"

	GetDepositQueue              = "get-deposit-queue"
	SignWithdrawalQueue          = "sign-withdrawal-queue"
	SubmitBitcoinWithdrawalQueue = "submit-bitcoin-withdrawal-queue"
	SubmitTransactionQueue       = "submit-transaction-queue"
)

var ErrorMaxResendReached = errors.New("max resend count reached")

type Producer interface {
	SendGetDepositRequest(request bridgeTypes.GetDepositRequest) error
	SendSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error
	SendSubmitBitcoinWithdrawalRequest(request bridgeTypes.BitcoinWithdrawalRequest) error
	SendSubmitTransactionRequest(request bridgeTypes.SubmitTransactionRequest) error
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

type Identifiable interface {
	Id() int64
}

type BatchProcessor[T Identifiable] interface {
	// ProcessBatch processes the batch and returns whether the batch should be reprocessed,
	// a callback to be called if the reprocessing fails, and an error.
	ProcessBatch(batch []T) (reprocessable bool, rprFailCallback func(ids ...int64) error, err error)
}
