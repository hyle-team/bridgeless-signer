package zano

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/chain"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
	"regexp"
)

var addressPattern = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{97}$`)

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
	return bridgeTypes.DefaultTransactionHashPattern.MatchString(hash)
}

func (p *proxy) SendBitcoins(m map[string]*big.Int) (txHash string, err error) {
	return "", bridgeTypes.ErrNotImplemented
}

func (p *proxy) GetSignHash(data data.DepositData) ([]byte, error) {
	return nil, bridgeTypes.ErrNotImplemented
}

func NewBridgeProxy(chain chain.Zano, logger *logan.Entry) bridgeTypes.Proxy {
	return &proxy{logger, chain}
}
