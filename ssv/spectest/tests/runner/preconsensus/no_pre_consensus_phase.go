package preconsensus

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoPreConsensusPhase tests a valid pre-consensus message but for duties without the pre-consensus phase
func NoPreConsensusPhase() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// MessageID
	msgIDBytes := testingutils.CommitteeMsgID(ks)
	var msgID types.MessageID
	copy(msgID[:], msgIDBytes)

	expectedErr := "no pre consensus phase for committee runner"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "no pre-consensus phase",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "attester",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsg(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), msgID)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.CommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsg(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), msgID)),
				},
				OutputMessages: []*types.PartialSignatureMessages{},
				ExpectedError:  expectedErr,
			},
		},
	}
}
