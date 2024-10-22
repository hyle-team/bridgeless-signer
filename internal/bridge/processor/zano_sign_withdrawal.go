package processor

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hyle-team/bridgeless-signer/internal/bridge/proxy/zano"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/pkg/errors"
)

func (p *Processor) ProcessZanoSignWithdrawalRequest(req bridgeTypes.WithdrawalRequest) (res *bridgeTypes.ZanoSignedWithdrawalRequest, reprocessable bool, err error) {
	defer func() { err = p.updateInvalidDepositStatus(err, reprocessable, req.DepositDbId) }()

	proxy, err := p.proxies.Proxy(req.Data.DestinationChainId)
	if err != nil {
		if errors.Is(err, bridgeTypes.ErrChainNotSupported) {
			return nil, false, bridgeTypes.ErrChainNotSupported
		}
		return nil, true, errors.Wrap(err, "failed to get proxy")
	}
	if proxy.Type() != bridgeTypes.ChainTypeZano {
		return nil, false, bridgeTypes.ErrChainNotSupported
	}
	zanoProxy, ok := proxy.(zano.BridgeProxy)
	if !ok {
		return nil, false, bridgeTypes.ErrChainNotSupported
	}

	unsignedTx, err := zanoProxy.EmitAssetUnsigned(req.Data)
	if err != nil {
		return nil, true, errors.Wrap(err, "failed to emit unsigned tx")
	}

	signData := hexutil.MustDecode(bridgeTypes.HexPrefix + unsignedTx.ExpectedTxHash)
	signature, err := p.signer.SignMessage(signData)
	if err != nil {
		return nil, true, errors.Wrap(err, "failed to sign message")
	}

	encodedSignature := hexutil.Encode(signature)
	// stripping redundant hex-prefix and recovery byte (two hex-characters)
	strippedSignature := encodedSignature[2 : len(encodedSignature)-2]
	signedTx := bridgeTypes.SignedTransaction{
		Signature:           strippedSignature,
		UnsignedTransaction: *unsignedTx,
	}

	return &bridgeTypes.ZanoSignedWithdrawalRequest{
		DepositDbId: req.DepositDbId,
		Data:        req.Data,
		Transaction: signedTx,
	}, false, nil
}
