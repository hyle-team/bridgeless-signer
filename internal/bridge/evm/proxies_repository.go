package evm

import (
	"fmt"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"

	"github.com/ethereum/go-ethereum/common"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

type proxiesRepository struct {
	proxies map[string]bridgeTypes.Proxy
}

func NewProxiesRepository(chains []chain.Chain, signer common.Address) (bridgeTypes.ProxiesRepository, error) {
	proxiesMap := make(map[string]bridgeTypes.Proxy)

	for _, c := range chains {
		proxy, err := NewBridgeProxy(c, signer)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create proxy for chain %s", c.Id.String()))
		}

		proxiesMap[c.Id.String()] = proxy
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
