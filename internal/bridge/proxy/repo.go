package proxy

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/btc"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/evm"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

type proxiesRepository struct {
	proxies map[string]bridgeTypes.Proxy
}

func NewProxiesRepository(chains []chain.Chain, signer common.Address) (proxyRepo bridgeTypes.ProxiesRepository, err error) {
	proxiesMap := make(map[string]bridgeTypes.Proxy)

	for _, ch := range chains {
		var proxy bridgeTypes.Proxy

		switch ch.Type {
		case bridgeTypes.ChainTypeEVM:
			var evmChain chain.Evm
			evmChain, err = ch.Evm()
			if err != nil {
				return nil, errors.Wrap(err, "failed to init evm chain")
			}
			proxy, err = evm.NewBridgeProxy(evmChain, signer)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to create proxy for chain %s", ch.Id))
			}
		case bridgeTypes.ChainTypeBitcoin:
			var bitcoinChain chain.Bitcoin
			bitcoinChain, err = ch.Bitcoin()
			if err != nil {
				return nil, errors.Wrap(err, "failed to init bitcoin")
			}
			proxy = btc.NewBridgeProxy(bitcoinChain)
		default:
			return nil, errors.Errorf("unknown chain type %s", ch.Type)
		}

		proxiesMap[ch.Id] = proxy
	}

	return &proxiesRepository{proxies: proxiesMap}, nil
}

func (p proxiesRepository) Proxy(chainId string) (bridgeTypes.Proxy, error) {
	proxy, ok := p.proxies[chainId]
	if !ok {
		return nil, bridgeTypes.ErrChainNotSupported
	}

	return proxy, nil
}

func (p proxiesRepository) SupportsChain(chainId string) bool {
	_, ok := p.proxies[chainId]
	return ok
}