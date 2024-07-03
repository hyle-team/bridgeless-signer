package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

var (
	ErrChainNotSupported      = errors.New("chain not supported")
	ErrTxPending              = errors.New("transaction is pending")
	ErrTxFailed               = errors.New("transaction failed")
	ErrTxNotFound             = errors.New("transaction not found")
	ErrDepositNotFound        = errors.New("deposit not found")
	ErrTxNotConfirmed         = errors.New("transaction not confirmed")
	ErrInvalidReceiverAddress = errors.New("invalid receiver address")
)

type Proxy interface {
	GetDepositData(id data.DepositIdentifier) (*DepositData, error)
	IsDepositLog(log *types.Log) bool
	GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error)
	FormWithdrawalTransaction(data DepositData) (*types.Transaction, error)
	SendWithdrawalTransaction(signedTx *types.Transaction) error
}

type ProxiesRepository interface {
	Proxy(chainId string) (Proxy, error)
	SupportsChain(chainId string) bool
}

type DepositData struct {
	data.DepositIdentifier

	DestinationChainId *big.Int
	SourceAddress      common.Address
	DestinationAddress string
	Amount             *big.Int
	TokenAddress       common.Address
}

func (d DepositData) OriginTxId() string {
	return d.DepositIdentifier.String()
}
