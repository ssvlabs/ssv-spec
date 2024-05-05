package testingutils

import (
	"github.com/ssvlabs/ssv-spec/qbft"
)

func UnknownDutyValueCheck() qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}
