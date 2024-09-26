package evm

import (
	"bytes"
	"context"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
	"regexp"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

const DepositEvent = "BridgeIn"

var txHashPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)

type proxy struct {
	chain          chain.Evm
	bridgeContract *contracts.Bridge
	contractABI    abi.ABI
	depositEvent   abi.Event
	signerAddr     common.Address
	signerNonce    uint64
	nonceM         sync.Mutex
	logger         *logan.Entry
}

// NewBridgeProxy creates a new bridge proxy for the given chain.
// We need signer address to obtain the nonce for the signer when forming a new transaction.
func NewBridgeProxy(chain chain.Evm, signerAddr common.Address, logger *logan.Entry) (bridgeTypes.Proxy, error) {
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
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer nonce")
	}

	return &proxy{
		chain:          chain,
		contractABI:    bridgeAbi,
		depositEvent:   depositEvent,
		bridgeContract: bridgeContract,
		signerAddr:     signerAddr,
		signerNonce:    nonce,
		logger:         logger,
	}, nil
}

func (p *proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeEVM
}

func (p *proxy) isDepositLog(log *types.Log) bool {
	if log == nil || len(log.Topics) == 0 {
		return false
	}

	return bytes.Equal(log.Topics[0].Bytes(), p.depositEvent.ID.Bytes()) && len(log.Topics) == 2
}

func (p *proxy) AddressValid(addr string) bool {
	return common.IsHexAddress(addr)
}

func (p *proxy) SendBitcoins(map[string]*big.Int) (txHash string, err error) {
	return "", bridgeTypes.ErrNotImplemented
}

func (p *proxy) TransactionHashValid(hash string) bool {
	return txHashPattern.MatchString(hash)
}
