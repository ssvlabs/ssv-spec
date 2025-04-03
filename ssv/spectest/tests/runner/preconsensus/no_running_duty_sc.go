package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv"

	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// noRunningDutySyncCommitteeContributionSC returns state comparison object for the NoRunningDuty SyncCommitteeContribution versioned spec test
func noRunningDutySyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			return testingutils.SyncCommitteeContributionRunner(ks)
		}(),
	}
}

// noRunningDutyValidatorRegistrationSC returns state comparison object for the NoRunningDuty ValidatorRegistration versioned spec test
func noRunningDutyValidatorRegistrationSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			return testingutils.ValidatorRegistrationRunner(ks)
		}(),
	}
}

// noRunningDutyVoluntaryExitSC returns state comparison object for the NoRunningDuty VoluntaryExit versioned spec test
func noRunningDutyVoluntaryExitSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			return testingutils.VoluntaryExitRunner(ks)
		}(),
	}
}

// noRunningDutyProposerSC returns state comparison object for the NoRunningDuty Proposer versioned spec test
func noRunningDutyProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			return testingutils.ProposerRunner(ks)
		}(),
	}
}

// noRunningDutyBlindedProposerSC returns state comparison object for the NoRunningDuty Blinded Proposer versioned spec test
func noRunningDutyBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			return testingutils.ProposerBlindedBlockRunner(ks)
		}(),
	}
}
