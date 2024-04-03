package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// PrepareQuorumTriggeredTwice tests triggering prepare quorum twice by sending > 2f+1 prepare messages
func PrepareQuorumTriggeredTwice() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance()
	sc := prepareQuorumTriggeredTwiceStateComparison()
	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[3], 3),

		testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[4], 4),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepared quorum committed twice",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], 1),
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], 1),
		},
	}
}

func prepareQuorumTriggeredTwiceStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	state := testingutils.BaseInstance().State
	state.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1))

	state.LastPreparedRound = 1
	state.LastPreparedValue = testingutils.TestingQBFTFullData

	state.ProposeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*types.SignedSSVMessage{
		qbft.FirstRound: {
			testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		},
	}}

	state.PrepareContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*types.SignedSSVMessage{
		qbft.FirstRound: {
			testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.NetworkKeys[3], types.OperatorID(3)),
			testingutils.TestingPrepareMessage(ks.NetworkKeys[4], types.OperatorID(4)),
		},
	}}

	state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*types.SignedSSVMessage{
		qbft.FirstRound: {
			testingutils.TestingCommitMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		},
	}}

	return &comparable.StateComparison{ExpectedState: state}
}
