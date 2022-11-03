package synccommitteeaggregator

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// AllAggregatorQuorum tests a quorum of selection proofs of which all are aggregator
func AllAggregatorQuorum() *SyncCommitteeAggregatorProofSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &SyncCommitteeAggregatorProofSpecTest{
		Name: "sync committee aggregator all are aggregators",
		Messages: []*types.Message{
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), types.PartialContributionProofSignatureMsgType),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2), types.PartialContributionProofSignatureMsgType),
			testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3), types.PartialContributionProofSignatureMsgType),
		},
		ProofRootsMap: map[string]bool{
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[0][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[1][:]): true,
			hex.EncodeToString(testingutils.TestingContributionProofsSigned[2][:]): true,
		},
		PostDutyRunnerStateRoot: "82af835fc5163cde0feca2992f215e722e7bf96b7b9622d4256ad7ed86a3e3ce",
	}
}
