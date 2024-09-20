package btc

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	bridgeTypes "github.com/hyle-team/bridgeless-signer/internal/bridge/types"
	"github.com/hyle-team/bridgeless-signer/internal/data"
	"github.com/pkg/errors"
	"math/big"
	"slices"
	"strings"
)

const (
	defaultDecimals                  = 8
	defaultDepositorAddressOutputIdx = 0
)

func (p *proxy) GetDepositData(id data.DepositIdentifier) (*data.DepositData, error) {
	var (
		depositIdx = id.TxEventId
		dstDataIdx = depositIdx + 1
	)

	tx, err := p.getTx(id.TxHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction")
	}

	if tx.BlockHash == "" {
		return nil, bridgeTypes.ErrTxPending
	}
	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode block hash")
	}
	block, err := p.chain.Rpc.GetBlockVerbose(blockHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}
	if tx.Confirmations < p.chain.Confirmations {
		return nil, bridgeTypes.ErrTxNotConfirmed
	}

	if len(tx.Vout) < dstDataIdx+1 || len(tx.Vin) == 0 {
		return nil, bridgeTypes.ErrDepositNotFound
	}

	amount, err := p.parseDepositOutput(tx.Vout[depositIdx])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get deposit amount")
	}

	addr, chainId, err := p.parseDestinationOutput(tx.Vout[dstDataIdx])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get destination address")
	}

	depositor, err := p.parseSenderAddress(tx.Vin[defaultDepositorAddressOutputIdx])
	if err != nil {
		return nil, errors.Wrap(err, "failed to get depositor")
	}

	return &data.DepositData{
		DepositIdentifier:  id,
		DestinationChainId: chainId,
		DestinationAddress: addr,
		SourceAddress:      depositor,
		DepositAmount:      amount,
		// no token address here
		Block: block.Height,
	}, nil
}

func (p *proxy) parseSenderAddress(in btcjson.Vin) (addr string, err error) {
	prevTx, err := p.getTx(in.Txid)
	if err != nil {
		return "", errors.Wrap(err, "failed to get previous transaction to identify sender")
	}

	if len(prevTx.Vout) < int(in.Vout)+1 {
		return "", errors.New("sender vout not found")
	}

	scriptRaw, err := hex.DecodeString(prevTx.Vout[in.Vout].ScriptPubKey.Hex)
	if err != nil {
		return "", errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, err.Error())
	}

	_, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptRaw, p.chain.Params)
	if err != nil {
		return "", errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, err.Error())
	}
	if len(addrs) == 0 {
		return "", errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, "empty sender address")
	}

	return addrs[0].String(), nil
}

func (p *proxy) parseDestinationOutput(out btcjson.Vout) (addr, chainId string, err error) {
	if len(out.ScriptPubKey.Hex) == 0 {
		return addr, chainId, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, "empty destination")
	}

	scriptRaw, err := hex.DecodeString(out.ScriptPubKey.Hex)
	if scriptRaw[0] != txscript.OP_RETURN && len(scriptRaw) <= 3 {
		return addr, chainId, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, "destination data missing")
	}

	// Omitting: OP_RETURN OP_PUSH [return data length] (first three bytes)
	dstData := scriptRaw[3:]

	params := strings.Split(string(dstData), "-")
	if len(params) != 2 {
		return addr, chainId, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, "invalid destination params count")
	}

	return params[0], params[1], nil
}

var supportedScriptTypes = []txscript.ScriptClass{
	txscript.PubKeyHashTy,
	txscript.WitnessV0PubKeyHashTy,
	txscript.WitnessV1TaprootTy,
}

func (p *proxy) parseDepositOutput(out btcjson.Vout) (*big.Int, error) {
	scriptRaw, err := hex.DecodeString(out.ScriptPubKey.Hex)
	if err != nil {
		return nil, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, err.Error())
	}

	stype, addrs, _, err := txscript.ExtractPkScriptAddrs(scriptRaw, p.chain.Params)
	if err != nil {
		return nil, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, err.Error())
	}
	if !slices.Contains(supportedScriptTypes, stype) || len(addrs) != 1 {
		return nil, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, fmt.Sprintf("unsupported type %s", stype))
	}
	if !p.bridgeAddr(addrs[0]) {
		return nil, errors.Wrap(bridgeTypes.ErrInvalidScriptPubKey, "receiver address is not bridge")
	}
	if out.Value == 0 {
		return nil, bridgeTypes.ErrInvalidDepositedAmount
	}

	return toBigint(out.Value, defaultDecimals), nil
}

func (p *proxy) bridgeAddr(addr btcutil.Address) bool {
	for _, receiver := range p.chain.Receivers {
		if addr.String() == receiver.String() {
			return true
		}
	}

	return false
}

func toBigint(val float64, decimals int64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	coin.SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result)

	return result
}
