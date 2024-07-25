package config

import (
	"github.com/spf13/cast"
	"reflect"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Rabbitter interface {
	RabbitMQConfig() Config
}

type configer struct {
	once   comfig.Once
	getter kv.Getter
}

func NewRabbitter(getter kv.Getter) Rabbitter {
	return &configer{
		getter: getter,
	}
}

func (c *configer) RabbitMQConfig() Config {
	return c.once.Do(func() interface{} {
		const rabbitmqConfigKey = "rabbitmq"
		var cfg Config

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, rabbitHooks).
			From(kv.MustGetStringMap(c.getter, rabbitmqConfigKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out rabbitmq config"))
		}

		return cfg
	}).(Config)
}

var rabbitHooks = figure.Hooks{
	"*amqp091.Connection": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case string:
			conn, err := amqp.Dial(v)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to dial amqp")
			}

			return reflect.ValueOf(conn), nil
		default:
			return reflect.Value{}, errors.Errorf("failed to cast %#v of type %T to *amqp.Connection", value, value)
		}
	},
	"[]int32": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case []interface{}:
			var a []int32
			for i, u := range v {
				int32Value, err := cast.ToInt32E(u)
				if err != nil {
					return reflect.Value{}, errors.Wrapf(err, "failed to cast slice element number %d: %#v of type %T into int32", i, value, value)
				}
				a = append(a, int32Value)
			}

			return reflect.ValueOf(a), nil
		default:
			return reflect.Value{}, errors.Errorf("failed to cast %#v of type %T to []int32", value, value)
		}
	},
}
