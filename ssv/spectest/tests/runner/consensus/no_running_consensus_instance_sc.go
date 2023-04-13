package consensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	ssvcomparable "github.com/bloxapp/ssv-spec/ssv/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// noRunningConsensusInstanceStateComparison returns state comparison object for the NoRunningConsensusInstance spec test
func noRunningConsensusInstanceStateComparison() *ssvcomparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	return &ssvcomparable.StateComparison{
		SyncCommitteeContribution: func() ssv.Runner {
			return testingutils.SyncCommitteeContributionRunner(ks)
		}(),
		SyncCommittee: func() ssv.Runner {
			return testingutils.SyncCommitteeRunner(ks)
		}(),
		Aggregator: func() ssv.Runner {
			return testingutils.AggregatorRunner(ks)
		}(),
		Proposer: func() ssv.Runner {
			return testingutils.ProposerRunner(ks)
		}(),
		BlindedProposer: func() ssv.Runner {
			return testingutils.ProposerBlindedBlockRunner(ks)
		}(),
		Attester: func() ssv.Runner {
			return testingutils.AttesterRunner(ks)
		}(),
	}
}
