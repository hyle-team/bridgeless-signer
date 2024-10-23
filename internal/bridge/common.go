package bridge

import "regexp"

const (
	HexPrefix                 = "0x"
	DefaultNativeTokenAddress = "0x0000000000000000000000000000000000000000"
)

var DefaultTransactionHashPattern = regexp.MustCompile("^0x[a-fA-F0-9]{64}$")
