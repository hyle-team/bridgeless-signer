package zano

import (
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	zanoTypes "github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
	"math/big"
)

func (p *proxy) WithdrawalAmountValid(amount *big.Int) bool {
	if amount.Cmp(big.NewInt(0)) != 1 {
		return false
	}

	return true
}

func (p *proxy) EmitAssetUnsigned(data data.DepositData) (*bridgeTypes.UnsignedTransaction, error) {
	destination := zanoTypes.Destination{
		Address: data.DestinationAddress,
		Amount:  data.WithdrawalAmount.Uint64(),
		// leaving empty here
		AssetID: "",
	}

	raw, err := p.chain.Client.EmitAsset(data.DestinationTokenAddress, destination)
	if err != nil {
		return nil, errors.Wrap(err, "failed to emit unsigned asset")
	}

	signingData := raw.DataForExternalSigning
	txDetails, err := p.chain.Client.TxDetails(
		signingData.OutputsAddresses,
		signingData.UnsignedTx,
		// leaving empty
		"",
		signingData.TxSecretKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse tx details")
	}

	return &bridgeTypes.UnsignedTransaction{
		ExpectedTxHash: txDetails.VerifiedTxID,
		FinalizedTx:    signingData.FinalizedTx,
		Data:           signingData.UnsignedTx,
	}, nil
}

func (p *proxy) EmitAssetSigned(signedTx bridgeTypes.SignedTransaction) (string, error) {
	_, err := p.chain.Client.SendExtSignedAssetTX(
		signedTx.Signature,
		signedTx.ExpectedTxHash,
		signedTx.FinalizedTx,
		signedTx.Data,
		// TODO: investigate
		true,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to emit signed asset")
	}

	return bridgeTypes.HexPrefix + signedTx.ExpectedTxHash, nil
}
