package evm

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *proxy) GetDepositData(id data.DepositIdentifier) (*data.DepositData, error) {
	txReceipt, err := p.GetTransactionReceipt(common.HexToHash(id.TxHash))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction receipt")
	}
	if txReceipt.Status != types.ReceiptStatusSuccessful {
		return nil, bridgeTypes.ErrTxFailed
	}

	if err = p.validateConfirmations(txReceipt); err != nil {
		return nil, errors.Wrap(err, "failed to validate confirmations")
	}

	if len(txReceipt.Logs) < id.TxEventId+1 {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	log := txReceipt.Logs[id.TxEventId]
	if !p.isDepositLog(log) {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	var event contracts.BridgeBridgeIn
	if err = p.contractABI.UnpackIntoInterface(&event, DepositEvent, log.Data); err != nil {
		return nil, errors.Wrap(err, "failed to unpack deposit event")
	}
	// parsing indexed event parameter that is always present and not in the parsed even data
	event.Token = common.HexToAddress(log.Topics[1].Hex())

	return &data.DepositData{
		DepositIdentifier:  id,
		DestinationChainId: event.ChainId,
		DestinationAddress: event.DstAddress,
		SourceAddress:      event.SrcAddress,
		Amount:             event.Amount,
		TokenAddress:       event.Token,
		Block:              int64(log.BlockNumber),
	}, nil
}

func (p *proxy) validateConfirmations(receipt *types.Receipt) error {
	curHeight, err := p.chain.EvmRpc.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get current block number")
	}

	// including the current block
	if receipt.BlockNumber.Uint64()+uint64(p.chain.Confirmations)-1 > curHeight {
		return bridgeTypes.ErrTxNotConfirmed
	}

	return nil
}
