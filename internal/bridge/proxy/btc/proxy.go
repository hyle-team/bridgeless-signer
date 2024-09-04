package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
)

// MinSatoshisPerOutput calculated for P2PKH
const MinSatoshisPerOutput = 547

type proxy struct {
	chain chain.Bitcoin
}

func NewBridgeProxy(ch chain.Bitcoin) bridgeTypes.Proxy {
	return &proxy{chain: ch}
}

func (*proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeBitcoin
}

func (p *proxy) FormWithdrawalTransaction(data data.DepositData) (*types.Transaction, error) {
	return nil, bridgeTypes.ErrNotImplemented
}

func (p *proxy) SendWithdrawalTransaction(signedTx *types.Transaction) error {
	return bridgeTypes.ErrNotImplemented
}

func (p *proxy) AddressValid(addr string) bool {
	_, err := btcutil.DecodeAddress(addr, p.chain.Params)
	return err == nil
}
