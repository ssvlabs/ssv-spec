package testingutilscomparable

import (
	"bytes"
	"github.com/bloxapp/ssv-spec/types"
)

func CompareCommittee(source, target []*types.Operator) []Difference {
	ret := make([]Difference, 0)

	if len(source) != len(target) {
		ret = append(ret, Sprintf("Committee length", "source %d != target %d", len(source), len(target)))
	}

	for i := range source {
		if source[i].OperatorID != target[i].OperatorID {
			ret = append(ret, Sprintf("OperatorID", "source %d != target %d", source[i].OperatorID, target[i].OperatorID))
		}

		if !bytes.Equal(source[i].PubKey, target[i].PubKey) {
			ret = append(ret, Sprintf("PubKey", "source %x != target %x", source[i].PubKey, target[i].PubKey))
		}
	}

	return ret
}
