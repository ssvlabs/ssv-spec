package synccommitteeaggregator

import (
	"encoding/hex"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoneAggregatorQuorum tests a quorum of selection proofs but none of which are aggregator
func NoneAggregatorQuorum() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	test := NewSyncCommitteeAggregatorProofSpecTest(
		"sync committee aggregator none is aggregator",
		testdoc.SyncCommitteeAggregatorProofNoneAggregatorDoc,
		[]*types.SignedSSVMessage{
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
			testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),
		},
		"05c705130d1fdb9401cc21dccfd35d21eeb4a5d541ff96af2b9c908d3f646100",
		"",
		map[string]bool{
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[0][:]): false,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[1][:]): false,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[2][:]): false,
		},
		"",
		ks,
	)

	return test
}
