package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PrepareJustificationNestedJustification tests a proposal for > 1 round, prepared previously, where a prepare
// message in its prepare justifications itself carries prepare justifications
func PrepareJustificationNestedJustification() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	valid, nested := prepareQuorums(ks)

	return nestedPrepareJustificationTest(
		"proposal prepare justification nested justification",
		testdoc.ProposalPrepareJustificationNestedJustificationDoc,
		ks, valid, nested,
		types.PrepareMessageInvalidErrorCode,
	)
}

// RCJustificationNestedPrepareJustification tests a proposal for > 1 round, prepared previously, where a prepare
// message in the justification of an embedded round change itself carries prepare justifications
func RCJustificationNestedPrepareJustification() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	valid, nested := prepareQuorums(ks)

	return nestedPrepareJustificationTest(
		"proposal rc justification nested prepare justification",
		testdoc.ProposalRCJustificationNestedPrepareJustificationDoc,
		ks, nested, valid,
		types.UnexpectedJustificationErrorCode,
	)
}

// prepareQuorums returns a valid prepare quorum, and the same quorum where operator 3's prepare carries a prepare
// justification
func prepareQuorums(ks *testingutils.TestKeySet) (valid, nested []*types.SignedSSVMessage) {
	valid = []*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	}
	nested = []*types.SignedSSVMessage{
		valid[0],
		valid[1],
		testingutils.SignQBFTMsg(ks.OperatorKeys[3], types.OperatorID(3), &qbft.Message{
			MsgType:              qbft.PrepareMsgType,
			Height:               qbft.FirstHeight,
			Round:                qbft.FirstRound,
			Identifier:           testingutils.TestingIdentifier,
			Root:                 testingutils.TestingQBFTRootData,
			PrepareJustification: testingutils.MarshalJustifications(valid[:1]),
		}),
	}
	return valid, nested
}

// nestedPrepareJustificationTest processes a round 2 proposal justified by prepared round changes that embed
// rcPrepares, with proposalPrepares as its prepare justifications
func nestedPrepareJustificationTest(
	name, doc string,
	ks *testingutils.TestKeySet,
	rcPrepares, proposalPrepares []*types.SignedSSVMessage,
	expectedErrorCode int,
) tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2

	rcJustifications := testingutils.MarshalJustifications(rcPrepares)
	rcMsgs := []*types.SignedSSVMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[1], types.OperatorID(1), 2, rcJustifications),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[2], types.OperatorID(2), 2, rcJustifications),
		testingutils.TestingRoundChangeMessageWithRoundAndRC(ks.OperatorKeys[3], types.OperatorID(3), 2, rcJustifications),
	}

	inputMessages := []*types.SignedSSVMessage{
		testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(proposalPrepares),
		),
	}

	return tests.NewMsgProcessingSpecTest(
		name,
		doc,
		pre,
		"",
		nil,
		inputMessages,
		nil,
		expectedErrorCode,
		nil,
		ks,
	)
}
