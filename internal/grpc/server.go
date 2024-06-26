package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyle-team/bridgeless-signer/docs"
	"github.com/hyle-team/bridgeless-signer/internal/config"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/ignite/cli/ignite/pkg/openapiconsole"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var _ types.ServiceServer = &Server{}

// Server is used to implement types.ServiceServer GRPC server.
type Server struct {
	listener    net.Listener
	gatewayAddr string
}

// NewServer creates a new GRPC server.
func NewServer(
	listener net.Listener,
	gatewayCfg config.HTTPGatewayConfig,
) *Server {
	return &Server{
		listener:    listener,
		gatewayAddr: gatewayCfg.Address,
	}
}

// RunGRPC starts the GRPC server.
func (s *Server) RunGRPC(_ context.Context) error {
	grpcServer := grpc.NewServer()
	types.RegisterServiceServer(grpcServer, s)
	return grpcServer.Serve(s.listener)
}

// RunHTTPGateway starts the HTTP gateway server.
func (s *Server) RunHTTPGateway(ctx context.Context) error {
	grpcGatewayRouter := runtime.NewServeMux()
	httpRouter := http.NewServeMux()

	if err := types.RegisterServiceHandlerServer(context.Background(), grpcGatewayRouter, s); err != nil {
		return errors.Wrap(err, "failed to register service handler")
	}

	httpRouter.Handle("/static/service.swagger.json", http.FileServer(http.FS(docs.Docs)))
	httpRouter.HandleFunc("/api", openapiconsole.Handler("TSS service", "/static/service.swagger.json"))
	httpRouter.Handle("/", grpcGatewayRouter)

	srv := &http.Server{Addr: s.gatewayAddr, Handler: httpRouter}
	defer func() {
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("failed to shutdown HTTP server", err)
		}
	}()

	return srv.ListenAndServe()
}

func (s *Server) Echo(ctx context.Context, message *types.StringMessage) (*types.StringMessage, error) {
	return message, nil
}
