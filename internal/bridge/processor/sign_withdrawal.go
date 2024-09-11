package processor

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
	"math/big"
)

func (p *Processor) ProcessSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	chainId, set := new(big.Int).SetString(req.Data.DestinationChainId, 10)
	if !set {
		return nil, false, errors.New("invalid destination chain id")
	}

	tx, err := p.signer.SignTx(req.Transaction, chainId)
	if err != nil {
		// TODO: should be reprocessable or not?
		return res, true, errors.Wrap(err, "failed to sign withdrawal transaction")
	}

	return &bridgeTypes.WithdrawalRequest{
		Data:        req.Data,
		DepositDbId: req.DepositDbId,
		Transaction: tx,
	}, false, nil
}
