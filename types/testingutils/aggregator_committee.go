package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var BaseAggregatorCommitteeWithRunner = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, runner ssv.Runner) *ssv.Committee {

	var keySetSample *TestKeySet
	for _, ks := range keySetMap {
		keySetSample = ks
		break
	}

	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	for valIdx, ks := range keySetMap {
		shareMap[valIdx] = TestingShare(ks, valIdx)
	}

	createRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
		return runner
	}

	return ssv.NewCommittee(
		*TestingCommitteeMember(keySetSample),
		shareMap,
		nil,
		createRunnerF,
	)
}
