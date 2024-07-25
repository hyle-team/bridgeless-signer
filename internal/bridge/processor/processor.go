package processor

import (
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/tokens"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
)

type Processor struct {
	proxies     bridgeTypes.ProxiesRepository
	db          data.DepositsQ
	signer      *signer.Signer
	tokenPairer tokens.TokenPairer
}

func New(
	proxies bridgeTypes.ProxiesRepository,
	db data.DepositsQ,
	signer *signer.Signer,
	tokenPairer tokens.TokenPairer) *Processor {
	return &Processor{proxies: proxies, db: db, signer: signer, tokenPairer: tokenPairer}
}

func (p *Processor) ProcessGetDepositRequest(req bridgeTypes.GetDepositRequest) (data *bridgeTypes.FormWithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	proxy, err := p.proxies.Proxy(req.DepositIdentifier.ChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return data, false, bridgeTypes.ErrChainNotSupported
		}
		return data, true, errors.Wrap(err, "failed to get proxy")
	}

	depositData, err := proxy.GetDepositData(req.DepositIdentifier)
	if err == nil {
		return &bridgeTypes.FormWithdrawalRequest{
			DepositDbId: req.DepositDbId,
			Data:        *depositData,
		}, false, nil
	}

	reprocessable = true
	if errors.Is(err, bridgeTypes.ErrTxFailed) ||
		errors.Is(err, bridgeTypes.ErrDepositNotFound) {
		reprocessable = false
	}

	return nil, reprocessable, errors.Wrap(err, "failed to get deposit data")
}

func (p *Processor) ProcessFormWithdrawalRequest(req bridgeTypes.FormWithdrawalRequest) (request *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId.String())
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}

	dstTokenAddress, err := p.tokenPairer.GetDestinationTokenAddress(
		req.Data.DepositIdentifier.GetChainId(),
		req.Data.TokenAddress,
		req.Data.DestinationChainId,
	)
	if err != nil {
		reprocessable = true
		if errors.Is(err, tokens.ErrSourceTokenNotSupported) ||
			errors.Is(err, tokens.ErrDestinationTokenNotSupported) {
			reprocessable = false
		}

		return nil, reprocessable, errors.Wrap(err, "failed to get destination token address")
	}
	req.Data.DestinationTokenAddress = &dstTokenAddress

	var tx *ethTypes.Transaction
	txConn := p.db.New()
	err = txConn.Transaction(func() error {
		tmpErr := txConn.SetDepositData(req.Data)
		if tmpErr != nil {
			return errors.Wrap(tmpErr, "failed to save deposit data")
		}

		tx, tmpErr = proxy.FormWithdrawalTransaction(req.Data)
		return errors.Wrap(tmpErr, "failed to form withdrawal transaction")
	})
	if err == nil {
		return &bridgeTypes.WithdrawalRequest{
			Data:        req.Data,
			DepositDbId: req.DepositDbId,
			Transaction: tx,
		}, false, nil
	}

	reprocessable = true
	if errors.Is(err, bridgeTypes.ErrInvalidReceiverAddress) {
		reprocessable = false
	}

	return nil, reprocessable, errors.Wrap(err, "failed to form withdrawal transaction")
}

func (p *Processor) ProcessSendWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	// ensure that withdrawal request was not already processed
	deposit, err := p.db.Get(req.Data.DepositIdentifier)
	if err != nil {
		return true, errors.Wrap(err, "failed to check if deposit already processed")
	}
	if deposit == nil {
		return true, errors.New("deposit was not found in the database")
	}
	if !deposit.WithdrawalAllowed() {
		return false, errors.New("withdrawal transaction was already sent")
	}

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId.String())
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return false, bridgeTypes.ErrChainNotSupported
		}
		return true, errors.Wrap(err, "failed to get proxy")
	}

	// rollback if transaction failed to be sent
	txConn := p.db.New()
	err = txConn.Transaction(func() error {
		if tempErr := txConn.SetWithdrawalTx(
			req.DepositDbId, req.Transaction.Hash().Hex(), req.Data.DestinationChainId.String(),
		); tempErr != nil {
			return errors.Wrap(tempErr, "failed to set withdrawal tx")
		}

		return errors.Wrap(proxy.SendWithdrawalTransaction(req.Transaction), "failed to send withdrawal transaction")
	})
	return err != nil, err
}

func (p *Processor) ProcessSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	tx, err := p.signer.SignTx(req.Transaction, req.Data.DestinationChainId)
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

func (p *Processor) SetWithdrawStatusFailed(id int64) error {
	return errors.Wrap(p.db.UpdateStatus(id, types.WithdrawalStatus_FAILED), "failed to update deposit status")
}

func (p *Processor) updateInvalidDepositStatus(id int64, err error, reprocessable bool) error {
	if err == nil || reprocessable {
		return err
	}

	if tempErr := p.db.UpdateStatus(id, types.WithdrawalStatus_INVALID); tempErr != nil {
		return errors.Wrap(tempErr, "failed to update deposit status")
	}

	return err
}
