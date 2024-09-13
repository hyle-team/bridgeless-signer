package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"
	"math/big"
)

func (p *proxy) SendBitcoins(data map[string]*big.Int) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty data")
	}

	amounts := make(map[btcutil.Address]btcutil.Amount, len(data))
	for adrRaw, amount := range data {
		addr, err := btcutil.DecodeAddress(adrRaw, p.chain.Params)
		if err != nil {
			return "", errors.Wrap(err, "failed to decode address")
		}
		if amount == nil {
			return "", errors.New("amount is nil")
		}
		value := amount.Int64()
		if value < minSatoshisPerOutput {
			return "", errors.New("amount is too small")
		}

		amounts[addr] = btcutil.Amount(value)
	}

	hash, err := p.chain.Rpc.SendMany("", amounts)
	if err != nil {
		return "", errors.Wrap(err, "failed to send transaction")
	}

	return hash.String(), nil
}

func (p *proxy) WithdrawalAmountValid(amount *big.Int) bool {
	if amount.IsInt64() && amount.Int64() < minSatoshisPerOutput {
		return false
	}

	return true
}
