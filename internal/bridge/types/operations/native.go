package operations

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
)

type WithdrawNativeContent struct {
	Amount   []byte
	Receiver []byte
	TxHash   []byte
	TxNonce  []byte
	ChainID  []byte
}

func (w WithdrawNativeContent) CalculateHash() []byte {
	return crypto.Keccak256(
		w.Amount,
		w.Receiver,
		w.TxHash,
		w.TxNonce,
		w.ChainID,
	)
}

func (w WithdrawNativeContent) Equals(other []byte) bool {
	return bytes.Equal(other, w.CalculateHash())
}
