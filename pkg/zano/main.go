package gosdk

import (
	"github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
)

const (
	searchForTransactionsMethod = "search_for_transactions"
	deployAssetMethod           = "deploy_asset"
	emitAssetMethod             = "emit_asset"
	transferMethod              = "transfer"
)

type ZanoSDK struct {
	client *Client
}

func NewZanoSDK(url string) *ZanoSDK {
	return &ZanoSDK{
		client: NewClient(url),
	}
}

// service []types.ServiceEntrie can be empty.
func (z ZanoSDK) Transfer(comment string, service []types.ServiceEntrie, destinations []types.Destination) (*types.TransferResponse, error) {
	if service == nil || len(service) == 0 {
		service = []types.ServiceEntrie{}
	}
	if destinations == nil || len(destinations) == 0 {
		return nil, errors.New("destinations must be non-empty")
	}
	req := types.TransferParams{
		Comment:                 comment,
		Destinations:            destinations,
		ServiceEntries:          service,
		Fee:                     10000000000,
		HideReceiver:            true,
		Mixin:                   15,
		PaymentID:               "",
		PushPayer:               false,
		ServiceEntriesPermanent: true,
	}

	resp := new(types.TransferResponse)
	if err := z.client.Call(transferMethod, resp, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// Pass a hash without 0x prefix
// If past non-existed tx hash node will return a node last tx
// If past empty string instead of a hash node will return all tx for this wallet
// If GetTxResponse contains nill in each field (in, out, pool) it means that transaction
// in pending and there arent ways to send some value from this asset
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
	if err := z.client.Call(searchForTransactionsMethod, resp, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// assetId must be non-empty and without prefix 0x
func (z ZanoSDK) EmitAsset(assetId string, destinations []types.Destination) (resp *types.EmitAssetResponse, err error) {
	req := types.EmitAssetParams{
		AssetID:                assetId,
		Destination:            destinations,
		DoNotSplitDestinations: false,
	}

	if err = z.client.Call(emitAssetMethod, resp, req); err != nil {
		return nil, err
	}

	return resp, nil
}

// https://docs.zano.org/docs/build/rpc-api/wallet-rpc-api/deploy_asset
// Asset ID inside destinations can be ommited
func (z ZanoSDK) DeployAsset(AssetDescriptor types.AssetDescriptor, destinations []types.Destination) (resp *types.DeployAssetResponse, err error) {
	req := types.DeployAssetParams{
		AssetDescriptor:        AssetDescriptor,
		Destinations:           destinations,
		DoNotSplitDestinations: false,
	}

	if err = z.client.Call(deployAssetMethod, resp, req); err != nil {
		return nil, err
	}

	return resp, nil
}
