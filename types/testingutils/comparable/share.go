package testingutilscomparable

import "github.com/bloxapp/ssv-spec/types"

func CompareShare(source, target *types.Share) []Difference {
	ret := make([]Difference, 0)

	if source.OperatorID != target.OperatorID {
		ret = append(ret, Sprintf("OperatorID source/ target %d <---> %d", source.OperatorID, target.OperatorID))
	}

	return ret
}
