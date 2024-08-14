package api

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyle-team/bridgeless-signer/docs"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/config"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/ignite/cli/ignite/pkg/openapiconsole"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var _ types.ServiceServer = &Server{}

type ServiceHandler interface {
	SubmitWithdrawal(ctx context.Context, request *types.WithdrawalRequest) error
	CheckWithdrawal(ctx context.Context, request *types.CheckWithdrawalRequest) (*types.CheckWithdrawalResponse, error)
}

// Server is a GRPC and HTTP gateway application server.
type Server struct {
	listener    net.Listener
	gatewayAddr string
	handler     ServiceHandler
}

// NewServer creates a new GRPC server.
func NewServer(
	listener net.Listener,
	gatewayCfg config.RESTGatewayConfig,
	handler ServiceHandler,
) *Server {
	return &Server{
		listener:    listener,
		gatewayAddr: gatewayCfg.Address,
		handler:     handler,
	}
}

// RunGRPC starts the GRPC server.
func (s *Server) RunGRPC(ctx context.Context) error {
	grpcServer := grpc.NewServer()
	types.RegisterServiceServer(grpcServer, s)

	// graceful shutdown
	go func() { <-ctx.Done(); grpcServer.Stop() }()
	return grpcServer.Serve(s.listener)
}

// RunRESTGateway starts the REST gateway server.
func (s *Server) RunRESTGateway(ctx context.Context) error {
	grpcGatewayRouter := runtime.NewServeMux()
	httpRouter := http.NewServeMux()

	if err := types.RegisterServiceHandlerServer(context.Background(), grpcGatewayRouter, s); err != nil {
		return errors.Wrap(err, "failed to register service handler")
	}

	httpRouter.Handle("/static/service.swagger.json", http.FileServer(http.FS(docs.Docs)))
	httpRouter.HandleFunc("/api", openapiconsole.Handler("TSS service", "/static/service.swagger.json"))
	httpRouter.Handle("/", grpcGatewayRouter)

	srv := &http.Server{Addr: s.gatewayAddr, Handler: httpRouter}

	// graceful shutdown
	go func() { <-ctx.Done(); srv.Shutdown(ctx) }()
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "failed to listen and serve")
	}

	return nil
}

func (s *Server) SubmitWithdrawal(ctx context.Context, request *types.WithdrawalRequest) (*types.Empty, error) {
	return &types.Empty{}, s.handler.SubmitWithdrawal(ctx, request)
}

func (s *Server) CheckWithdrawal(ctx context.Context, request *types.CheckWithdrawalRequest) (*types.CheckWithdrawalResponse, error) {
	return s.handler.CheckWithdrawal(ctx, request)
}
