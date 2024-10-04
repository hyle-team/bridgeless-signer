package operations

import (
	"math/big"
)

func To32Bytes(arr []byte) []byte {
	if len(arr) >= 32 {
		return arr[:32]
	}

	res := make([]byte, 32-len(arr))
	return append(res, arr...)
}

func IntTo32Bytes(amount int) []byte {
	return To32Bytes(big.NewInt(int64(amount)).Bytes())
}

func BoolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	}

	return []byte{0}
}
