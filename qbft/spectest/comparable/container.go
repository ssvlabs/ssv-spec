package comparable

import (
	"github.com/bloxapp/ssv-spec/qbft"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func CompareQBFTMessageContainer(source, target *qbft.MsgContainer) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)
	return ret
}
