package operations

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"math/big"
)

func ToBytes32(arr []byte) []byte {
	if len(arr) >= 32 {
		return arr[:32]
	}

	res := make([]byte, 32-len(arr))
	return append(res, arr...)
}

func IntToBytes32(amount int) []byte {
	return ToBytes32(big.NewInt(int64(amount)).Bytes())
}

func BoolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	}

	return []byte{0}
}

func HexToBytes32(val string) ([]byte, error) {
	if Has0xPrefix(val) {
		val = val[2:]
	}

	// each byte is encoded by two hexadecimal digits
	if len(val) != 32*2 {
		return nil, errors.New("invalid hex string length")
	}

	return hex.DecodeString(val)
}

func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
