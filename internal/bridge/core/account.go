package core

import (
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common/hexutil"
	secp256k1 "github.com/hyle-team/bridgeless-core/crypto/ethsecp256k1"
	"github.com/pkg/errors"
)

const hrp = "bridge"

type Account struct {
	privKey *secp256k1.PrivKey
	addr    string
}

func NewAccount(privKey string) (*Account, error) {
	key := &secp256k1.PrivKey{Key: hexutil.MustDecode(privKey)}
	address, err := bech32.ConvertAndEncode(hrp, key.PubKey().Address().Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert and encode address")
	}

	return &Account{
		privKey: key,
		addr:    address,
	}, nil
}

func (a *Account) PrivateKey() *secp256k1.PrivKey {
	return a.privKey
}

func (a *Account) PublicKey() types.PubKey {
	return a.privKey.PubKey()
}

func (a *Account) CosmosAddress() string {
	return a.addr
}
