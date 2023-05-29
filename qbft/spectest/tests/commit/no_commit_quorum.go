package commit

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// NoCommitQuorum tests the state of the QBFT instance when received commit messages don't create a quorum
func NoCommitQuorum() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	sc := NoCommitQuorumStateComparison()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMessage(ks.Shares[2], 2),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "no commit quorum",
		Pre:           pre,
		PostRoot:      sc.Root(),
		PostState:     sc.ExpectedState,
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
	}
}

func NoCommitQuorumStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()

	state := testingutils.BaseInstance().State
	state.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	state.LastPreparedRound = 1
	state.LastPreparedValue = testingutils.TestingQBFTFullData
	state.Decided = false

	state.ProposeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
		},
	}}

	state.PrepareContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
			testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
		},
	}}

	state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[2], types.OperatorID(2)),
		},
	}}

	return &comparable.StateComparison{ExpectedState: state}
}
