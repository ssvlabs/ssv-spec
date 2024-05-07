package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
)

var BaseCommitteeWithRunnerSample = func(keySet *TestKeySet, runnerSample *ssv.CommitteeRunner) *ssv.Committee {

	createRunnerF := func() *ssv.CommitteeRunner {
		return ssv.NewCommitteeRunner(runnerSample.BaseRunner.BeaconNetwork,
			runnerSample.BaseRunner.Share,
			qbft.NewController(runnerSample.BaseRunner.QBFTController.Identifier,
				runnerSample.BaseRunner.QBFTController.Share,
				runnerSample.BaseRunner.QBFTController.GetConfig()),
			runnerSample.GetBeaconNode(),
			runnerSample.GetNetwork(),
			runnerSample.GetSigner(),
			runnerSample.GetOperatorSigner(),
			runnerSample.GetValCheckF(),
		).(*ssv.CommitteeRunner)
	}

	return ssv.NewCommittee(
		*TestingOperator(keySet),
		NewTestingVerifier(),
		createRunnerF,
	)
}
