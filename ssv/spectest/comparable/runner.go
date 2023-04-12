package comparable

import (
	"bytes"
	comparable2 "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/ssv"
	testingutilscomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

func CompareRunner(source, target ssv.Runner) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if diff := testingutilscomparable.NestedCompare("BaseRunner", compareBaseRunner(source.GetBaseRunner(), target.GetBaseRunner())); len(diff) > 0 {
		ret = append(ret, diff)
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

func compareBaseRunner(source, target *ssv.BaseRunner) []testingutilscomparable.Difference {
	ret := make([]testingutilscomparable.Difference, 0)

	if shareDiff := testingutilscomparable.NestedCompare("State", compareBaseRunnerState(source.State, target.State)); len(shareDiff) > 0 {
		ret = append(ret, shareDiff)
	}

	if shareDiff := testingutilscomparable.NestedCompare("Share", testingutilscomparable.CompareShare(source.Share, target.Share)); len(shareDiff) > 0 {
		ret = append(ret, shareDiff)
	}

	if shareDiff := testingutilscomparable.NestedCompare("QBFTController", comparable2.CompareQBFTController(source.QBFTController, target.QBFTController)); len(shareDiff) > 0 {
		ret = append(ret, shareDiff)
	}

	if source.BeaconNetwork != target.BeaconNetwork {
		ret = append(ret, testingutilscomparable.Sprintf("BeaconNetwork", "source %d != target %d", source.BeaconNetwork, target.BeaconNetwork))
	}

	if source.BeaconRoleType != target.BeaconRoleType {
		ret = append(ret, testingutilscomparable.Sprintf("BeaconRoleType", "source %d != target %d", source.BeaconRoleType, target.BeaconRoleType))
	}

	return ret
}
