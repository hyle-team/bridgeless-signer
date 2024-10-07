package operations

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
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

func NewWithdrawERC20Content(event data.DepositData) (*WithdrawERC20Content, error) {
	destinationChainID, ok := new(big.Int).SetString(event.DestinationChainId, 10)
	if !ok {
		return nil, errors.New("invalid chain id")
	}

	return &WithdrawERC20Content{
		Amount:                  To32Bytes(event.DepositAmount.Bytes()),
		Receiver:                hexutil.MustDecode(event.DestinationAddress),
		TxHash:                  hexutil.MustDecode(event.TxHash),
		TxNonce:                 IntTo32Bytes(event.TxEventId),
		ChainID:                 To32Bytes(destinationChainID.Bytes()),
		DestinationTokenAddress: event.DestinationTokenAddress.Bytes(),
		IsWrapped:               BoolToBytes(event.IsWrappedToken),
	}, nil
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
