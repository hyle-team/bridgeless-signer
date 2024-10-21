package btc

import (
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"math/big"
	"strings"
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

	if err := p.ensureWalletLoaded(); err != nil {
		return "", errors.Wrap(err, "failed to ensure wallet loaded")
	}

	hash, err := p.chain.Rpc.SendMany("", amounts)
	if err != nil {
		return "", errors.Wrap(err, "failed to send transaction")
	}

	return bridgeTypes.HexPrefix + hash.String(), nil
}

func (p *proxy) WithdrawalAmountValid(amount *big.Int) bool {
	if amount.Cmp(big.NewInt(minSatoshisPerOutput)) == -1 {
		return false
	}

	return true
}

func (p *proxy) ensureWalletLoaded() error {
	info, err := p.chain.Rpc.GetWalletInfo()
	if err != nil {
		if !strings.HasPrefix(err.Error(), fmt.Sprintf("%v", btcjson.ErrRPCWalletNotFound)) {
			return errors.Wrap(err, "failed to get wallet info")
		}
	} else {
		if info.WalletName == p.chain.Wallet {
			return nil
		}
	}

	_, err = p.chain.Rpc.LoadWallet(p.chain.Wallet)

	return errors.Wrap(err, "failed to load wallet")
}
