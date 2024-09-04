package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"
	"math/big"
)

func (p *proxy) SendBitcoins(data map[string]*big.Int) (string, error) {
	var amounts map[btcutil.Address]btcutil.Amount
	for adrRaw, amount := range data {
		addr, err := btcutil.DecodeAddress(adrRaw, p.chain.Params)
		if err != nil {
			return "", errors.Wrap(err, "failed to decode address")
		}
		if amount == nil {
			return "", errors.New("amount is nil")
		}
		value := amount.Int64()
		if value < MinSatoshisPerOutput {
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
