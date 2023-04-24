package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// FullFlowAfterDecided tests a decided msg for round 1 followed by a full proposal, prepare, commit for round 2
func FullFlowAfterDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMultiSignerMessage(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
			[]types.OperatorID{1, 2, 3},
		),

		testingutils.TestingProposalMessageWithParams(
			ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), nil,
		),

		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),

		testingutils.TestingCommitMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingCommitMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingCommitMessageWithRound(ks.Shares[4], types.OperatorID(4), 2),
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
				ControllerPostRoot: "d8ad104b049b410a9c6462af1fce5ff279f04f0820ac79116d5713e3525ab3ca",
			},
		},
	}
}
