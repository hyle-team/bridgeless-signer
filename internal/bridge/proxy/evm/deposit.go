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
	txReceipt, from, err := p.GetTransactionReceipt(common.HexToHash(id.TxHash))
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
			eventBody := new(contracts.BridgeDepositedNative)
			if err = p.contractABI.UnpackIntoInterface(eventBody, eventName, log.Data); err != nil {
				p.logger.Debug(errors.Wrap(err, "failed to unpack event"))
				continue
			}

			unpackedData = &data.DepositData{
				DepositIdentifier:  id,
				DestinationChainId: eventBody.Network,
				DestinationAddress: eventBody.Receiver,
				DepositAmount:      eventBody.Amount,
				Block:              int64(log.BlockNumber),
				SourceAddress:      from.String(),
			}

			break

		case DepositedERC20:
			eventBody := new(contracts.BridgeDepositedERC20)
			if err = p.contractABI.UnpackIntoInterface(eventBody, eventName, log.Data); err != nil {
				p.logger.Debug(errors.Wrap(err, "failed to unpack event"))
				continue
			}

			unpackedData = &data.DepositData{
				DepositIdentifier:  id,
				DestinationChainId: eventBody.Network,
				DestinationAddress: eventBody.Receiver,
				DepositAmount:      eventBody.Amount,
				TokenAddress:       eventBody.Token,
				IsWrappedToken:     eventBody.IsWrapped,
				Block:              int64(log.BlockNumber),
				SourceAddress:      from.String(),
			}

			break
		}
	}

	if unpackedData == nil {
		return nil, bridgeTypes.ErrFailedUnpackLogs
	}

	return unpackedData, nil
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
