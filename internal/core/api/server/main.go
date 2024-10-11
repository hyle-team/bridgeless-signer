package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyle-team/bridgeless-signer/docs"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/ctx"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/middlewares"
	api "github.com/hyle-team/bridgeless-signer/internal/core/api/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/ignite/cli/ignite/pkg/openapiconsole"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"time"
)

var _ types.ServiceServer = grpcImplementation{}

type grpcImplementation struct{}

type server struct {
	grpc net.Listener
	http net.Listener

	logger       *logan.Entry
	ctxExtenders []func(context.Context) context.Context
}

// NewServer creates a new GRPC server.
func NewServer(
	grpc net.Listener,
	http net.Listener,
	db data.DepositsQ,
	proxies bridgeTypes.ProxiesRepository,
	producer rabbitTypes.Producer,
	logger *logan.Entry,
) api.Server {
	return &server{
		grpc:   grpc,
		http:   http,
		logger: logger,

		ctxExtenders: []func(context.Context) context.Context{
			ctx.LoggerProvider(logger),
			ctx.DBProvider(db),
			ctx.ProxiesProvider(proxies),
			ctx.ProducerProvider(producer),
		},
	}
}

func (s *server) RunGRPC(ctx context.Context) error {
	grpcServer := s.grpcServer()

	// graceful shutdown
	go func() { <-ctx.Done(); grpcServer.GracefulStop(); s.logger.Info("grpc serving stopped") }()

	s.logger.Info("grpc serving started")
	return grpcServer.Serve(s.grpc)
}

func (s *server) RunHTTP(ctxt context.Context) error {
	srv := &http.Server{Handler: s.httpRouter(ctxt)}

	// graceful shutdown
	go func() {
		<-ctxt.Done()
		shutdownDeadline, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownDeadline); err != nil {
			s.logger.WithError(err).Error("failed to shutdown http server")
		}
		s.logger.Info("http serving stopped")
	}()

	s.logger.Info("http serving started")
	if err := srv.Serve(s.http); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *server) httpRouter(ctxt context.Context) http.Handler {
	router := chi.NewRouter()
	router.Use(ape.LoganMiddleware(s.logger), ape.RecoverMiddleware(s.logger))

	// pointing to grpc implementation
	grpcGatewayRouter := runtime.NewServeMux()
	_ = types.RegisterServiceHandlerServer(ctxt, grpcGatewayRouter, grpcImplementation{})

	// grpc interceptor not working here
	router.With(ape.CtxMiddleware(s.ctxExtenders...)).Mount("/", grpcGatewayRouter)
	router.With(
		ape.CtxMiddleware(s.ctxExtenders...),
		// extending with websocket middleware
		middlewares.HijackedConnectionCloser(ctxt),
	).Get("/ws/check/{origin_tx_id}", CheckWithdrawalWs)

	router.Mount("/static/service.swagger.json", http.FileServer(http.FS(docs.Docs)))
	router.HandleFunc("/api", openapiconsole.Handler("Signer service", "/static/service.swagger.json"))

	return router
}

func (s *server) grpcServer() *grpc.Server {
	interceptor := grpc.UnaryInterceptor(api.GRPCContextExtenderInterceptor(s.ctxExtenders...))
	grpcServer := grpc.NewServer(interceptor)

	types.RegisterServiceServer(grpcServer, grpcImplementation{})
	reflection.Register(grpcServer)

	return grpcServer
}
