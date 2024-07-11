package config

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/tokens"
)

type configTokenPairer struct {
	tokenPairInfo map[string]map[string]common.Address
}

func NewTokenPairer(tokenPairs []TokenPairs) tokens.TokenPairer {
	tokenPairInfo := make(map[string]map[string]common.Address)
	for _, token := range tokenPairs {
		pairs := make(map[string]common.Address)
		for _, pair := range token.Pairs {
			pairs[pair.ChainId.String()] = pair.Address
		}

		tokenPairInfo[formTokenKey(token.Token.ChainId, token.Token.Address)] = pairs
	}

	return &configTokenPairer{tokenPairInfo: tokenPairInfo}
}

func (p *configTokenPairer) GetDestinationTokenAddress(
	srcChainId *big.Int,
	srcTokenAddr common.Address,
	dstChainId *big.Int,
) (common.Address, error) {
	key := formTokenKey(srcChainId, srcTokenAddr)

	pairs, ok := p.tokenPairInfo[key]
	if !ok {
		return common.Address{}, tokens.ErrSourceTokenNotSupported
	}

	dstTokenAddr, ok := pairs[dstChainId.String()]
	if !ok {
		return common.Address{}, tokens.ErrDestinationTokenNotSupported
	}

	return dstTokenAddr, nil
}

func formTokenKey(srcChainId *big.Int, srcTokenAddr common.Address) string {
	return fmt.Sprintf("%s-%s", srcChainId.String(), strings.ToLower(srcTokenAddr.String()))
}
