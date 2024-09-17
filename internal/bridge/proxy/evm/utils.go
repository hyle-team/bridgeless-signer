package evm

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
)

func IsAddressEmpty(addr common.Address) bool {
	return bytes.Equal(addr.Bytes(), new(common.Address).Bytes())
}
