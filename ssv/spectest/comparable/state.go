package comparable

import (
	"bytes"
	comparable2 "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/ssv"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func compareBaseRunnerState(source, target *ssv.State) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if (source == nil && target != nil) || (source != nil && target == nil) {
		ret = append(ret, testingutilscomparable.Sprintf("State nil?", "source %t != target %t", source == nil, target == nil))
		return ret
	}

	if diff := testingutilscomparable.NestedCompare("PreConsensusContainer", compareRunnerPartialSigContainer(source.PreConsensusContainer, target.PreConsensusContainer)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("PostConsensusContainer", compareRunnerPartialSigContainer(source.PostConsensusContainer, target.PostConsensusContainer)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("RunningInstance", comparable2.CompareQBFTInstance(source.RunningInstance, target.RunningInstance)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("DecidedValue", testingutilscomparable.CompareConsensusData(source.DecidedValue, target.DecidedValue)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if diff := testingutilscomparable.NestedCompare("StartingDuty", testingutilscomparable.CompareDuty(source.StartingDuty, target.StartingDuty)); len(diff) > 0 {
		ret = append(ret, diff)
	}

	if source.Finished != target.Finished {
		ret = append(ret, testingutilscomparable.Sprintf("Finished", "source %t != target %t", source.Finished, target.Finished))
	}

	r1, err := source.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	r2, err := target.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	if !bytes.Equal(r1[:], r2[:]) {
		ret = append(ret, testingutilscomparable.Sprintf("Root", "source %x != target %x", r1, r2))
	}

	return ret
}
