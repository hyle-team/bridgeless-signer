package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hyle-team/bridgeless-signer/docs"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/ctx"
	grpcServer "github.com/hyle-team/bridgeless-signer/internal/core/api/grpc"
	"github.com/hyle-team/bridgeless-signer/internal/core/api/middlewares"
	api "github.com/hyle-team/bridgeless-signer/internal/core/api/types"
	rabbitTypes "github.com/hyle-team/bridgeless-signer/internal/core/rabbitmq/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
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

var _ grpcServer.ServiceServer = grpcImplementation{}

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
	srv := s.grpcServer()

	// graceful shutdown
	go func() { <-ctx.Done(); srv.GracefulStop(); s.logger.Info("grpc serving stopped: context canceled") }()

	s.logger.Info("grpc serving started")
	return srv.Serve(s.grpc)
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
		s.logger.Info("http serving stopped: context canceled")
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
	_ = grpcServer.RegisterServiceHandlerServer(ctxt, grpcGatewayRouter, grpcImplementation{})

	// grpc interceptor not working here
	router.With(ape.CtxMiddleware(s.ctxExtenders...)).Mount("/", grpcGatewayRouter)
	router.With(
		ape.CtxMiddleware(s.ctxExtenders...),
		// extending with websocket middleware
		middlewares.HijackedConnectionCloser(ctxt),
	).Get("/ws/check/{origin_tx_id}", CheckWithdrawalWs)

	router.Mount("/static/api.swagger.json", http.FileServer(http.FS(docs.Docs)))
	router.HandleFunc("/api", openapiconsole.Handler("Signer service", "/static/api.swagger.json"))

	return router
}

func (s *server) grpcServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			ContextExtenderInterceptor(s.ctxExtenders...),
			LoggerInterceptor(s.logger),
			// RecoveryInterceptor should be the last one
			RecoveryInterceptor(s.logger),
		),
	)

	grpcServer.RegisterServiceServer(srv, grpcImplementation{})
	reflection.Register(srv)

	return srv
}
