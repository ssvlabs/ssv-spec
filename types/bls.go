package types

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"math/big"
	"sync"
)

var (
	curveOrder = new(big.Int)
)

var initOnce = sync.Once{}

// InitBLS initializes BLS
func InitBLS() {
	initOnce.Do(func() {
		_ = bls.Init(bls.BLS12_381)
		_ = bls.SetETHmode(bls.EthModeDraft07)

		curveOrder, _ = curveOrder.SetString(bls.GetCurveOrder(), 10)
	})
}
