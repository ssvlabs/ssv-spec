package roundchange

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// The tests in this file run an instance whose identifier is the committee's COMMITTEE message ID. The other
// identifier is the same committee's AGGREGATOR_COMMITTEE message ID. Prepares with the other identifier are
// otherwise valid: same height, round, root and signers, with valid signatures.

// JustificationDifferentIdentifier tests a prepared round change whose prepare justifications have another
// instance's Identifier and MsgID
func JustificationDifferentIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := testingutils.CommitteeMsgID(ks), testingutils.AggregatorCommitteeMsgID(ks)

	return justificationIdentifierTest(
		"justification different identifier",
		testdoc.RoundChangeJustificationDifferentIdentifierDoc,
		ks, own, other, other,
	)
}

// JustificationDifferentMsgID tests a prepared round change whose prepare justifications have the instance
// identifier as Identifier but another instance's identifier as MsgID
func JustificationDifferentMsgID() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := testingutils.CommitteeMsgID(ks), testingutils.AggregatorCommitteeMsgID(ks)

	return justificationIdentifierTest(
		"justification different msg id",
		testdoc.RoundChangeJustificationDifferentMsgIDDoc,
		ks, own, own, other,
	)
}

// JustificationDifferentQBFTIdentifier tests a prepared round change whose prepare justifications have the instance
// identifier as MsgID but another instance's identifier as Identifier
func JustificationDifferentQBFTIdentifier() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	own, other := testingutils.CommitteeMsgID(ks), testingutils.AggregatorCommitteeMsgID(ks)

	return justificationIdentifierTest(
		"justification different qbft identifier",
		testdoc.RoundChangeJustificationDifferentQBFTIdentifierDoc,
		ks, own, other, own,
	)
}

// justificationIdentifierTest processes a round 2 round change with the instance identifier own, prepared in round 1
// and justified by a quorum of prepares with the given Identifier and MsgID
func justificationIdentifierTest(name, doc string, ks *testingutils.TestKeySet, own, identifier, msgID []byte) tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ID = own
	pre.State.Round = 2

	prepares := make([]*types.SignedSSVMessage, 0, 3)
	for _, id := range []types.OperatorID{1, 2, 3} {
		prepares = append(prepares, testingutils.SignQBFTMsgWithMsgID(ks.OperatorKeys[id], id, &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: identifier,
			Root:       testingutils.TestingQBFTRootData,
		}, msgID))
	}

	rc := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    2,
		Identifier:               own,
		Root:                     testingutils.TestingQBFTRootData,
		DataRound:                qbft.FirstRound,
		RoundChangeJustification: testingutils.MarshalJustifications(prepares),
	})
	rc.FullData = testingutils.TestingQBFTFullData

	return tests.NewMsgProcessingSpecTest(
		name,
		doc,
		pre,
		"",
		nil,
		[]*types.SignedSSVMessage{rc},
		nil,
		types.MessageIdentifierInvalidErrorCode,
		nil,
		ks,
	)
}
