package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidPrepareJustificationValue tests a proposal for > 1 round, prepared previously but one of the prepare justifications has value != highest prepared
func InvalidPrepareJustificationValue() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessageWithParams(
			// TODO: different value instead of wrong root
			ks.Shares[3], types.OperatorID(3), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingIdentifier, testingutils.DifferentRoot),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[1], types.OperatorID(1), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[2], types.OperatorID(2), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.Shares[3], types.OperatorID(3), 2,
			testingutils.MarshalJustifications(prepareMsgs)),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "invalid prepare justification value",
		Pre:            pre,
		PostRoot:       "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: change round msg not valid: round change justification invalid: proposed data mistmatch",
	}
}
