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

var (
	DepositNative  = "DepositedNative"
	DepositedERC20 = "DepositedERC20"
)

var events = []string{
	DepositNative,
	DepositedERC20,
}

func (p *proxy) GetDepositData(id data.DepositIdentifier) (*data.DepositData, error) {
	txReceipt, err := p.GetTransactionReceipt(common.HexToHash(id.TxHash))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction receipt")
	}
	if txReceipt.Status != types.ReceiptStatusSuccessful {
		return nil, bridgeTypes.ErrTxFailed
	}

	if len(txReceipt.Logs) < id.TxEventId+1 {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	log := txReceipt.Logs[id.TxEventId]
	if !p.isDepositLog(log) {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	if err = p.validateConfirmations(txReceipt); err != nil {
		return nil, errors.Wrap(err, "failed to validate confirmations")
	}

	var unpackedData *data.DepositData

	for _, eventName := range events {
		switch eventName {
		case DepositNative:
			event := new(contracts.BridgeDepositedNative)
			unpackLog, err := p.unpackLog(event, eventName, log, id)
			if err != nil && unpackLog == nil {
				p.logger.Debug(errors.Wrap(err, "failed to unpack event"))
				continue
			}

			unpackedData = unpackLog
			break

		case DepositedERC20:
			event := new(contracts.BridgeDepositedERC20)
			unpackLog, err := p.unpackLog(event, eventName, log, id)
			if err != nil && unpackLog == nil {
				p.logger.Debug(errors.Wrap(err, "failed to unpack event"))
				continue
			}

			unpackedData = unpackLog
			break
		}
	}

	if unpackedData == nil {
		return nil, bridgeTypes.ErrFailedUnpackLogs
	}

	return unpackedData, nil
}

func (p *proxy) unpackLog(eventBody interface{}, eventName string, log *types.Log, id data.DepositIdentifier) (*data.DepositData, error) {
	if err := p.contractABI.UnpackIntoInterface(&eventBody, eventName, log.Data); err != nil {
		return nil, errors.Wrap(err, "failed to unpack deposit event")
	}

	switch eventName {
	case DepositedERC20:
		return &data.DepositData{
			DepositIdentifier:  id,
			DestinationChainId: eventBody.(contracts.BridgeDepositedERC20).Network,
			DestinationAddress: eventBody.(contracts.BridgeDepositedERC20).Receiver,
			DepositAmount:      eventBody.(contracts.BridgeDepositedERC20).Amount,
			TokenAddress:       eventBody.(contracts.BridgeDepositedERC20).Token,
			Block:              int64(log.BlockNumber),
		}, nil

	case DepositNative:
		return &data.DepositData{
			DepositIdentifier:  id,
			DestinationChainId: eventBody.(contracts.BridgeDepositedNative).Network,
			DestinationAddress: eventBody.(contracts.BridgeDepositedNative).Receiver,
			DepositAmount:      eventBody.(contracts.BridgeDepositedNative).Amount,
			Block:              int64(log.BlockNumber),
		}, nil
	default:
		return nil, errors.Errorf("unknown event %s", eventName)
	}
}

func (p *proxy) validateConfirmations(receipt *types.Receipt) error {
	curHeight, err := p.chain.Rpc.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get current block number")
	}

	// including the current block
	if receipt.BlockNumber.Uint64()+p.chain.Confirmations-1 > curHeight {
		return bridgeTypes.ErrTxNotConfirmed
	}

	return nil
}
