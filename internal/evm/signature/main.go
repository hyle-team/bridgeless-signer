package signature

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Signer struct {
	privKey *ecdsa.PrivateKey
}

func NewSigner(privKey *ecdsa.PrivateKey) *Signer {
	return &Signer{privKey: privKey}
}

func (s *Signer) SignTx(tx *types.Transaction, chainID int64) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(chainID)), s.privKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
