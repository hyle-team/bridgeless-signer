package zano

import (
	"encoding/base64"
	"encoding/json"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	zanoTypes "github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
	"math/big"
)

type destinationData struct {
	Address string `json:"destination_address"`
	ChainId string `json:"destination_chain_id"`
}

func (p *proxy) GetDepositData(id data.DepositIdentifier) (*data.DepositData, error) {
	transaction, _, err := p.GetTransaction(id.TxHash, true, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction")
	}
	if transaction == nil {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	if err = p.validateConfirmations(transaction.Height); err != nil {
		return nil, errors.Wrap(err, "failed to validate confirmations")
	}

	if !transaction.Ado.IsValidAssetBurn() {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	if len(transaction.ServiceEntries) < id.TxEventId+1 {
		return nil, bridgeTypes.ErrDepositNotFound
	}
	addr, chainId, err := parseDestinationData(transaction.ServiceEntries[id.TxEventId])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse destination data")
	}

	var depositor string
	if len(transaction.RemoteAddresses) > 0 {
		depositor = transaction.RemoteAddresses[0]
	}

	return &data.DepositData{
		DepositIdentifier:  id,
		DestinationChainId: chainId,
		DestinationAddress: addr,
		SourceAddress:      depositor,
		// FIXME: find real deposit (burn) amount
		DepositAmount: new(big.Int).SetUint64(transaction.Amount),
		TokenAddress:  *transaction.Ado.OptAssetId,
		Block:         int64(transaction.Height),
	}, nil
}

func (p *proxy) validateConfirmations(txHeight uint64) error {
	if txHeight == 0 {
		return bridgeTypes.ErrTxPending
	}

	currentHeight, err := p.chain.Client.CurrentHeight()
	if err != nil {
		return errors.Wrap(err, "failed to get current height")
	}

	if currentHeight-txHeight < p.chain.Confirmations {
		return bridgeTypes.ErrTxNotConfirmed
	}

	return nil
}

func parseDestinationData(entry zanoTypes.ServiceEntry) (addr, chainId string, err error) {
	raw, err := base64.StdEncoding.DecodeString(entry.Body)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to decode base64 body")
	}

	var dstData destinationData
	if err = json.Unmarshal(raw, &dstData); err != nil {
		return "", "", errors.Wrap(err, "failed to unmarshal json data")
	}

	return dstData.Address, dstData.ChainId, nil
}
