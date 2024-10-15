package gosdk

import (
	"github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
)

const (
	// wallet methods
	searchForTransactionsMethod = "search_for_transactions"
	deployAssetMethod           = "deploy_asset"
	emitAssetMethod             = "emit_asset"
	burnAssetMethod             = "burn_asset"
	transferMethod              = "transfer"

	sendExtSignedAssetTxMethod = "send_ext_signed_asset_tx"
	// node methods
	decryptTxDetailsMethod = "decrypt_tx_details"

	defaultMixin = 15
)

type ZanoSDK struct {
	client *Client
}

func NewZanoSDK(walletRPC, nodeRPC string) *ZanoSDK {
	return &ZanoSDK{
		client: NewClient(walletRPC, nodeRPC),
	}
}

// Transfer Make new payment transaction from the wallet
// service []types.ServiceEntry can be empty.
// wallet rpc api method
func (z ZanoSDK) Transfer(comment string, service []types.ServiceEntry, destinations []types.Destination) (*types.TransferResponse, error) {
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
		Fee:                     "10000000000",
		HideReceiver:            true,
		Mixin:                   defaultMixin,
		PaymentID:               "",
		PushPayer:               false,
		ServiceEntriesPermanent: true,
	}

	resp := new(types.TransferResponse)
	if err := z.client.Call(transferMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetTransaction Search for transactions in the wallet by few parameters
// Pass a hash without 0x prefix
// If past non-existed tx hash node will return a node last tx
// If past empty string instead of a hash node will return all tx for this wallet
// If GetTxResponse contains nill in each field (in, out, pool) it means that transaction
// in pending and there aren`t ways to send some value from this asset
// wallet rpc api method
func (z ZanoSDK) GetTransaction(txid string) (*types.GetTxResponse, error) {
	req := types.GetTxParams{
		FilterByHeight: false,
		In:             true,
		MaxHeight:      0,
		MinHeight:      0,
		Out:            true,
		Pool:           false,
		TxID:           txid,
	}
	resp := new(types.GetTxResponse)
	if err := z.client.Call(searchForTransactionsMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// EmitAsset Emmit new coins of the asset, that is controlled by this wallet.
// assetId must be non-empty and without prefix 0x
// wallet rpc api method
func (z ZanoSDK) EmitAsset(assetId string, destinations []types.Destination) (*types.EmitAssetResponse, error) {
	req := types.EmitAssetParams{
		AssetID:                assetId,
		Destinations:           destinations,
		DoNotSplitDestinations: false,
	}

	resp := new(types.EmitAssetResponse)
	if err := z.client.Call(emitAssetMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// BurnAsset Burn some owned amount of the coins for the given asset.
// https://docs.zano.org/docs/build/rpc-api/wallet-rpc-api/burn_asset/
// assetId must be non-empty and without prefix 0x
// wallet rpc api method
func (z ZanoSDK) BurnAsset(assetId string, amount string) (*types.BurnAssetResponse, error) {
	req := types.BurnAssetParams{
		AssetID:    assetId,
		BurnAmount: amount,
	}

	resp := new(types.BurnAssetResponse)
	if err := z.client.Call(burnAssetMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// DeployAsset Deploy new asset in the system.
// https://docs.zano.org/docs/build/rpc-api/wallet-rpc-api/deploy_asset
// Asset ID inside destinations can be omitted
// wallet rpc api method
func (z ZanoSDK) DeployAsset(assetDescriptor types.AssetDescriptor, destinations []types.Destination) (*types.DeployAssetResponse, error) {
	req := types.DeployAssetParams{
		AssetDescriptor:        assetDescriptor,
		Destinations:           destinations,
		DoNotSplitDestinations: false,
	}

	resp := new(types.DeployAssetResponse)
	if err := z.client.Call(deployAssetMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}

// TxDetails Decrypts transaction private information. Should be used only with your own local daemon for security reasons.
// node rpc api method
func (z ZanoSDK) TxDetails(outputAddress []string, txBlob, txID, txSecretKey string) (*types.DecryptTxDetailsResponse, error) {
	req := types.DecryptTxDetailsParams{
		OutputsAddresses: outputAddress,
		TxBlob:           txBlob,
		TxID:             txID,
		TxSecretKey:      txSecretKey,
	}

	resp := new(types.DecryptTxDetailsResponse)
	if err := z.client.Call(decryptTxDetailsMethod, resp, req, false); err != nil {
		return nil, err
	}

	return resp, nil
}

// SendExtSignedAssetTX Inserts externally made asset ownership signature into the given transaction and broadcasts it.
// wallet rpc api method
func (z ZanoSDK) SendExtSignedAssetTX(ethSig, expectedTXID, finalizedTx, unsignedTx string, unlockTransfersOnFail bool) (*types.SendExtSignedAssetTXResult, error) {
	req := types.SendExtSignedAssetTXParams{
		EthSig:                ethSig,
		ExpectedTxID:          expectedTXID,
		FinalizedTx:           finalizedTx,
		UnlockTransfersOnFail: unlockTransfersOnFail,
		UnsignedTx:            unsignedTx,
	}

	resp := new(types.SendExtSignedAssetTXResult)
	if err := z.client.Call(sendExtSignedAssetTxMethod, resp, req, true); err != nil {
		return nil, err
	}

	return resp, nil
}
