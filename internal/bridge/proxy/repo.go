package proxy

import (
	"fmt"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/btc"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/evm"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type proxiesRepository struct {
	proxies map[string]bridgeTypes.Proxy
}

func NewProxiesRepository(chains []chain.Chain, logger *logan.Entry) (proxyRepo bridgeTypes.ProxiesRepository, err error) {
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
			proxy, err = evm.NewBridgeProxy(evmChain, logger)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to create proxy for chain %s", ch.Id))
			}
		case bridgeTypes.ChainTypeBitcoin:
			var bitcoinChain chain.Bitcoin
			bitcoinChain, err = ch.Bitcoin()
			if err != nil {
				return nil, errors.Wrap(err, "failed to init bitcoin")
			}
			proxy = btc.NewBridgeProxy(bitcoinChain, logger)
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
