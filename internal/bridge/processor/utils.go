package processor

import "math/big"

func transformAmount(amount *big.Int, currentDecimals uint64, targetDecimals uint64) {
	if currentDecimals == targetDecimals {
		return
	}

	if currentDecimals < targetDecimals {
		for i := uint64(0); i < targetDecimals-currentDecimals; i++ {
			amount.Mul(amount, new(big.Int).SetInt64(10))
		}
		return
	}

	for i := uint64(0); i < currentDecimals-targetDecimals; i++ {
		amount.Div(amount, new(big.Int).SetInt64(10))
	}

	return
}
