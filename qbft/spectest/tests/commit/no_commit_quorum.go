package commit

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// NoCommitQuorum tests the state of the QBFT instance when received commit messages don't create a quorum
func NoCommitQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := NoCommitQuorumStateComparison()

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], 3),

		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[2], 2),
	}

	outputMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingCommitMessage(ks.OperatorKeys[1], 1),
	}

	test := tests.NewMsgProcessingSpecTest(
		"no commit quorum",
		testdoc.CommitTestNoCommitQuorumDoc,
		pre,
		sc.Root(),
		sc.ExpectedState,
		msgs,
		outputMsgs,
		0,
		nil,
		ks,
	)

	return test
}

func NoCommitQuorumStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	state := testingutils.BaseInstance().State
	state.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))

	state.LastPreparedRound = 1
	state.LastPreparedValue = testingutils.TestingQBFTFullData
	state.Decided = false

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
		},
	}}

	state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
		qbft.FirstRound: {
			testingutils.ToProcessingMessage(testingutils.TestingCommitMessage(ks.OperatorKeys[1], types.OperatorID(1))),
			testingutils.ToProcessingMessage(testingutils.TestingCommitMessage(ks.OperatorKeys[2], types.OperatorID(2))),
		},
	}}

	return &comparable.StateComparison{ExpectedState: state}
}
