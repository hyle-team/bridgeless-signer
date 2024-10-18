package ctx

import (
	"context"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	dbKey ctxKey = iota
	logKey
	proxiesKey
	producerKey
)

func DBProvider(q data.DepositsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbKey, q)
	}
}

// DB always returns unique connection
func DB(ctx context.Context) data.DepositsQ {
	return ctx.Value(dbKey).(data.DepositsQ).New()
}

func LoggerProvider(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logKey, entry)
	}
}

func Logger(ctx context.Context) *logan.Entry {
	return ctx.Value(logKey).(*logan.Entry)
}

func ProxiesProvider(proxies bridgeTypes.ProxiesRepository) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, proxiesKey, proxies)
	}
}

func Proxies(ctx context.Context) bridgeTypes.ProxiesRepository {
	return ctx.Value(proxiesKey).(bridgeTypes.ProxiesRepository)
}

func ProducerProvider(producer rabbitTypes.Publisher) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, producerKey, producer)
	}
}

func Producer(ctx context.Context) rabbitTypes.Publisher {
	return ctx.Value(producerKey).(rabbitTypes.Publisher)
}
