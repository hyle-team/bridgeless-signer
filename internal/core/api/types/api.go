package types

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal           = status.Error(codes.Internal, "internal error")
	ErrTxAlreadySubmitted = status.Error(codes.AlreadyExists, "transaction already submitted")
	ErrInvalidOriginTxId  = errors.New("invalid origin tx id")
)

type Server interface {
	RunGRPC(ctx context.Context) error
	RunHTTP(ctx context.Context) error
}

func GRPCContextExtenderInterceptor(extenders ...func(context.Context) context.Context) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		for _, extend := range extenders {
			ctx = extend(ctx)
		}

		return handler(ctx, req)
	}
}
