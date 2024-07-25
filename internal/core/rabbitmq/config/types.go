package config

import (
	"runtime"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Connection        *amqp.Connection `fig:"url,required"`
	ConsumerInstances uint             `fig:"consumer_instances,required"`
	ResendParams      ResendParams     `fig:"resend_params,required"`
}

type ResendParams struct {
	Delays        []int32 `fig:"delays,required"`
	MaxRetryCount uint    `fig:"max_retry_count,required"`
}

func (c *Config) Validate() error {
	if len(c.ResendParams.Delays) == 0 {
		return errors.New("delays should not be empty")
	}

	if c.ConsumerInstances == 0 {
		c.ConsumerInstances = uint(runtime.NumCPU())
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
