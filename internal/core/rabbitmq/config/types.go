package config

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Connection   *amqp.Connection `fig:"url,required"`
	ResendParams ResendParams     `fig:"resend_params,required"`
}

type ResendParams struct {
	Delays        []int64 `fig:"delays,required"`
	MaxRetryCount uint    `fig:"max_retry_count,required"`
}

func (c Config) Validate() error {
	if len(c.ResendParams.Delays) == 0 {
		return errors.New("delays should not be empty")
	}

	if c.ResendParams.MaxRetryCount == 0 {
		return errors.New("max_retry_count should be greater than 0")
	}

	return nil
}
