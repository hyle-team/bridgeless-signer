package processor

import (
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
	switch {
	case err == nil:
		data = &bridgeTypes.FormWithdrawalRequest{
			DepositDbId: req.DepositDbId,
			Data:        *depositData,
		}
	case errors.Is(err, bridgeTypes.ErrTxFailed),
		errors.Is(err, bridgeTypes.ErrDepositNotFound):
		reprocessable = false
	case errors.Is(err, bridgeTypes.ErrTxNotConfirmed),
		errors.Is(err, bridgeTypes.ErrTxPending):
		// explicitly marking as reprocessable
		reprocessable = true
	default:
		// unexpected error occurred - marking as reprocessable
		reprocessable = true
	}

	err = errors.Wrap(err, "failed to get deposit data")

	return
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
	switch {
	case errors.Is(err, tokens.ErrSourceTokenNotSupported),
		errors.Is(err, tokens.ErrDestinationTokenNotSupported):
		return nil, false, err
	case err != nil:
		return nil, true, errors.Wrap(err, "failed to get destination token address")
	default:
		req.Data.DestinationTokenAddress = &dstTokenAddress
	}

	tx, err := proxy.FormWithdrawalTransaction(req.Data)
	switch {
	case err == nil:
		request = &bridgeTypes.WithdrawalRequest{
			Data:        req.Data,
			DepositDbId: req.DepositDbId,
			Transaction: tx,
		}
	case errors.Is(err, bridgeTypes.ErrInvalidReceiverAddress):
		reprocessable = false
	default:
		reprocessable = true
	}

	err = errors.Wrap(err, "failed to form withdrawal transaction")

	return
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
	if deposit.WithdrawalTxHash != nil {
		return false, errors.New("withdrawal transaction was already sent")
	}

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId.String())
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return false, bridgeTypes.ErrChainNotSupported
		}
		return true, errors.Wrap(err, "failed to get proxy")
	}

	err = errors.Wrap(proxy.SendWithdrawalTransaction(req.Transaction), "failed to send withdrawal transaction")
	if err == nil {
		err = errors.Wrap(p.db.SetWithdrawalTx(
			req.DepositDbId, req.Transaction.Hash().Hex(), req.Data.DestinationChainId.String(),
		), "failed to set withdrawal tx")
	}

	// TODO: should be reprocessable or not?
	return true, err
}

func (p *Processor) ProcessSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.WithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	tx, err := p.signer.SignTx(req.Transaction, req.Data.DestinationChainId)
	if err == nil {
		res = &bridgeTypes.WithdrawalRequest{
			Data:        req.Data,
			DepositDbId: req.DepositDbId,
			Transaction: tx,
		}
	}

	// TODO: should be reprocessable or not?
	return res, true, errors.Wrap(err, "failed to sign withdrawal transaction")
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
