package testingutils

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
)

func UnknownDutyValueCheck() qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}
