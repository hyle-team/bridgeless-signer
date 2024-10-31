package types

import (
	"context"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/processor"

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

	GetDepositQueue         = "get-deposit-queue"
	EthSignWithdrawalQueue  = "eth-sign-withdrawal-queue"
	BtcSendWithdrawalQueue  = "btc-send-withdrawal-queue"
	ZanoSignWithdrawalQueue = "zano-sign-withdrawal-queue"
	ZanoSendWithdrawalQueue = "zano-send-withdrawal-queue"
	SubmitTransactionQueue  = "submit-transaction-queue"
)

var (
	ErrMaxResendReached = errors.New("max resend count reached")
	ErrConnectionClosed = errors.New("RabbitMQ connection was closed")
)

type Producer interface {
	PublishGetDepositRequest(request bridgeTypes.GetDepositRequest) error
	PublishSubmitTransactionRequest(request bridgeTypes.SubmitTransactionRequest) error

	PublishEthereumSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error

	PublishBitcoinSendWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error

	PublishZanoSignWithdrawalRequest(request bridgeTypes.WithdrawalRequest) error
	PublishZanoSendWithdrawalRequest(request bridgeTypes.ZanoSignedWithdrawalRequest) error
	DeliveryResender
}

type DeliveryResender interface {
	ResendDelivery(queue string, msg amqp.Delivery) error
}

type Consumer interface {
	Name() string
	Consume(ctx context.Context, queue string) error
}
