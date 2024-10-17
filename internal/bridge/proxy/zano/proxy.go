package zano

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/zano"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
	"regexp"
)

var addressPattern = regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{97}$`)

type proxy struct {
	logger *logan.Entry
	client *zano.Sdk
}

func (p *proxy) Type() bridgeTypes.ChainType {
	return bridgeTypes.ChainTypeZano
}

func (p *proxy) GetTransactionStatus(txHash string) (bridgeTypes.TransactionStatus, error) {
	//TODO implement me
	panic("implement me")
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

func NewBridgeProxy(logger *logan.Entry, sdk *zano.Sdk) bridgeTypes.Proxy {
	return &proxy{logger, sdk}
}