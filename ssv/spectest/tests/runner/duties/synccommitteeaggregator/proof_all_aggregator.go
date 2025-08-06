package synccommitteeaggregator

import (
	"encoding/hex"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// AllAggregatorQuorum tests a quorum of selection proofs of which all are aggregator
func AllAggregatorQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	test := NewSyncCommitteeAggregatorProofSpecTest(
		"sync committee aggregator all are aggregators",
		testdoc.SyncCommitteeAggregatorProofAllAggregatorDoc,
		[]*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
		},
		"b11cb2d6684fd612da4b885cd52080cd2d5a2c79f5a42297623205603f1f6599",
		"",
		map[string]bool{
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[0][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[1][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[2][:]): true,
		},
		"",
		ks,
	)

	return test
}
