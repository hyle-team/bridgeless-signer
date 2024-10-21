package zano

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	zanoTypes "github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
	"strings"
)

func (p *proxy) GetTransactionStatus(txHash string) (bridgeTypes.TransactionStatus, error) {
	tx, inPool, err := p.GetTransaction(txHash, true, true, true)
	if err != nil {
		return bridgeTypes.TransactionStatusUnknown, err
	}
	if tx == nil {
		return bridgeTypes.TransactionStatusNotFound, nil
	}
	if inPool {
		return bridgeTypes.TransactionStatusPending, nil
	}

	return bridgeTypes.TransactionStatusSuccessful, nil
}

func (p *proxy) GetTransaction(txHash string, searchIn, searchOut, searchPool bool) (res *zanoTypes.Transaction, pool bool, err error) {
	txHash = strings.TrimPrefix(txHash, bridgeTypes.HexPrefix)
	resp, err := p.chain.Client.GetTransactions(txHash)
	if err != nil {
		return res, false, errors.Wrap(err, "failed to get transaction")
	}

	if searchIn && len(resp.In) != 0 {
		return &resp.In[0], false, nil
	}
	if searchOut && len(resp.Out) != 0 {
		return &resp.Out[0], false, nil
	}
	// TODO: ASK ZANO SIDE TO FIX THAT
	if searchPool && len(resp.Pool) != 0 {
		for _, tx := range resp.Pool {
			if tx.TxHash == txHash {
				return &tx, true, nil
			}
		}
	}

	return res, false, nil
}
