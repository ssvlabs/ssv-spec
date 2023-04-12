package comparable

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/qbft"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func CompareQBFTInstance(source, target *qbft.Instance) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if (source == nil && target != nil) || (source != nil && target == nil) {
		ret = append(ret, testingutilscomparable.Sprintf("Instance nil?", "source %t != target %t", source == nil, target == nil))
		return ret
	}

	if diff := testingutilscomparable.NestedCompare("State", CompareQBFTState(source.State, target.State)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if !bytes.Equal(source.StartValue, target.StartValue) {
		ret = append(ret, testingutilscomparable.Sprintf("StartValue", "source %x != target %x", source.StartValue, target.StartValue))
	}

	return ret
}
