package zano

import (
	"github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
)

const (
	defaultMixin = 15
	defaultFee   = "10000000000"
)

type Sdk struct {
	client *Client
}

func NewSDK(walletRPC, nodeRPC string) *Sdk {
	return &Sdk{
		client: NewClient(walletRPC, nodeRPC),
	}
}

// Transfer Make new payment transaction from the wallet
// service []types.ServiceEntry can be empty.
// wallet rpc api method
func (z Sdk) Transfer(comment string, service []types.ServiceEntry, destinations []types.Destination) (*types.TransferResponse, error) {
	if service == nil || len(service) == 0 {
		service = []types.ServiceEntry{}
	}
	if destinations == nil || len(destinations) == 0 {
		return nil, errors.New("destinations must be non-empty")
	}
	req := types.TransferParams{
		Comment:                 comment,
		Destinations:            destinations,
		ServiceEntries:          service,
		Fee:                     defaultFee,
		HideReceiver:            true,
		Mixin:                   defaultMixin,
		PaymentID:               "",
		PushPayer:               false,
		ServiceEntriesPermanent: true,
	}

	resp := new(types.TransferResponse)
	if err := z.client.Call(types.TransferMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetTransactions Search for transactions in the wallet by few parameters
// Pass a hash without 0x prefix
// If past empty string instead of a hash node will return all tx for this wallet
// wallet rpc api method
func (z Sdk) GetTransactions(txid string) (*types.GetTxResponse, error) {
	req := types.GetTxParams{
		FilterByHeight: false,
		In:             true,
		MaxHeight:      0,
		MinHeight:      0,
		Out:            true,
		Pool:           true,
		TxID:           txid,
	}
	resp := new(types.GetTxResponse)
	if err := z.client.Call(types.SearchForTransactionsMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// EmitAsset Emmit new coins of the asset, that is controlled by this wallet.
// assetId must be non-empty and without prefix 0x
// wallet rpc api method
func (z Sdk) EmitAsset(assetId string, destinations ...types.Destination) (*types.EmitAssetResponse, error) {
	if len(destinations) == 0 {
		return nil, errors.New("destinations must be non-empty")
	}

	req := types.EmitAssetParams{
		AssetID:                assetId,
		Destinations:           destinations,
		DoNotSplitDestinations: false,
	}

	resp := new(types.EmitAssetResponse)
	if err := z.client.Call(types.EmitAssetMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// BurnAsset Burn some owned amount of the coins for the given asset.
// https://docs.zano.org/docs/build/rpc-api/wallet-rpc-api/burn_asset/
// assetId must be non-empty and without prefix 0x
// wallet rpc api method
func (z Sdk) BurnAsset(assetId string, amount string) (*types.BurnAssetResponse, error) {
	req := types.BurnAssetParams{
		AssetID:    assetId,
		BurnAmount: amount,
	}

	resp := new(types.BurnAssetResponse)
	if err := z.client.Call(types.BurnAssetMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// DeployAsset Deploy new asset in the system.
// https://docs.zano.org/docs/build/rpc-api/wallet-rpc-api/deploy_asset
// Asset ID inside destinations can be omitted
// wallet rpc api method
func (z Sdk) DeployAsset(assetDescriptor types.AssetDescriptor, destinations []types.Destination) (*types.DeployAssetResponse, error) {
	req := types.DeployAssetParams{
		AssetDescriptor:        assetDescriptor,
		Destinations:           destinations,
		DoNotSplitDestinations: false,
	}

	resp := new(types.DeployAssetResponse)
	if err := z.client.Call(types.DeployAssetMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// TxDetails Decrypts transaction private information. Should be used only with your own local daemon for security reasons.
// node rpc api method
func (z Sdk) TxDetails(outputAddress []string, txBlob, txID, txSecretKey string) (*types.DecryptTxDetailsResponse, error) {
	req := types.DecryptTxDetailsParams{
		OutputsAddresses: outputAddress,
		TxBlob:           txBlob,
		TxID:             txID,
		TxSecretKey:      txSecretKey,
	}

	resp := new(types.DecryptTxDetailsResponse)
	if err := z.client.Call(types.DecryptTxDetailsMethod, resp, req, false); err != nil {
		return nil, err
	}

	return resp, nil
}

// SendExtSignedAssetTX Inserts externally made asset ownership signature into the given transaction and broadcasts it.
// wallet rpc api method
func (z Sdk) SendExtSignedAssetTX(ethSig, expectedTXID, finalizedTx, unsignedTx string, unlockTransfersOnFail bool) (*types.SendExtSignedAssetTXResult, error) {
	req := types.SendExtSignedAssetTXParams{
		EthSig:                ethSig,
		ExpectedTxID:          expectedTXID,
		FinalizedTx:           finalizedTx,
		UnlockTransfersOnFail: unlockTransfersOnFail,
		UnsignedTx:            unsignedTx,
	}

	resp := new(types.SendExtSignedAssetTXResult)
	if err := z.client.Call(types.SendExtSignedAssetTxMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

func (z Sdk) CurrentHeight() (uint64, error) {
	resp := new(types.GetHeightResponse)
	err := z.client.CallRaw(types.GetHeightMethod, resp)

	return resp.Height, err
}
