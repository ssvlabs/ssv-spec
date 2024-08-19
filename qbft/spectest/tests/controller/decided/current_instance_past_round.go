package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CurrentInstancePastRound tests a decided msg received for current running instance for a past round
func CurrentInstancePastRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[1], 1, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], 2, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], 3, 2),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], 1),

		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], 1),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], 2),
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs, []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.OperatorKeys[1], 1, 2,
			testingutils.MarshalJustifications(rcMsgs)),

		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[1], 1, 2),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[2], 2, 2),
		testingutils.TestingPrepareMessageWithRound(ks.OperatorKeys[3], 3, 2),

		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[1], 1, 2),
		testingutils.TestingCommitMessageWithRound(ks.OperatorKeys[2], 2, 2),
		testingutils.TestingCommitMultiSignerMessage([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}),
	}...)

	return &tests.ControllerSpecTest{
		Name: "decide current instance past round",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedCnt: 1,
					DecidedVal: testingutils.TestingQBFTFullData,
				},
			},
		},
	}
}
