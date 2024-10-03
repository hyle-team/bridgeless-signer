package operations

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
)

type WithdrawERC20Content struct {
	DestinationTokenAddress []byte
	Amount                  []byte
	Receiver                []byte
	TxHash                  []byte
	TxNonce                 []byte
	ChainID                 []byte
	IsWrapped               []byte
}

func (w WithdrawERC20Content) CalculateHash() []byte {
	return crypto.Keccak256(
		w.DestinationTokenAddress,
		w.Amount,
		w.Receiver,
		w.TxHash,
		w.TxNonce,
		w.ChainID,
		w.IsWrapped,
	)
}

func (w WithdrawERC20Content) Equals(other []byte) bool {
	return bytes.Equal(other, w.CalculateHash())
}
