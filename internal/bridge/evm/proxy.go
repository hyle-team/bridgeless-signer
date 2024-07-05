package evm

import (
	"bytes"
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

const DepositEvent = "BridgeIn"

type bridgeProxy struct {
	chain          chain.Chain
	bridgeContract *contracts.Bridge
	contractABI    abi.ABI
	depositEvent   abi.Event
	signerAddr     common.Address
	signerNonce    uint64
}

func NewBridgeProxy(chain chain.Chain, signerAddr common.Address) (bridgeTypes.Proxy, error) {
	bridgeAbi, err := abi.JSON(strings.NewReader(contracts.BridgeMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse bridge ABI")
	}

	depositEvent, ok := bridgeAbi.Events[DepositEvent]
	if !ok {
		return nil, errors.New("wrong bridge ABI events")
	}

	bridgeContract, err := contracts.NewBridge(chain.BridgeAddress, chain.Rpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bridge contract")
	}

	nonce, err := chain.Rpc.PendingNonceAt(context.Background(), signerAddr)

	return &bridgeProxy{
		chain:          chain,
		contractABI:    bridgeAbi,
		depositEvent:   depositEvent,
		bridgeContract: bridgeContract,
		signerAddr:     signerAddr,
		signerNonce:    nonce,
	}, nil
}

func (p *bridgeProxy) IsDepositLog(log *types.Log) bool {
	if log == nil {
		return false
	}

	return bytes.Equal(log.Topics[0].Bytes(), p.depositEvent.ID.Bytes()) && len(log.Topics) == 2
}
