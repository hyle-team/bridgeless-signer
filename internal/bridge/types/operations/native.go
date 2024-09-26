package operations

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type WithdrawNativeContent struct {
	amount   []byte
	receiver []byte
	txhash   []byte
	txnonce  []byte
	chainID  []byte
}

func (w WithdrawNativeContent) CalculateHash() []byte {
	return crypto.Keccak256(w.amount, w.receiver, w.txhash, w.txnonce, w.chainID)
}
