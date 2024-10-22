package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
)

// minSatoshisPerOutput calculated for P2PKH
const minSatoshisPerOutput = 547

type BridgeProxy interface {
	bridgeTypes.Proxy
	SendBitcoins(map[string]*big.Int) (txHash string, err error)
}

type proxy struct {
	chain  chain.Bitcoin
	logger *logan.Entry
}

func NewBridgeProxy(ch chain.Bitcoin, logger *logan.Entry) BridgeProxy {
	return &proxy{chain: ch, logger: logger}
}

func (*proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeBitcoin
}

func (p *proxy) AddressValid(addr string) bool {
	_, err := btcutil.DecodeAddress(addr, p.chain.Params)
	return err == nil
}

func (p *proxy) TransactionHashValid(hash string) bool {
	return bridgeTypes.DefaultTransactionHashPattern.MatchString(hash)
}
