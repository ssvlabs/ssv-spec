package prepare

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// PrepareQuorumTriggeredTwice tests triggering prepare quorum twice by sending > 2f+1 prepare messages
func PrepareQuorumTriggeredTwice() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	sc := prepareQuorumTriggeredTwiceStateComparison()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[4], 4),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepared quorum committed twice",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
			testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		},
	}
}

func prepareQuorumTriggeredTwiceStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	state := testingutils.BaseInstance().State
	state.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	state.LastPreparedRound = 1
	state.LastPreparedValue = testingutils.TestingQBFTFullData

	state.ProposeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
		qbft.FirstRound: {
			testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))),
		},
	}}

	state.PrepareContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
		qbft.FirstRound: {
			testingutils.ToProcessingMessage(testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1))),
			testingutils.ToProcessingMessage(testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2))),
			testingutils.ToProcessingMessage(testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3))),
			testingutils.ToProcessingMessage(testingutils.TestingPrepareMessage(ks.OperatorKeys[4], types.OperatorID(4))),
		},
	}}

	state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
		qbft.FirstRound: {
			testingutils.ToProcessingMessage(testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1))),
		},
	}}

	return &comparable.StateComparison{ExpectedState: state}
}
