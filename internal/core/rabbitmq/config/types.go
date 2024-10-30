package config

import (
	"runtime"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Connection            *amqp.Connection   `fig:"url,required"`
	BaseConsumerInstances uint               `fig:"base_consumer_instances,required"`
	ResendParams          ResendParams       `fig:"resend_params,required"`
	TxSubmitterOpts       BatchConsumingOpts `fig:"tx_submitter,required"`
	BitcoinSubmitterOpts  BatchConsumingOpts `fig:"bitcoin_submitter,required"`
}

type BatchConsumingOpts struct {
	MaxSize int           `fig:"max_size,required"`
	Period  time.Duration `fig:"period,required"`
}

type ResendParams struct {
	Delays        []int32 `fig:"delays,required"`
	MaxRetryCount uint    `fig:"max_retry_count,required"`
}

func (c *Config) Validate() error {
	if len(c.ResendParams.Delays) == 0 {
		return errors.New("delays should not be empty")
	}

	if c.BaseConsumerInstances == 0 {
		c.BaseConsumerInstances = uint(runtime.NumCPU())
	}

	if c.ResendParams.MaxRetryCount == 0 {
		return errors.New("max_retry_count should be greater than 0")
	}

	return nil
}

func (c *Config) NewChannel() *amqp.Channel {
	ch, err := c.Connection.Channel()
	if err != nil {
		panic(errors.Wrap(err, "failed to open channel"))
	}

	return ch
}
