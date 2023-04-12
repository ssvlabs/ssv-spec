package comparable

import (
	"bytes"
	"fmt"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func compareNestedSigMap(source, target map[types.OperatorID][]byte) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	for operatorID := range source {
		if !bytes.Equal(source[operatorID], target[operatorID]) {
			ret = append(ret, testingutilscomparable.Sprintf(fmt.Sprintf("OperatorID %d Data", operatorID), "source %x != target %x", source[operatorID], target[operatorID]))
		}
	}

	return ret
}

func compareRunnerPartialSigContainer(source, target *ssv.PartialSigContainer) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if (source == nil && target != nil) || (source != nil && target == nil) {
		ret = append(ret, testingutilscomparable.Sprintf("Container nil?", "source %t != target %t", source == nil, target == nil))
		return ret
	}

	if source.Quorum != target.Quorum {
		ret = append(ret, testingutilscomparable.Sprintf("Quorum", "source %t != target %t", source.Quorum, target.Quorum))
	}

	if len(source.Signatures) != len(target.Signatures) {
		ret = append(ret, testingutilscomparable.Sprintf("Signatures", "length source %d != target %d", len(source.Signatures), len(target.Signatures)))
		return ret
	}

	for root := range source.Signatures {
		if len(source.Signatures[root]) != len(target.Signatures[root]) {
			ret = append(ret, testingutilscomparable.Sprintf(fmt.Sprintf("Root %s", root), "source %d != target %d", len(source.Signatures[root]), len(target.Signatures[root])))
			continue
		}
		if diff := testingutilscomparable.NestedCompare(
			fmt.Sprintf("Root %s", root),
			compareNestedSigMap(source.Signatures[root], target.Signatures[root])); len(diff) > 0 {
			ret = append(ret, diff)
		}
	}

	return ret
}
