package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var BaseAggregatorCommittee = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet) *ssv.Committee {

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
		return AggregatorCommitteeRunnerWithShareMap(shareMap)
	}
	return ssv.NewCommittee(
		*TestingCommitteeMember(keySetSample),
		shareMap,
		createRunnerF,
	)
}

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
		createRunnerF,
	)
}

var BaseAggregatorCommitteeWithCreatorFieldsFromRunner = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, runnerSample *ssv.AggregatorCommitteeRunner) *ssv.Committee {
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
		runner, err := ssv.NewAggregatorCommitteeRunner(
			runnerSample.BaseRunner.BeaconNetwork,
			shareMap,
			qbft.NewController(
				runnerSample.BaseRunner.QBFTController.Identifier,
				runnerSample.BaseRunner.QBFTController.CommitteeMember,
				runnerSample.BaseRunner.QBFTController.GetConfig(),
				TestingOperatorSigner(keySetSample),
			),
			NewTestingBeaconNode(),
			NewTestingNetwork(1, keySetSample.OperatorKeys[1]),
			runnerSample.GetSigner(),
			runnerSample.GetOperatorSigner(),
			runnerSample.GetValCheckF(),
		)
		if err != nil {
			panic(err)
		}
		return runner
	}

	return ssv.NewCommittee(
		*TestingCommitteeMember(keySetSample),
		shareMap,
		createRunnerF,
	)
}
