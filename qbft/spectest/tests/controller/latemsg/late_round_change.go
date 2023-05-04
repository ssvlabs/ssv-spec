package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// LateRoundChange tests process late round change msg for an instance which just decided
func LateRoundChange() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := lateRoundChangeStateComparison()

	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.TestingRoundChangeMessage(ks.Shares[4], 4))

	return &tests.ControllerSpecTest{
		Name: "late round change",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
					),
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
	}
}

func lateRoundChangeStateComparison() *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)
	_ = contr.StartNewInstance([]byte{1, 2, 3, 4})

	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.TestingRoundChangeMessage(ks.Shares[4], 4))

	state := testingutils.BaseInstance().State
	state.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))
	state.LastPreparedRound = 1
	state.LastPreparedValue = testingutils.TestingQBFTFullData
	state.Decided = true
	state.DecidedValue = testingutils.TestingQBFTFullData

	state.ProposeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			msgs[0],
		},
	}}
	state.PrepareContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			msgs[1],
			msgs[2],
			msgs[3],
		},
	}}
	state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			msgs[4],
			msgs[5],
			msgs[6],
		},
	}}
	state.RoundChangeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
		qbft.FirstRound: {
			msgs[7],
		},
	}}

	contr.StoredInstances[0].State = state

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
