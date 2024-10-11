package gosdk

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"log"
	"testing"
)

func Test_transfer(t *testing.T) {
	zano := NewZanoSDK("http://localhost:10500/json_rpc", "")
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
	zano := NewZanoSDK("http://localhost:10500/json_rpc", "")

	res, err := zano.GetTransaction("9b6fcc0a46d6b5c05b2ab145e63ffa63fa050c41c2c737f42f559877c07ed539")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func Test_deployAsset(t *testing.T) {
	zano := NewZanoSDK("http://localhost:10500/json_rpc", "")

	description := types.AssetDescriptor{
		DecimalPoint:   12,
		FullName:       "test asset",
		HiddenSupply:   false,
		MetaInfo:       "TEST",
		Owner:          "",
		Ticker:         "TEST",
		TotalMaxSupply: "1000000000000000000",
		CurrentSupply:  "1000000000000",
		OwnerEthPubKey: "03619c6bc485172e5852a7266e070880d44c6377b64c3c7aa4a3e9435be0cc10ef",
	}

	dst := make([]types.Destination, 0)
	dst = append(dst, types.Destination{
		Address: "ZxCrkj75u8b218WAiFMiwccwJ6q6GLDLBd4acpHqiyp178qPLQbJR3e1MWHd8hc1h5e4oVWKs9t2VLThQooFidFi1MVAk3kVL",
		Amount:  1,
		AssetID: "",
	})
	res, err := zano.DeployAsset(description, dst)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func Test_emitAsset(t *testing.T) {
	// pass nodeRPC url instead of ""
	zano := NewZanoSDK("http://localhost:10500/json_rpc", "http://37.27.100.59:10505/json_rpc")

	dst := make([]types.Destination, 0)
	dst = append(dst, types.Destination{
		Address: "ZxDphM9gFU359BXfg2BsPi4xrfapivmTi1c1pvogvD3dbAdha4iCosCWup8YkyitrvdAH15Cin65C2AFpA3AF6cJ2amvcNF7w",
		Amount:  100,
		AssetID: "",
	})

	// emit asset returns raw data what can be used to decrypt tx data
	res, err := zano.EmitAsset("42fa39099df8e7ca53140fcce14d9f74c1250a1338b0a9762f677448b9ee08ad", dst)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := zano.TxDetails(res.DataForExternalSigning.OutputsAddresses, res.DataForExternalSigning.UnsignedTx, "", res.DataForExternalSigning.TxSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// sign tx
	signature, err := SignMessage(nil, []byte(setHexPrefix(tx.VerifiedTxID)))
	if err != nil {
		log.Fatal(err)
	}

	strSignature := hexutil.Encode(signature)
	// Remove the last byte (2 hex chars)
	strSignature = strSignature[:len(strSignature)-2]

	// submit signed tx
	submittedTX, err := zano.SendExtSignedAssetTX(strSignature, tx.VerifiedTxID, res.DataForExternalSigning.FinalizedTx, res.DataForExternalSigning.UnsignedTx, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(submittedTX)
}

func Test_BurnAsset(t *testing.T) {
	// pass nodeRPC url instead of ""
	zano := NewZanoSDK("http://localhost:10500/json_rpc", "http://37.27.100.59:10505/json_rpc")

	// emit asset returns raw data what can be used to decrypt tx data
	res, err := zano.BurnAsset("7d3f348fbebfffc4e61a3686189cf870ea393e1c88b8f636acbfdacf9e4b2db2", "1")
	if err != nil {
		log.Fatal(err)
	}

	tx, err := zano.TxDetails(res.DataForExternalSigning.OutputsAddresses, res.DataForExternalSigning.UnsignedTx, "", res.DataForExternalSigning.TxSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// sign tx
	signature, err := SignMessage(nil, []byte(setHexPrefix(tx.VerifiedTxID)))
	if err != nil {
		log.Fatal(err)
	}

	strSignature := hexutil.Encode(signature)
	// Remove the last byte (2 hex chars)
	strSignature = strSignature[:len(strSignature)-2]

	// submit signed tx
	submittedTX, err := zano.SendExtSignedAssetTX(strSignature, tx.VerifiedTxID, res.DataForExternalSigning.FinalizedTx, res.DataForExternalSigning.UnsignedTx, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(submittedTX)
}

func SignMessage(key *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	sig, err := crypto.Sign(message, key)
	if err != nil {
		return nil, err
	}
	//sig[64] += 27

	return sig, nil
}

func setHexPrefix(s string) string {
	return "0x" + s
}
