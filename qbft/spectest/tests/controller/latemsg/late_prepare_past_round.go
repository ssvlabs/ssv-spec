package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LatePreparePastRound tests process late prepare msg for an instance which just decided for a round < decided round
func LatePreparePastRound() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),

		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs, []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(rcMsgs)),

		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),

		testingutils.TestingCommitMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingCommitMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingCommitMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),

		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
	}...)

	return &tests.ControllerSpecTest{
		Name: "late prepare past round",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithRound(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
						2,
					),
				},
				ControllerPostRoot: "04eba8412e564e9b000f61035cfbc663216671d98d39d74119612614916542ad",
			},
		},
		ExpectedError: "could not process msg: invalid signed message: past round",
	}
}
