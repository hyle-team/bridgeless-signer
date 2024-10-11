package btc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"strings"
)

func (p *proxy) GetTransactionStatus(txHash string) (bridgeTypes.TransactionStatus, error) {
	tx, err := p.getTransaction(txHash)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrTxNotFound) {
			return bridgeTypes.TransactionStatusNotFound, nil
		}

		return bridgeTypes.TransactionStatusUnknown, errors.Wrap(err, "failed to get raw transaction")
	}

	// At least one confirmation means that block is mined
	if tx.Confirmations > 0 {
		return bridgeTypes.TransactionStatusSuccessful, nil
	} else {
		return bridgeTypes.TransactionStatusPending, nil
	}
}

func (p *proxy) getTransaction(txHash string) (*btcjson.TxRawResult, error) {
	txHash = strings.TrimPrefix(txHash, "0x")
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse tx hash")
	}

	tx, err := p.chain.Rpc.GetRawTransactionVerbose(hash)
	if err != nil {
		if strings.Contains(err.Error(), "No such mempool or blockchain transaction") {
			return nil, bridgeTypes.ErrTxNotFound
		}
		return nil, errors.Wrap(err, "failed to get raw transaction")
	}

	return tx, nil
}
