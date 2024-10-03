package signer

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
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

func (s *Signer) setPrefix(message []byte) []byte {
	lenMessage := []byte(fmt.Sprintf("%d", len(message)))
	prefix := []byte("\x19Ethereum Signed Message:\n")
	prefixedMessage := bytes.Join([][]byte{prefix, lenMessage, message}, nil)
	return crypto.Keccak256(prefixedMessage)
}

func (s *Signer) SignMessage(message []byte) ([]byte, error) {
	sig, err := crypto.Sign(s.setPrefix(message)[:], s.privKey)
	if err != nil {
		return nil, err
	}
	sig[64] += 27

	return sig, nil
}

func (s *Signer) Address() common.Address {
	return s.addr
}
