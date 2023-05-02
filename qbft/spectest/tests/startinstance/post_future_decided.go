package startinstance

import (
	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreviousDecided tests starting an instance when the previous one decided
func PreviousDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "start instance prev decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
					testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks),
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
				},
				ControllerPostRoot:  previousDecidedStateComparison(qbft.FirstHeight, true).Root(),
				ControllerPostState: previousDecidedStateComparison(qbft.FirstHeight, true).ExpectedState,
			},
			{
				InputValue:          []byte{1, 2, 3, 4},
				ControllerPostRoot:  previousDecidedStateComparison(1, false).Root(),
				ControllerPostState: previousDecidedStateComparison(1, false).ExpectedState,
			},
		},
	}
}

func previousDecidedStateComparison(height qbft.Height, decidedState bool) *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)

	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	for i := 0; i <= int(height); i++ {
		contr.Height = qbft.Height(i)
		_ = contr.StartNewInstance([]byte{1, 2, 3, 4})

		state := testingutils.BaseInstance().State
		state.Height = qbft.Height(i)

		// last height
		if !decidedState && i == int(height) {
			contr.StoredInstances[0].State = state
			break
		}

		state.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.Height(i), testingutils.TestingQBFTRootData, nil, nil)
		state.ProposalAcceptedForCurrentRound.Message.Height = qbft.Height(i)
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

		contr.StoredInstances[0].State = state
	}

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
