package latemsg

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// LateRoundChangePastRound tests process late round change msg for an instance which just decided for a round < decided round
func LateRoundChangePastRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessage(ks.NetworkKeys[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.NetworkKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.NetworkKeys[2], types.OperatorID(2)),
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs, []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.NetworkKeys[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(rcMsgs)),

		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[2], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),

		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),

		testingutils.TestingRoundChangeMessage(ks.NetworkKeys[4], types.OperatorID(4)),
	}...)

	return &tests.ControllerSpecTest{
		Name:          "late round change past round",
		ExpectedError: "could not process msg: invalid signed message: past round",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithRound(
						[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
						[]types.OperatorID{1, 2, 3},
						2,
					),
				},
				ControllerPostRoot: "94ff35842a25bfb2379bb5419e83ad32f4e4cc079a2b8cfa6f6e183fd44b8f52",
			},
		},
	}
}
