package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	bridgetypes "github.com/hyle-team/bridgeless-core/x/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"strings"
)

func (c *Connector) GetDestinationTokenInfo(
	srcChainId string,
	srcTokenAddr common.Address,
	dstChainId string,
) (bridgetypes.TokenInfo, error) {
	req := bridgetypes.QueryGetTokenPair{
		SrcChain:   srcChainId,
		SrcAddress: strings.ToLower(srcTokenAddr.String()),
		DstChain:   dstChainId,
	}

	resp, err := c.querier.GetTokenPair(context.Background(), &req)
	if err != nil {
		if errors.Is(err, bridgetypes.ErrTokenPairNotFound.GRPCStatus().Err()) {
			return bridgetypes.TokenInfo{}, types.ErrPairNotFound
		}

		return bridgetypes.TokenInfo{}, errors.Wrap(err, "failed to get token pair")
	}

	return resp.Info, nil

}
