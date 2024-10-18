package evm

import (
	"bytes"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

const (
	EventDepositedNative = "DepositedNative"
	EventDepositedERC20  = "DepositedERC20"
)

var events = []string{
	EventDepositedNative,
	EventDepositedERC20,
}

type proxy struct {
	chain          chain.Evm
	bridgeContract *contracts.Bridge
	contractABI    abi.ABI
	depositEvents  []abi.Event
	nonceM         sync.Mutex
	logger         *logan.Entry
}

// NewBridgeProxy creates a new bridge proxy for the given chain.
func NewBridgeProxy(chain chain.Evm, logger *logan.Entry) (bridgeTypes.Proxy, error) {
	bridgeAbi, err := abi.JSON(strings.NewReader(contracts.BridgeMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse bridge ABI")
	}

	depositEvents := make([]abi.Event, len(events))
	for i, event := range events {
		depositEvent, ok := bridgeAbi.Events[event]
		if !ok {
			return nil, errors.New("wrong bridge ABI events")
		}
		depositEvents[i] = depositEvent
	}

	bridgeContract, err := contracts.NewBridge(chain.BridgeAddress, chain.Rpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bridge contract")
	}

	return &proxy{
		chain:          chain,
		contractABI:    bridgeAbi,
		depositEvents:  depositEvents,
		bridgeContract: bridgeContract,
		logger:         logger,
	}, nil
}

func (p *proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeEVM
}

func (p *proxy) getDepositLogType(log *types.Log) string {
	if log == nil || len(log.Topics) == 0 {
		return ""
	}

	for _, event := range p.depositEvents {
		isEqual := bytes.Equal(log.Topics[0].Bytes(), event.ID.Bytes())
		if isEqual {
			return event.Name
		}
	}

	return ""
}

func (p *proxy) AddressValid(addr string) bool {
	return common.IsHexAddress(addr)
}

func (p *proxy) SendBitcoins(_ map[string]*big.Int) (txHash string, err error) {
	return "", bridgeTypes.ErrNotImplemented
}

func (p *proxy) TransactionHashValid(hash string) bool {
	return bridgeTypes.DefaultTransactionHashPattern.MatchString(hash)
}

func (p *proxy) EmitAssetUnsigned(data data.DepositData) (*bridgeTypes.UnsignedTransaction, error) {
	return nil, bridgeTypes.ErrNotImplemented
}

func (p *proxy) EmitAssetSigned(transaction bridgeTypes.SignedTransaction) (txHash string, err error) {
	return "", bridgeTypes.ErrNotImplemented
}
