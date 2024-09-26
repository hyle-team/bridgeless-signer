package signer

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	privKey *ecdsa.PrivateKey
	addr    common.Address
}

func NewSigner(privKey *ecdsa.PrivateKey) *Signer {
	if privKey == nil {
		return nil
	}

	return &Signer{privKey: privKey, addr: crypto.PubkeyToAddress(privKey.PublicKey)}
}

func (s *Signer) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), s.privKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (s *Signer) SignMessage(message []byte) ([]byte, error) {
	sig, err := crypto.Sign(message[:], s.privKey)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (s *Signer) Address() common.Address {
	return s.addr
}
