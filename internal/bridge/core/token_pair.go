package core

import (
	"context"
	bridgetypes "github.com/hyle-team/bridgeless-core/v12/x/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"strings"
)

func (c *Connector) GetDestinationTokenInfo(
	srcChainId string,
	srcTokenAddr string,
	dstChainId string,
) (bridgetypes.TokenInfo, error) {
	req := bridgetypes.QueryGetTokenPair{
		SrcChain:   srcChainId,
		SrcAddress: strings.ToLower(srcTokenAddr),
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
