package types

import (
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
	"regexp"
)

var DefaultTransactionHashPattern = regexp.MustCompile("^0x[a-fA-F0-9]{64}$")

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
	ErrInvalidScriptPubKey    = errors.New("invalid script pub key")
	ErrFailedUnpackLogs       = errors.New("failed to unpack logs")
	ErrUnsupportedEvent       = errors.New("unsupported event")
	ErrUnsupportedContract    = errors.New("unsupported contract")
)

type ChainType string

const (
	ChainTypeEVM     ChainType = "evm"
	ChainTypeBitcoin ChainType = "bitcoin"
	ChainTypeZano    ChainType = "zano"
	ChainTypeOther   ChainType = "other"
)

func (c ChainType) Validate() error {
	switch c {
	case ChainTypeEVM, ChainTypeBitcoin, ChainTypeZano, ChainTypeOther:
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
	TransactionHashValid(hash string) bool
	WithdrawalAmountValid(amount *big.Int) bool

	// Ethereum-specific methods
	GetSignHash(data data.DepositData) ([]byte, error)

	// Bitcoin-specific methods
	SendBitcoins(map[string]*big.Int) (txHash string, err error)

	// Zano-specific methods
	// EmitAsset(data data.DepositData)
}

type ProxiesRepository interface {
	Proxy(chainId string) (Proxy, error)
	SupportsChain(chainId string) bool
}
