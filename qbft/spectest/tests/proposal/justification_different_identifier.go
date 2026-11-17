package proposal

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// The tests in this file run an instance whose identifier is the committee's COMMITTEE message ID. The other
// identifier is the same committee's AGGREGATOR_COMMITTEE message ID. Messages with the other identifier are
// otherwise valid: same height, round, root and signers, with valid signatures.

// RCJustificationDifferentIdentifier tests a proposal for round > 1 justified by a quorum of round changes whose
// Identifier and MsgID are another instance's identifier
func RCJustificationDifferentIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := identifiers(ks)

	rcMsgs := []*types.SignedSSVMessage{
		roundChangeWithIdentifiers(ks, 1, other, other, nil),
		roundChangeWithIdentifiers(ks, 2, other, other, nil),
		roundChangeWithIdentifiers(ks, 3, other, other, nil),
	}

	return differentIdentifierTest(
		"proposal rc justification different identifier",
		testdoc.ProposalRCJustificationDifferentIdentifierDoc,
		ks, own, rcMsgs, nil,
		types.MessageIdentifierInvalidErrorCode,
	)
}

// RCJustificationDifferentMsgID tests a proposal for round > 1 justified by a quorum of round changes whose
// Identifier is the instance identifier but whose MsgID is another instance's identifier
func RCJustificationDifferentMsgID() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := identifiers(ks)

	rcMsgs := []*types.SignedSSVMessage{
		roundChangeWithIdentifiers(ks, 1, own, other, nil),
		roundChangeWithIdentifiers(ks, 2, own, other, nil),
		roundChangeWithIdentifiers(ks, 3, own, other, nil),
	}

	return differentIdentifierTest(
		"proposal rc justification different msg id",
		testdoc.ProposalRCJustificationDifferentMsgIDDoc,
		ks, own, rcMsgs, nil,
		types.MessageIdentifierInvalidErrorCode,
	)
}

// RCJustificationDifferentQBFTIdentifier tests a proposal for round > 1 justified by a quorum of round changes whose
// MsgID is the instance identifier but whose Identifier is another instance's identifier
func RCJustificationDifferentQBFTIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := identifiers(ks)

	rcMsgs := []*types.SignedSSVMessage{
		roundChangeWithIdentifiers(ks, 1, other, own, nil),
		roundChangeWithIdentifiers(ks, 2, other, own, nil),
		roundChangeWithIdentifiers(ks, 3, other, own, nil),
	}

	return differentIdentifierTest(
		"proposal rc justification different qbft identifier",
		testdoc.ProposalRCJustificationDifferentQBFTIdentifierDoc,
		ks, own, rcMsgs, nil,
		types.MessageIdentifierInvalidErrorCode,
	)
}

// PrepareJustificationDifferentIdentifier tests a proposal for round > 1, prepared previously, whose prepare
// justifications have another instance's identifier. The proposal justification check reports any invalid prepare
// justification as PrepareMessageInvalidErrorCode.
func PrepareJustificationDifferentIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := identifiers(ks)

	rcMsgs := preparedRoundChanges(ks, own, preparesWithIdentifiers(ks, own, own))

	return differentIdentifierTest(
		"proposal prepare justification different identifier",
		testdoc.ProposalPrepareJustificationDifferentIdentifierDoc,
		ks, own, rcMsgs, preparesWithIdentifiers(ks, other, other),
		types.PrepareMessageInvalidErrorCode,
	)
}

// RCJustificationNestedPrepareDifferentIdentifier tests a proposal for round > 1, prepared previously, whose
// round change justifications embed prepares with another instance's identifier
func RCJustificationNestedPrepareDifferentIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := identifiers(ks)

	rcMsgs := preparedRoundChanges(ks, own, preparesWithIdentifiers(ks, other, other))

	return differentIdentifierTest(
		"proposal rc justification nested prepare different identifier",
		testdoc.ProposalRCJustificationNestedPrepareDifferentIdentifierDoc,
		ks, own, rcMsgs, preparesWithIdentifiers(ks, own, own),
		types.MessageIdentifierInvalidErrorCode,
	)
}

// identifiers returns the instance identifier and another identifier of the same committee, for a different role
func identifiers(ks *testingutils.TestKeySet) (own, other []byte) {
	return testingutils.CommitteeMsgID(ks), testingutils.AggregatorCommitteeMsgID(ks)
}

// preparesWithIdentifiers returns a quorum of round 1 prepares with the given Identifier and MsgID
func preparesWithIdentifiers(ks *testingutils.TestKeySet, identifier, msgID []byte) []*types.SignedSSVMessage {
	ret := make([]*types.SignedSSVMessage, 0, 3)
	for _, id := range []types.OperatorID{1, 2, 3} {
		ret = append(ret, testingutils.SignQBFTMsgWithMsgID(ks.OperatorKeys[id], id, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: identifier,
			Root:       testingutils.TestingQBFTRootData,
		}, msgID))
	}
	return ret
}

// preparedRoundChanges returns a quorum of round 2 round changes with the instance identifier, prepared in round 1
// and justified by prepares
func preparedRoundChanges(ks *testingutils.TestKeySet, own []byte, prepares []*types.SignedSSVMessage) []*types.SignedSSVMessage {
	return []*types.SignedSSVMessage{
		roundChangeWithIdentifiers(ks, 1, own, own, prepares),
		roundChangeWithIdentifiers(ks, 2, own, own, prepares),
		roundChangeWithIdentifiers(ks, 3, own, own, prepares),
	}
}

// roundChangeWithIdentifiers returns a round 2 round change with the given Identifier and MsgID. It is prepared in
// round 1 and justified by prepares if prepares isn't nil, and not prepared otherwise.
func roundChangeWithIdentifiers(
	ks *testingutils.TestKeySet,
	id types.OperatorID,
	identifier, msgID []byte,
	prepares []*types.SignedSSVMessage,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     qbft.FirstHeight,
		Round:      2,
		Identifier: identifier,
		Root:       testingutils.TestingQBFTRootData,
		DataRound:  qbft.NoRound,
	}
	if prepares != nil {
		msg.DataRound = qbft.FirstRound
		msg.RoundChangeJustification = testingutils.MarshalJustifications(prepares)
	}
	ret := testingutils.SignQBFTMsgWithMsgID(ks.OperatorKeys[id], id, msg, msgID)
	if prepares != nil {
		ret.FullData = testingutils.TestingQBFTFullData
	}
	return ret
}

// differentIdentifierTest processes a round 2 proposal with the instance identifier, justified by rcMsgs and
// prepares
func differentIdentifierTest(
	name, doc string,
	ks *testingutils.TestKeySet,
	own []byte,
	rcMsgs, prepares []*types.SignedSSVMessage,
	expectedErrorCode int,
) tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ID = own
	pre.State.Round = 2

	msg := &qbft.Message{
		MsgType:                  qbft.ProposalMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    2,
		Identifier:               own,
		Root:                     testingutils.TestingQBFTRootData,
		RoundChangeJustification: testingutils.MarshalJustifications(rcMsgs),
	}
	if prepares != nil {
		msg.PrepareJustification = testingutils.MarshalJustifications(prepares)
	}
	proposal := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), msg)
	proposal.FullData = testingutils.TestingQBFTFullData

	return tests.NewMsgProcessingSpecTest(
		name,
		doc,
		pre,
		"",
		nil,
		[]*types.SignedSSVMessage{proposal},
		nil,
		expectedErrorCode,
		nil,
		ks,
	)
}
