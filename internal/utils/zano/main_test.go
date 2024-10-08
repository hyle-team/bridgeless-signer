package gosdk

import (
	"dosdk/types"
	"fmt"
	"log"
	"testing"
)

func Test_transfer(t *testing.T) {
	zano := NewZanoSDK("http://localhost:12111/json_rpc")
	dst := make([]types.Destination, 0)
	dst = append(dst, types.Destination{
		Address: "ZxCrkj75u8b218WAiFMiwccwJ6q6GLDLBd4acpHqiyp178qPLQbJR3e1MWHd8hc1h5e4oVWKs9t2VLThQooFidFi1MVAk3kVL",
		Amount:  1,
		AssetID: "7d3f348fbebfffc4e61a3686189cf870ea393e1c88b8f636acbfdacf9e4b2db2",
	})

	res, err := zano.Transfer("bruh", nil, dst)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func Test_getTx(t *testing.T) {
	zano := NewZanoSDK("http://localhost:12111/json_rpc")

	res, err := zano.GetTransaction("9b6fcc0a46d6b5c05b2ab145e63ffa63fa050c41c2c737f42f559877c07ed539")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func Test_deployAsset(t *testing.T) {
	zano := NewZanoSDK("http://localhost:12111/json_rpc")

	description := types.AssetDescriptor{
		DecimalPoint:   12,
		FullName:       "bruh asset",
		HiddenSupply:   false,
		MetaInfo:       "Stable and private",
		Owner:          "ZxDphM9gFU359BXfg2BsPi4xrfapivmTi1c1pvogvD3dbAdha4iCosCWup8YkyitrvdAH15Cin65C2AFpA3AF6cJ2amvcNF7w",
		Ticker:         "ZABC",
		TotalMaxSupply: 1000000000000000000,
		CurrentSupply:  1000000000000,
		//OwnerEthPubKey: "03619c6bc485172e5852a7266e070880d44c6377b64c3c7aa4a3e9435be0cc10ef",
	}

	dst := make([]types.Destination, 0)
	dst = append(dst, types.Destination{
		Address: "ZxCrkj75u8b218WAiFMiwccwJ6q6GLDLBd4acpHqiyp178qPLQbJR3e1MWHd8hc1h5e4oVWKs9t2VLThQooFidFi1MVAk3kVL",
		Amount:  10,
		AssetID: "",
	})
	res, err := zano.DeployAsset(description, dst)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
