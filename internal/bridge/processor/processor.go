package processor

import (
	"github.com/hyle-team/bridgeless-signer/internal/bridge/signer"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/hyle-team/bridgeless-signer/pkg/types"
	"github.com/pkg/errors"
)

type Processor struct {
	proxies bridgeTypes.ProxiesRepository
	db      data.DepositsQ
	signer  *signer.Signer
}

func New(proxies bridgeTypes.ProxiesRepository, db data.DepositsQ, signer *signer.Signer) *Processor {
	return &Processor{proxies: proxies, db: db, signer: signer}
}

func (p *Processor) ProcessGetDepositRequest(req bridgeTypes.GetDepositRequest) (data *bridgeTypes.FormWithdrawRequest, reprocessable bool, err error) {
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
		data = &bridgeTypes.FormWithdrawRequest{
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

func (p *Processor) ProcessFormWithdrawRequest(req bridgeTypes.FormWithdrawRequest) (request *bridgeTypes.WithdrawRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId.String())
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}

	tx, err := proxy.FormWithdrawalTransaction(req.Data)
	switch {
	case err == nil:
		request = &bridgeTypes.WithdrawRequest{
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

func (p *Processor) ProcessSendWithdrawRequest(req bridgeTypes.WithdrawRequest) (reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

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

func (p *Processor) ProcessSignWithdrawRequest(req bridgeTypes.WithdrawRequest) (res *bridgeTypes.WithdrawRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(req.DepositDbId, err, reprocessable) }()

	tx, err := p.signer.SignTx(req.Transaction, req.Data.DestinationChainId)
	if err == nil {
		res = &bridgeTypes.WithdrawRequest{
			Data:        req.Data,
			DepositDbId: req.DepositDbId,
			Transaction: tx,
		}
	}

	// TODO: should be reprocessable or not?
	return res, true, errors.Wrap(err, "failed to sign withdrawal transaction")
}

func (p *Processor) SetWithdrawStatusFailed(id int64) error {
	return errors.Wrap(p.db.UpdateStatus(id, types.WithdrawStatus_FAILED), "failed to update deposit status")
}

func (p *Processor) updateInvalidDepositStatus(id int64, err error, reprocessable bool) error {
	if err == nil || reprocessable {
		return err
	}

	if tempErr := p.db.UpdateStatus(id, types.WithdrawStatus_INVALID); tempErr != nil {
		return errors.Wrap(tempErr, "failed to update deposit status")
	}

	return err
}
