package types

type ServiceEntrie struct {
	Body        string `json:"body"`
	Flags       int    `json:"flags"`
	Instruction string `json:"instruction"`
	Security    string `json:"security"`
	ServiceID   string `json:"service_id"`
}

type Destination struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
	AssetID string `json:"asset_id"`
}

type TransferParams struct {
	Comment                 string          `json:"comment"`
	Destinations            []Destination   `json:"destinations"`
	Fee                     int64           `json:"fee"`
	HideReceiver            bool            `json:"hide_receiver"`
	Mixin                   int             `json:"mixin"`
	PaymentID               string          `json:"payment_id"`
	PushPayer               bool            `json:"push_payer"`
	ServiceEntries          []ServiceEntrie `json:"service_entries"`
	ServiceEntriesPermanent bool            `json:"service_entries_permanent"`
}

type TransferResponse struct {
	TxHash        string `json:"tx_hash"`
	TxSize        int    `json:"tx_size"`
	TxUnsignedHex string `json:"tx_unsigned_hex"`
}

type GetTxParams struct {
	FilterByHeight bool   `json:"filter_by_height"`
	In             bool   `json:"in"`
	MaxHeight      int    `json:"max_height"`
	MinHeight      int    `json:"min_height"`
	Out            bool   `json:"out"`
	Pool           bool   `json:"pool"`
	TxID           string `json:"tx_id"`
}

type GetTxResponse struct {
	In   []Transaction `json:"in"`
	Out  []Transaction `json:"out"`
	Pool []Transaction `json:"pool"`
}

type Transaction struct {
	Amount                int64          `json:"amount"`
	Comment               string         `json:"comment"`
	Contract              []Contract     `json:"contract"`
	EmployedEntries       interface{}    `json:"employed_entries"`
	Fee                   int64          `json:"fee"`
	Height                int            `json:"height"`
	IsIncome              bool           `json:"is_income"`
	IsMining              bool           `json:"is_mining"`
	IsMixing              bool           `json:"is_mixing"`
	IsService             bool           `json:"is_service"`
	PaymentID             string         `json:"payment_id"`
	RemoteAddresses       []string       `json:"remote_addresses"`
	RemoteAliases         []string       `json:"remote_aliases"`
	ServiceEntries        []ServiceEntry `json:"service_entries"`
	ShowSender            bool           `json:"show_sender"`
	Subtransfers          []SubTransfer  `json:"subtransfers"`
	Timestamp             int            `json:"timestamp"`
	TransferInternalIndex int            `json:"transfer_internal_index"`
	TxBlobSize            int            `json:"tx_blob_size"`
	TxHash                string         `json:"tx_hash"`
	TxType                int            `json:"tx_type"`
	UnlockTime            int            `json:"unlock_time"`
}

type Contract struct {
	CancelExpirationTime int             `json:"cancel_expiration_time"`
	ContractID           string          `json:"contract_id"`
	ExpirationTime       int             `json:"expiration_time"`
	Height               int             `json:"height"`
	IsA                  bool            `json:"is_a"`
	PaymentID            string          `json:"payment_id"`
	PrivateDetailes      PrivateDetailes `json:"private_detailes"`
	State                int             `json:"state"`
	Timestamp            int             `json:"timestamp"`
}

type PrivateDetailes struct {
	AAddr   string `json:"a_addr"`
	APledge int    `json:"a_pledge"`
	BAddr   string `json:"b_addr"`
	BPledge int    `json:"b_pledge"`
	C       string `json:"c"`
	T       string `json:"t"`
	ToPay   int    `json:"to_pay"`
}

type ServiceEntry struct {
	Body        string `json:"body"`
	Flags       int    `json:"flags"`
	Instruction string `json:"instruction"`
	Security    string `json:"security"`
	ServiceID   string `json:"service_id"`
}

type SubTransfer struct {
	Amount   int64  `json:"amount"`
	AssetID  string `json:"asset_id"`
	IsIncome bool   `json:"is_income"`
}

type EmitAssetParams struct {
	AssetID                string        `json:"asset_id"`
	Destination            []Destination `json:"destination"`
	DoNotSplitDestinations bool          `json:"do_not_split_destinations"`
}

type EmitAssetResponse struct {
	TxID string `json:"tx_id"`
}

type AssetDescriptor struct {
	DecimalPoint   int    `json:"decimal_point"`
	FullName       string `json:"full_name"`
	HiddenSupply   bool   `json:"hidden_supply"`
	MetaInfo       string `json:"meta_info"`
	Owner          string `json:"owner"`
	Ticker         string `json:"ticker"`
	TotalMaxSupply int    `json:"total_max_supply"`
	OwnerEthPubKey string `json:"owner_eth_pub_key"`
	CurrentSupply  int    `json:"current_supply"`
}

type DeployAssetParams struct {
	AssetDescriptor        `json:"asset_descriptor"`
	Destinations           []Destination `json:"destinations"`
	DoNotSplitDestinations bool          `json:"do_not_split_destinations"`
}

type DeployAssetResponse struct {
	NewAssetId string `json:"new_asset_id"`
	TxID       string `json:"tx_id"`
}
