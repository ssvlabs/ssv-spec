package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// noRunningDutySyncCommitteeContributionSC returns state comparison object for the NoRunningDuty SyncCommitteeContribution versioned spec test
func noRunningDutySyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			return testingutils.SyncCommitteeContributionRunner(ks)
		}(),
	}
}

// noRunningDutyAggregatorSC returns state comparison object for the NoRunningDuty Aggregator versioned spec test
func noRunningDutyAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			return testingutils.AggregatorRunner(ks)
		}(),
	}
}

// noRunningDutyValidatorRegistrationSC returns state comparison object for the NoRunningDuty ValidatorRegistration versioned spec test
func noRunningDutyValidatorRegistrationSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			return testingutils.ValidatorRegistrationRunner(ks)
		}(),
	}
}

// noRunningDutyProposerSC returns state comparison object for the NoRunningDuty Proposer versioned spec test
func noRunningDutyProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			return testingutils.ProposerRunner(ks)
		}(),
	}
}

// noRunningDutyBlindedProposerSC returns state comparison object for the NoRunningDuty Blinded Proposer versioned spec test
func noRunningDutyBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	return &comparable.StateComparison{
		ExpectedState: func() types.Root {
			return testingutils.ProposerBlindedBlockRunner(ks)
		}(),
	}
}
