package evm

import (
	"bytes"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	"github.com/hyle-team/bridgeless-signer/internal/evm/chain"
	"github.com/pkg/errors"
)

var (
	ErrChainNotSupported = errors.New("chain not supported")
	ErrTxPending         = errors.New("transaction is pending")
	ErrTxFailed          = errors.New("transaction failed")
	ErrDepositNotFound   = errors.New("deposit not found")
	ErrTxNotConfirmed    = errors.New("transaction not confirmed")
)

type bridgeProxy struct {
	chains       map[string]chain.Chain
	abi          abi.ABI
	depositEvent abi.Event
}

func NewBridgeProxy(chains []chain.Chain) BridgeProxy {
	chainsMap := make(map[string]chain.Chain)

	for _, c := range chains {
		chainsMap[c.Id.String()] = c
	}

	bridgeAbi, err := abi.JSON(strings.NewReader(contracts.BridgeMetaData.ABI))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse bridge ABI"))
	}
	depositEvent, ok := bridgeAbi.Events[DepositEvent]
	if !ok {
		panic(errors.New("wrong bridge ABI events"))
	}

	return &bridgeProxy{chains: chainsMap, abi: bridgeAbi, depositEvent: depositEvent}
}

func (p *bridgeProxy) SupportsChain(chainId string) bool {
	_, ok := p.chains[chainId]
	return ok
}

func (p *bridgeProxy) IsDepositLog(log *types.Log) bool {
	if log == nil {
		return false
	}

	return bytes.Equal(log.Topics[0].Bytes(), p.depositEvent.ID.Bytes())
}
