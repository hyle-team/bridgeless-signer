package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/contracts"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
)

func (p *bridgeProxy) GetDepositData(id data.DepositIdentifier) (*bridgeTypes.DepositData, error) {
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
	if !p.IsDepositLog(log) {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	var event contracts.BridgeBridgeIn
	if err = p.contractABI.UnpackIntoInterface(&event, DepositEvent, log.Data); err != nil {
		return nil, errors.Wrap(err, "failed to unpack deposit event")
	}
	// parsing indexed event parameter that is always present and not in the parsed even data
	event.Token = common.HexToAddress(log.Topics[1].Hex())

	return &bridgeTypes.DepositData{
		DepositIdentifier: id,
		// TODO: change uint8 to uint256
		DestinationChainId: big.NewInt(int64(event.ChainId)),
		DestinationAddress: event.DstAddress,
		SourceAddress:      event.SrcAddress,
		Amount:             event.Amount,
		TokenAddress:       event.Token,
	}, nil
}

func (p *bridgeProxy) validateConfirmations(receipt *types.Receipt) error {
	curHeight, err := p.chain.Rpc.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get current block number")
	}

	// including the current block
	if receipt.BlockNumber.Uint64()+uint64(p.chain.Confirmations)-1 > curHeight {
		return bridgeTypes.ErrTxNotConfirmed
	}

	return nil
}
