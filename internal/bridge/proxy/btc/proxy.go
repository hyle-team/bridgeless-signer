package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

// minSatoshisPerOutput calculated for P2PKH
const minSatoshisPerOutput = 547

type proxy struct {
	chain  chain.Bitcoin
	logger *logan.Entry
}

func NewBridgeProxy(ch chain.Bitcoin, logger *logan.Entry) bridgeTypes.Proxy {
	return &proxy{chain: ch, logger: logger}
}

func (*proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeBitcoin
}

func (p *proxy) GetSignHash(_ data.DepositData) ([]byte, error) {
	return nil, bridgeTypes.ErrNotImplemented
}

func (p *proxy) AddressValid(addr string) bool {
	_, err := btcutil.DecodeAddress(addr, p.chain.Params)
	return err == nil
}

func (p *proxy) TransactionHashValid(hash string) bool {
	return bridgeTypes.DefaultTransactionHashPattern.MatchString(hash)
}

func (p *proxy) EmitAssetUnsigned(data data.DepositData) (*bridgeTypes.UnsignedTransaction, error) {
	return nil, bridgeTypes.ErrNotImplemented
}

func (p *proxy) EmitAssetSigned(transaction bridgeTypes.SignedTransaction) (txHash string, err error) {
	return "", bridgeTypes.ErrNotImplemented
}
