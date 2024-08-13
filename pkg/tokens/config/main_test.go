package config

import (
	"errors"
	"github.com/hyle-team/bridgeless-signer/pkg/tokens"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/kit/kv"
)

func Test_TokenPairer(t *testing.T) {
	getter := kv.NewViperFile("config_test.yaml")
	pairer := NewConfigTokenPairerConfiger(getter).TokenPairer()

	testCases := []struct {
		name       string
		srcChainId *big.Int
		srcToken   common.Address
		dstChainId *big.Int
		expected   common.Address
		err        error
	}{
		{
			name:       "success for chain 80002",
			srcChainId: big.NewInt(80002),
			srcToken:   common.HexToAddress("0x9c9b83Ed9dd4cF8A385b6e318Fb97Cdfc320b627"),
			dstChainId: big.NewInt(1),
			expected:   common.HexToAddress("0x9c9b83Ed9dd4cF8A385b6e318Fb97Cdfc320b627"),
			err:        nil,
		},
		{
			name:       "success for chain 80003",
			srcChainId: big.NewInt(80003),
			srcToken:   common.HexToAddress("0x9c9b83Ed9dd4cF8A385b6e318Fb97Cdfc320b627"),
			dstChainId: big.NewInt(3),
			expected:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
			err:        nil,
		},
		{
			name:       "src token not supported",
			srcChainId: big.NewInt(80003),
			srcToken:   common.HexToAddress("0x555b83Ed9dd4cF8A385b6e318Fb97Cdfc320b555"),
			dstChainId: big.NewInt(3),
			expected:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
			err:        tokens.ErrSourceTokenNotSupported,
		},
		{
			name:       "dst token not supported",
			srcChainId: big.NewInt(80003),
			srcToken:   common.HexToAddress("0x9c9b83Ed9dd4cF8A385b6e318Fb97Cdfc320b627"),
			dstChainId: big.NewInt(7),
			expected:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
			err:        tokens.ErrDestinationTokenNotSupported,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := pairer.GetDestinationTokenAddress(tc.srcChainId, tc.srcToken, tc.dstChainId)
			if err == nil {
				if addr != tc.expected {
					t.Fatalf("expected address %v, got %v", tc.expected, addr)
				}
			} else {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
			}
		})
	}
}
