package attestationsynccommittee

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidConsensusData sends a proposal message to a maximal committee runner with a valid BeaconVote consensus data
func ValidConsensusData() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	multiSpecTest := &tests.MultiMsgProcessingSpecTest{
		Name: "valid consensus data",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "500 attestations 500 sync committees",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, validatorIndexList(500), validatorIndexList(500)),
				Messages: []*types.SignedSSVMessage{
					testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.OperatorKeys[1], types.OperatorID(1), testingutils.CommitteeMsgID(ks), testingutils.TestBeaconVoteByts,
						qbft.Height(testingutils.TestingDutySlot)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
			},
		},
	}

	return multiSpecTest
}
