package evm

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	zeroAddress = "0x0000000000000000000000000000000000000000"
)

func IsAddressEmpty(addr common.Address) bool {
	if bytes.Compare(addr.Bytes(), hexutil.MustDecode(zeroAddress)) != 0 {
		return false
	}

	return true
}
