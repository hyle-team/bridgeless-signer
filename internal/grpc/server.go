package grpc

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyle-team/bridgeless-signer/docs"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/ignite/cli/ignite/pkg/openapiconsole"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var _ types.ServiceServer = &Server{}

type ServiceHandler interface {
	SubmitWithdraw(ctx context.Context, request *types.WithdrawRequest) error
	CheckWithdraw(ctx context.Context, request *types.CheckWithdrawRequest) (*types.CheckWithdrawResponse, error)
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
	gatewayCfg RESTGatewayConfig,
	handler ServiceHandler,
) *Server {
	return &Server{
		listener:    listener,
		gatewayAddr: gatewayCfg.Address,
		handler:     handler,
	}
}

// RunGRPC starts the GRPC server.
func (s *Server) RunGRPC(_ context.Context) error {
	grpcServer := grpc.NewServer()
	types.RegisterServiceServer(grpcServer, s)
	return grpcServer.Serve(s.listener)
}

// RunRESTGateway starts the REST gateway server.
func (s *Server) RunRESTGateway(ctx context.Context) (err error) {
	grpcGatewayRouter := runtime.NewServeMux()
	httpRouter := http.NewServeMux()

	if err = types.RegisterServiceHandlerServer(context.Background(), grpcGatewayRouter, s); err != nil {
		return errors.Wrap(err, "failed to register service handler")
	}

	httpRouter.Handle("/static/service.swagger.json", http.FileServer(http.FS(docs.Docs)))
	httpRouter.HandleFunc("/api", openapiconsole.Handler("TSS service", "/static/service.swagger.json"))
	httpRouter.Handle("/", grpcGatewayRouter)

	srv := &http.Server{Addr: s.gatewayAddr, Handler: httpRouter}
	defer func() {
		if tmpErr := srv.Shutdown(ctx); tmpErr != nil {
			err = errors.Wrap(tmpErr, "failed to shutdown rest server")
		}
	}()

	return srv.ListenAndServe()
}

func (s *Server) SubmitWithdraw(ctx context.Context, request *types.WithdrawRequest) (*types.Empty, error) {
	return &types.Empty{}, s.handler.SubmitWithdraw(ctx, request)
}

func (s *Server) CheckWithdraw(ctx context.Context, request *types.CheckWithdrawRequest) (*types.CheckWithdrawResponse, error) {
	return s.handler.CheckWithdraw(ctx, request)
}
