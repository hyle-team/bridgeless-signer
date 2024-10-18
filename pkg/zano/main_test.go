package zano

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyle-team/bridgeless-signer/pkg/zano/types"
	"github.com/pkg/errors"
	"log"
	"testing"
	"time"
)

const (
	firstWallet  = "http://localhost:10500/json_rpc"
	secondWallet = "http://localhost:10505/json_rpc"
	secondKey    = ""
	firstKey     = ""
)

func Test_sign(t *testing.T) {
	// sign tx
	kp, err := crypto.HexToECDSA(secondKey)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Priv"))
	}
	const resultSig = "0x0ffc63e5113ee4d8da4262f1d9f4a0c6d6340e7042070c5e72bb748430ce17b456dcf27bf375be49a0071f9234778123e6dbb4f8ffd65f2989b5593864a2c0b41c"

	raw := hexutil.MustDecode("0x70dd03634d73880375109e0e6a57fb2769f83562d229646b32ab8a8362d932bf")
	signature, err := SignMessage(kp, raw)
	if err != nil {
		log.Fatal(3, err)
	}

	strSignature := hexutil.Encode(signature)
	fmt.Println(strSignature)
	fmt.Println(strSignature == resultSig)
}

func Test_transfer(t *testing.T) {
	zano := NewSDK(secondWallet, "")
	dst := make([]types.Destination, 0)
	dst = append(dst, types.Destination{
		Address: "ZxCuASh1nm3PzJThzA5fJf6BpnpBVZCj8iyKYEtnMsVBXzjxrvuh5X9TmWzDxezSPKjJLzAscgtgFWNJRKMT2WZL16fiGvfdm",
		Amount:  123,
		AssetID: "d6329b5b1f7c0805b5c345f4957554002a2f557845f64d7645dae0e051a6498a",
	})

	res, err := zano.Transfer("second creation", nil, dst)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.TxHash)

	txs, err := zano.GetTransactions(res.TxHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(txs.In))
	fmt.Println(len(txs.Out))
	fmt.Println(len(txs.Pool))

	time.Sleep(time.Second * 2)

	txs, err = zano.GetTransactions(res.TxHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(txs.In))
	fmt.Println(len(txs.Out))
	fmt.Println(len(txs.Pool))
}

func Test_getTx(t *testing.T) {
	zano := NewSDK(firstWallet, "")

	res, err := zano.GetTransactions("299fb99efaec9feb5a7bc8d4c2f81d193386192b27eed3a03e68b3d2e7ad2a7c")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(res.In))
	fmt.Println(len(res.Out))
	fmt.Println(len(res.Pool))
}

func Test_deployAsset(t *testing.T) {
	zano := NewSDK(secondWallet, "")

	description := types.AssetDescriptor{
		DecimalPoint:   12,
		FullName:       "test asset second",
		HiddenSupply:   false,
		MetaInfo:       "TESTS",
		Owner:          "",
		Ticker:         "TESTS",
		TotalMaxSupply: "1000000000000000000",
		CurrentSupply:  "1000000000000",
		OwnerEthPubKey: "035288e584e37f7479355f69f968ec88cc617947834eea63e94db535c8e717d774",
	}

	dst := make([]types.Destination, 0)
	dst = append(dst, types.Destination{
		Address: "ZxCuASh1nm3PzJThzA5fJf6BpnpBVZCj8iyKYEtnMsVBXzjxrvuh5X9TmWzDxezSPKjJLzAscgtgFWNJRKMT2WZL16fiGvfdm",
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
	zano := NewSDK(secondWallet, "http://37.27.100.59:10505/json_rpc")

	dst := []types.Destination{{
		Address: "ZxDphM9gFU359BXfg2BsPi4xrfapivmTi1c1pvogvD3dbAdha4iCosCWup8YkyitrvdAH15Cin65C2AFpA3AF6cJ2amvcNF7w",
		Amount:  15151515,
		AssetID: "",
	}}

	// emit asset returns raw data what can be used to decrypt tx data
	res, err := zano.EmitAsset("cab92cb5338d7b9f533c404c884cadb3ba579601074a6216e28f0d4da13e2c14", dst...)
	if err != nil {
		log.Fatal(1, err)
	}

	tx, err := zano.TxDetails(res.DataForExternalSigning.OutputsAddresses, res.DataForExternalSigning.UnsignedTx, "", res.DataForExternalSigning.TxSecretKey)
	if err != nil {
		log.Fatal(2, err)
	}

	// sign tx
	kp, err := crypto.HexToECDSA(secondKey)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Sign"))
	}
	// sign tx
	signature, err := SignMessage(kp, hexutil.MustDecode(setHexPrefix(tx.VerifiedTxID)))
	if err != nil {
		log.Fatal(3, err)
	}

	strSignature := hexutil.Encode(signature)
	// Remove the last byte (2 hex chars)
	strSignature = strSignature[2 : len(strSignature)-2]
	// submit signed tx
	submittedTX, err := zano.SendExtSignedAssetTX(strSignature, tx.VerifiedTxID, res.DataForExternalSigning.FinalizedTx, res.DataForExternalSigning.UnsignedTx, false)
	if err != nil {
		log.Fatal(4, err)
	}

	log.Println(submittedTX)
}

func Test_BurnAsset(t *testing.T) {
	// pass nodeRPC url instead of ""
	zano := NewSDK(secondWallet, "http://37.27.100.59:10505/json_rpc")

	// emit asset returns raw data what can be used to decrypt tx data
	res, err := zano.BurnAsset("7d3f348fbebfffc4e61a3686189cf870ea393e1c88b8f636acbfdacf9e4b2db2", "123")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Burn"))
	}

	tx, err := zano.TxDetails(res.DataForExternalSigning.OutputsAddresses, res.DataForExternalSigning.UnsignedTx, "", res.DataForExternalSigning.TxSecretKey)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Details"))
	}

	// sign tx
	kp, err := crypto.HexToECDSA(secondKey)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Sign"))
	}

	signature, err := SignMessage(kp, []byte(setHexPrefix(tx.VerifiedTxID)))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Sign"))
	}

	strSignature := hexutil.Encode(signature)
	// Remove the last byte (2 hex chars)
	strSignature = strSignature[:len(strSignature)-2]

	// submit signed tx
	submittedTX, err := zano.SendExtSignedAssetTX(strSignature, tx.VerifiedTxID, res.DataForExternalSigning.FinalizedTx, res.DataForExternalSigning.UnsignedTx, false)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Send"))
	}

	log.Println(submittedTX)
}

func SignMessage(key *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	sig, err := crypto.Sign(message, key)
	if err != nil {
		return nil, err
	}
	sig[64] += 27

	return sig, nil
}

func setHexPrefix(s string) string {
	return "0x" + s
}
