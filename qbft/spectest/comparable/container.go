package comparable

import (
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func CompareQBFTMessageContainer(source, target *qbft.MsgContainer) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if len(source.Msgs) != len(target.Msgs) {
		ret = append(ret, testingutilscomparable.Sprintf("Committee length source/ target %d <---> %d", len(source.Msgs), len(target.Msgs)))
		return ret
	}

	for round := range source.Msgs {
		if len(source.Msgs[round]) != len(target.Msgs[round]) {
			ret = append(ret, testingutilscomparable.Sprintf("Round %d length source/ target %d <---> %d", round, len(source.Msgs), len(target.Msgs)))
			return ret
		}
		for i := range source.Msgs[round] {
			if diff := testingutilscomparable.NestedCompare(
				fmt.Sprintf("Round %d", round),
				CompareQBFTSignedMessage(source.Msgs[round][i], target.Msgs[round][i])); len(diff) > 0 {
				ret = append(ret, diff)
			}
		}
	}

	return ret
}
