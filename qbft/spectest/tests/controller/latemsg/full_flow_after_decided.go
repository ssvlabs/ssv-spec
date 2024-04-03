package latemsg

import (
	"crypto/rsa"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FullFlowAfterDecided tests a decided msg for round 1 followed by a full proposal, prepare, commit for round 2
func FullFlowAfterDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingCommitMultiSignerMessage(
			[]*rsa.PrivateKey{ks.NetworkKeys[1], ks.NetworkKeys[2], ks.NetworkKeys[3]},
			[]types.OperatorID{1, 2, 3},
		),

		testingutils.TestingProposalMessageWithParams(
			ks.NetworkKeys[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),

		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.NetworkKeys[3], types.OperatorID(3), 2),

		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[1], types.OperatorID(1), 2),
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[2], types.OperatorID(2), 2),
		testingutils.TestingCommitMessageWithRound(ks.NetworkKeys[4], types.OperatorID(4), 2),
	}

	return &tests.ControllerSpecTest{
		Name: "full flow after decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
				},
				ControllerPostRoot: "ac8cbc7f074702dd643312cf277d894b35079706116237e3eb41a82914ee69e9",
			},
		},
	}
}
