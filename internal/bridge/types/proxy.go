package types

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
)

var (
	ErrChainNotSupported      = errors.New("chain not supported")
	ErrTxPending              = errors.New("transaction is pending")
	ErrTxFailed               = errors.New("transaction failed")
	ErrTxNotFound             = errors.New("transaction not found")
	ErrDepositNotFound        = errors.New("deposit not found")
	ErrTxNotConfirmed         = errors.New("transaction not confirmed")
	ErrInvalidReceiverAddress = errors.New("invalid receiver address")
	ErrInvalidDepositedAmount = errors.New("invalid deposited amount")
	ErrNotImplemented         = errors.New("not implemented")
	ErrInvalidScriptPubKey    = errors.New("invalid script pub keyå")
)

type ChainType string

const (
	ChainTypeEVM     ChainType = "evm"
	ChainTypeBitcoin ChainType = "bitcoin"
	ChainTypeOther   ChainType = "other"
)

func (c ChainType) Validate() error {
	switch c {
	case ChainTypeEVM, ChainTypeBitcoin, ChainTypeOther:
		return nil
	default:
		return errors.New("invalid chain type")
	}
}

type TransactionStatus int8

const (
	TransactionStatusPending TransactionStatus = iota
	TransactionStatusSuccessful
	TransactionStatusFailed
	TransactionStatusNotFound
	TransactionStatusUnknown
)

type Proxy interface {
	Type() ChainType
	GetTransactionStatus(txHash string) (TransactionStatus, error)
	GetDepositData(id data.DepositIdentifier) (*data.DepositData, error)
	AddressValid(addr string) bool

	// Ethereum-specific methods
	FormWithdrawalTransaction(data data.DepositData) (*types.Transaction, error)
	SendWithdrawalTransaction(signedTx *types.Transaction) error

	// Bitcoin-specific methods
	SendBitcoins(map[string]*big.Int) (txHash string, err error)
}

type ProxiesRepository interface {
	Proxy(chainId string) (Proxy, error)
	SupportsChain(chainId string) bool
}
