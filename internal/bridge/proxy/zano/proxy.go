package zano

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"regexp"
)

var addressPattern = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{97}$`)

type BridgeProxy interface {
	bridgeTypes.Proxy
	EmitAssetUnsigned(data data.DepositData) (*UnsignedTransaction, error)
	EmitAssetSigned(transaction SignedTransaction) (txHash string, err error)
}

type proxy struct {
	logger *logan.Entry
	chain  chain.Zano
}

func (p *proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeZano
}

func (p *proxy) GetDepositData(id data.DepositIdentifier) (*data.DepositData, error) {
	//TODO implement me
	panic("implement me")
}

func (p *proxy) AddressValid(addr string) bool {
	return addressPattern.MatchString(addr)
}

func (p *proxy) TransactionHashValid(hash string) bool {
	return bridge.DefaultTransactionHashPattern.MatchString(hash)
}

func NewBridgeProxy(chain chain.Zano, logger *logan.Entry) BridgeProxy {
	return &proxy{logger, chain}
}
