package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var BaseCommittee = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet) *ssv.Committee {

	var keySetSample *TestKeySet
	for _, ks := range keySetMap {
		keySetSample = ks
		break
	}

	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	for valIdx, ks := range keySetMap {
		shareMap[valIdx] = TestingShare(ks, valIdx)
	}

	createCommitteeRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
		return CommitteeRunnerWithShareMap(shareMap)
	}

	createAggregatorRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
		return AggregatorCommitteeRunnerWithShareMap(shareMap)
	}

	return ssv.NewCommittee(
		*TestingCommitteeMember(keySetSample),
		shareMap,
		createCommitteeRunnerF,
		createAggregatorRunnerF,
	)
}

var BaseCommitteeWithRunner = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet, runner ssv.Runner) *ssv.Committee {

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
		nil,
	)
}

var BaseCommitteeWithCreatorFieldsFromRunner = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet) *ssv.Committee {

	var keySetSample *TestKeySet
	for _, ks := range keySetMap {
		keySetSample = ks
		break
	}

	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	for valIdx, ks := range keySetMap {
		shareMap[valIdx] = TestingShare(ks, valIdx)
	}

	commRunnerSample := CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)
	createCommitteeRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
		runner, err := ssv.NewCommitteeRunner(
			commRunnerSample.BaseRunner.BeaconNetwork,
			shareMap,
			qbft.NewController(
				commRunnerSample.BaseRunner.QBFTController.Identifier,
				commRunnerSample.BaseRunner.QBFTController.CommitteeMember,
				commRunnerSample.BaseRunner.QBFTController.GetConfig(),
				TestingOperatorSigner(keySetSample),
			),
			NewTestingBeaconNode(),
			NewTestingNetwork(1, keySetSample.OperatorKeys[1]),
			commRunnerSample.GetSigner(),
			commRunnerSample.GetOperatorSigner(),
			commRunnerSample.GetValCheckF(),
		)
		if err != nil {
			panic(err)
		}
		return runner
	}

	aggRunnerSample := AggregatorCommitteeRunnerWithShareMap(shareMap).(*ssv.AggregatorCommitteeRunner)
	createAggCommRunnerF := func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
		runner, err := ssv.NewAggregatorCommitteeRunner(
			aggRunnerSample.BaseRunner.BeaconNetwork,
			shareMap,
			qbft.NewController(
				aggRunnerSample.BaseRunner.QBFTController.Identifier,
				aggRunnerSample.BaseRunner.QBFTController.CommitteeMember,
				aggRunnerSample.BaseRunner.QBFTController.GetConfig(),
				TestingOperatorSigner(keySetSample),
			),
			NewTestingBeaconNode(),
			NewTestingNetwork(1, keySetSample.OperatorKeys[1]),
			aggRunnerSample.GetSigner(),
			aggRunnerSample.GetOperatorSigner(),
			aggRunnerSample.GetValCheckF(),
		)
		if err != nil {
			panic(err)
		}
		return runner
	}

	return ssv.NewCommittee(
		*TestingCommitteeMember(keySetSample),
		shareMap,
		createCommitteeRunnerF,
		createAggCommRunnerF,
	)
}
