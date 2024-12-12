package lib

import (
	"fmt"
	"math/big"
)

func BigIntFromStr(s string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(s, 10)

	if !ok {
		return nil, fmt.Errorf("error converting to bigint: %v", s)
	}

	return n, nil
}
