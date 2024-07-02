package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *bridgeProxy) GetDepositData(id data.DepositIdentifier) (*DepositData, error) {
	txReceipt, err := p.GetTransactionReceipt(common.HexToHash(id.TxHash), id.ChainId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction receipt")
	}
	if txReceipt.Status != types.ReceiptStatusSuccessful {
		return nil, ErrTxFailed
	}

	if err = p.validateConfirmations(txReceipt, id.ChainId); err != nil {
		return nil, errors.Wrap(err, "failed to validate confirmations")
	}

	if len(txReceipt.Logs) < id.TxEventId+1 {
		return nil, ErrDepositNotFound
	}

	log := txReceipt.Logs[id.TxEventId]
	if !p.IsDepositLog(log) {
		return nil, ErrDepositNotFound
	}

	var event contracts.BridgeBridgeIn
	if err = p.abi.UnpackIntoInterface(&event, DepositEvent, log.Data); err != nil {
		return nil, errors.Wrap(err, "failed to unpack deposit event")
	}
	event.Token = common.HexToAddress(log.Topics[1].Hex())

	return &DepositData{
		DepositIdentifier: id,
		// TODO: change uint8 to uint256
		DestinationChainId: big.NewInt(int64(event.ChainId)),
		DestinationAddress: event.DstAddress,
		SourceAddress:      event.SrcAddress,
		Amount:             event.Amount,
		TokenAddress:       event.Token,
	}, nil
}

func (p *bridgeProxy) validateConfirmations(receipt *types.Receipt, chainId string) error {
	chain, ok := p.chains[chainId]
	if !ok {
		return ErrChainNotSupported
	}

	curHeight, err := chain.Rpc.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get current block number")
	}
	// including the current block
	if receipt.BlockNumber.Uint64()+uint64(chain.Confirmations)-1 > curHeight {
		return ErrTxNotConfirmed
	}

	return nil
}
