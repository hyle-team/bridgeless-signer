package operations

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type WithdrawERC20Content struct {
	tokenAddress []byte
	amount       []byte
	receiver     []byte
	txhash       []byte
	txnonce      []byte
	chainID      []byte
	isWrapped    []byte
}

func (w WithdrawERC20Content) CalculateHash() []byte {
	return crypto.Keccak256(w.tokenAddress, w.amount, w.receiver, w.txhash, w.txnonce, w.chainID, w.isWrapped)
}
