package preconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

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
