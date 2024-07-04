package handler

import (
	"errors"

	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const internalError = "internal error"

var (
	ErrInternal           = status.Error(codes.Internal, internalError)
	ErrTxAlreadySubmitted = status.Error(codes.AlreadyExists, "transaction already submitted")
	ErrChainNotSupported  = errors.New("chain not supported")
	ErrInvalidOriginTxId  = errors.New("invalid origin tx id")
)

// ServiceHandler is an implementation of the API interface.
type ServiceHandler struct {
	db      data.DepositsQ
	logger  *logan.Entry
	proxies bridgeTypes.ProxiesRepository
}

func NewServiceHandler(
	db data.DepositsQ,
	logger *logan.Entry,
	proxies bridgeTypes.ProxiesRepository,
) *ServiceHandler {
	return &ServiceHandler{
		db:      db,
		logger:  logger,
		proxies: proxies,
	}
}
