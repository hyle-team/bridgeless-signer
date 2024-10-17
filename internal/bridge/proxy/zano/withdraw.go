package zano

import (
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"math/big"
)

func (p *proxy) WithdrawalAmountValid(amount *big.Int) bool {
	if amount.Cmp(big.NewInt(0)) != 1 {
		return false
	}

	return true
}

func (p *proxy) GetSignHash(data data.DepositData) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
