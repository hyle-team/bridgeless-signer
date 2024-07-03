package evm

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/evm/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	bridgeContractAddress = "0xB4C4736D411eE984F0956771a724526E637E2058"
	confirmations         = 10

	depositTxHash    = "0xab0450db07012ccdae9a50e7751250ac3874a02cd95591bb4f7eed2a7ec29c0e"
	depositTxEventId = 2
)

var sepoliaRPCUrl = os.Getenv("SEPOLIA_RPC_URL")
var sepoliaChainConfigurer = func(rpc *ethclient.Client) chain.Chain {
	return chain.Chain{
		Id:            big.NewInt(11155111),
		Rpc:           rpc,
		BridgeAddress: common.HexToAddress(bridgeContractAddress),
		Confirmations: confirmations,
	}
}
var rpcConfigurer = func(t *testing.T, rpc string) *ethclient.Client {
	cli, err := ethclient.Dial(rpc)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to dial rpc"))
	}

	return cli
}

func TestBridgeProxy_GetDepositData(t *testing.T) {
	proxy, err := NewBridgeProxy(sepoliaChainConfigurer(rpcConfigurer(t, sepoliaRPCUrl)))
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to create bridge proxy"))
	}

	testCases := []struct {
		name              string
		depositIdentifier data.DepositIdentifier
		expected          bridgeTypes.DepositData
		err               error
	}{
		{
			name: "must get correct deposit data",
			depositIdentifier: data.DepositIdentifier{
				TxHash:    depositTxHash,
				TxEventId: depositTxEventId,
				ChainId:   "11155111",
			},
			expected: bridgeTypes.DepositData{
				DepositIdentifier: data.DepositIdentifier{
					TxHash:    depositTxHash,
					TxEventId: depositTxEventId,
					ChainId:   "11155111",
				},
				TokenAddress:       common.HexToAddress("0xC0a9A8D46F3859D33aA4067bC89B65Eb793713a4"),
				SourceAddress:      common.HexToAddress("0x274C84fa8c035Cce4443805aF5A2D5B3Fe5Ff5a2"),
				Amount:             big.NewInt(2000000000000000000),
				DestinationChainId: big.NewInt(101),
				DestinationAddress: "ZxCN9egC9BkeL22RobRQD2R9YL7tsfSWxHstXUoGwrPDcfckPNLeaND5EV86fsV5mAhnmPuFSSpiFegeytzbA91T2UrT1rV76",
			},
			err: nil,
		},
		{
			name: "must return deposit not found error",
			depositIdentifier: data.DepositIdentifier{
				TxHash:    "0xab0450db07012ccdae9a50e7751250ac3874a02cd95591bb4f7eed2a7ec29c0e",
				TxEventId: 3,
				ChainId:   "11155111",
			},
			err: bridgeTypes.ErrDepositNotFound,
		},
		{
			name: "must return tx not found error",
			depositIdentifier: data.DepositIdentifier{
				TxHash:    "0xab0450db07012ccdae9a50e7751250ac3874a02cd95591bb4f7eed2a7ec29c01",
				TxEventId: 3,
				ChainId:   "11155111",
			},
			err: bridgeTypes.ErrTxNotFound,
		},
		{
			name: "must return tx failed error",
			depositIdentifier: data.DepositIdentifier{
				TxHash:    "0xf6b05d531c13c2cace02afa65426085ff57813c8ffbcca3a31587b74ebf0da6a",
				TxEventId: 3,
				ChainId:   "11155111",
			},
			err: bridgeTypes.ErrTxFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := proxy.GetDepositData(tc.depositIdentifier)
			if err != nil {
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
			} else {
				assert.Equal(t, tc.expected.TokenAddress, result.TokenAddress)
				assert.Equal(t, tc.expected.SourceAddress, result.SourceAddress)
				assert.Equal(t, tc.expected.Amount.String(), result.Amount.String())
				assert.Equal(t, tc.expected.DestinationChainId.String(), result.DestinationChainId.String())
				assert.Equal(t, tc.expected.DestinationAddress, result.DestinationAddress)
				assert.Equal(t, tc.expected.OriginTxId(), result.OriginTxId())
			}
		})
	}

	t.Run("must return tx not confirmed error", func(t *testing.T) {
		hugeConfirmationsChain := sepoliaChainConfigurer(rpcConfigurer(t, sepoliaRPCUrl))
		hugeConfirmationsChain.Confirmations = 100_000_000_000_000
		proxy, err = NewBridgeProxy(hugeConfirmationsChain)
		if err != nil {
			t.Fatal(errors.Wrap(err, "failed to create bridge proxy"))
		}

		_, err = proxy.GetDepositData(data.DepositIdentifier{
			TxHash:    depositTxHash,
			TxEventId: depositTxEventId,
			ChainId:   "11155111",
		})

		assert.True(t, errors.Is(err, bridgeTypes.ErrTxNotConfirmed))
	})
}
