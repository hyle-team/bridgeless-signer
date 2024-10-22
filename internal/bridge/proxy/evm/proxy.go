package evm

import (
	"bytes"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"strings"
)

const (
	EventDepositedNative = "DepositedNative"
	EventDepositedERC20  = "DepositedERC20"
)

var events = []string{
	EventDepositedNative,
	EventDepositedERC20,
}

type BridgeProxy interface {
	bridgeTypes.Proxy
	GetSignHash(data data.DepositData) ([]byte, error)
}

type proxy struct {
	chain         chain.Evm
	contractABI   abi.ABI
	depositEvents []abi.Event
	logger        *logan.Entry
}

// NewBridgeProxy creates a new bridge proxy for the given chain.
func NewBridgeProxy(chain chain.Evm, logger *logan.Entry) (BridgeProxy, error) {
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

	return &proxy{
		chain:         chain,
		contractABI:   bridgeAbi,
		depositEvents: depositEvents,
		logger:        logger,
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

func (p *proxy) TransactionHashValid(hash string) bool {
	return bridgeTypes.DefaultTransactionHashPattern.MatchString(hash)
}
