package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidRoundChangeJustificationPrepared tests a proposal for > 1 round, prepared previously but one of the round change justifications has validRoundChange != nil
// TODO: make sure it does what it used to do before (https://github.com/ssvlabs/ssv-spec/pull/156#discussion_r1113040758)
func InvalidRoundChangeJustificationPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithParams(
			ks.OperatorKeys[1], types.OperatorID(2), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs),
		),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.OperatorKeys[2], types.OperatorID(2), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs),
		),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.OperatorKeys[3], types.OperatorID(3), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs),
		),
	}

	msgs := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal rc msg invalid (prepared)",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*types.SignedSSVMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: change round msg not valid: msg signature invalid: crypto/rsa: verification error",
	}
}
