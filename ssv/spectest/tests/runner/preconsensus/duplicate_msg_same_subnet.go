package preconsensus

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DuplicateMsgSameSubnet ensures only one selection proof per subnet is produced
func DuplicateMsgSameSubnet() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// Subnet is computed as index / (512/4):
	// 		- 0, 0, 1 will return 0
	// 		- 129, 129, 130 will return 1
	// Final message should only have 2 (root, selection proof) pairs, not 4 or 6
	// List with 6 indices, but only 2 unique subnets
	syncCommitteeIndices := []types.ValidatorSyncCommitteeIndex{0, 0, 1, 129, 129, 130}

	// Duty with indices 6 indices
	duty := testingutils.TestingSyncCommitteeContributionDutyWithValidatorContributionIndices(syncCommitteeIndices)

	// Message built with only 2 indices (one per subnet), as it should be the same as if we built it with all 6 indices!
	equivalentSyncCommitteeIndices := []types.ValidatorSyncCommitteeIndex{1, 130}
	outMsg := testingutils.PreConsensusContributionProofWithValidatorSyncCommitteeIndices(ks.Shares[1], ks.Shares[1], 1, 1, equivalentSyncCommitteeIndices)

	test := tests.NewMultiMsgProcessingSpecTest(
		"pre consensus duplicate selection proofs same subnet",
		testdoc.PreConsensusDuplicatedContributionSubnetDoc,
		[]*tests.MsgProcessingSpecTest{
			{
				Name:     "sync committee aggregator selection proof same subnet",
				Runner:   testingutils.AggregatorCommitteeRunner(ks),
				Duty:     duty,
				Messages: []*types.SignedSSVMessage{},
				OutputMessages: []*types.PartialSignatureMessages{
					outMsg,
				},
			},
		},
		ks,
	)

	return test
}
