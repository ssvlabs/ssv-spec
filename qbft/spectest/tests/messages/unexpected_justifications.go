package messages

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// RoundChangeNotPreparedWithJustifications tests a not prepared round change that carries round change justifications
func RoundChangeNotPreparedWithJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	justifications := testingutils.MarshalJustifications([]*types.SignedSSVMessage{
		testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.OperatorKeys[3], types.OperatorID(3)),
	})
	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   qbft.FirstHeight,
		Round:                    2,
		Identifier:               testingutils.TestingIdentifier,
		Root:                     testingutils.TestingQBFTRootData,
		DataRound:                qbft.NoRound,
		RoundChangeJustification: justifications,
	})

	return tests.NewMsgSpecTest(
		"rc not prepared with justifications",
		testdoc.MessagesRCNotPreparedWithJustificationsDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.UnexpectedJustificationErrorCode,
		ks,
	)
}

// PrepareWithPrepareJustifications tests a prepare message that carries prepare justifications
func PrepareWithPrepareJustifications() tests.SpecTest {
	return withPrepareJustifications(qbft.PrepareMsgType, "prepare with prepare justifications")
}

// CommitWithPrepareJustifications tests a commit message that carries prepare justifications
func CommitWithPrepareJustifications() tests.SpecTest {
	return withPrepareJustifications(qbft.CommitMsgType, "commit with prepare justifications")
}

// RoundChangeWithPrepareJustifications tests a round change message that carries prepare justifications
func RoundChangeWithPrepareJustifications() tests.SpecTest {
	return withPrepareJustifications(qbft.RoundChangeMsgType, "rc with prepare justifications")
}

func withPrepareJustifications(msgType qbft.MessageType, name string) tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    msgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: testingutils.TestingIdentifier,
		Root:       testingutils.TestingQBFTRootData,
		PrepareJustification: testingutils.MarshalJustifications([]*types.SignedSSVMessage{
			testingutils.TestingPrepareMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		}),
	})

	return tests.NewMsgSpecTest(
		name,
		testdoc.MessagesNonProposalWithPrepareJustificationsDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.UnexpectedJustificationErrorCode,
		ks,
	)
}

// CommitWithRoundChangeJustifications tests a commit message that carries round change justifications
func CommitWithRoundChangeJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.SignQBFTMsg(ks.OperatorKeys[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: testingutils.TestingIdentifier,
		Root:       testingutils.TestingQBFTRootData,
		RoundChangeJustification: testingutils.MarshalJustifications([]*types.SignedSSVMessage{
			testingutils.TestingRoundChangeMessage(ks.OperatorKeys[1], types.OperatorID(1)),
		}),
	})

	return tests.NewMsgSpecTest(
		"commit with rc justifications",
		testdoc.MessagesCommitWithRoundChangeJustificationsDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.UnexpectedJustificationErrorCode,
		ks,
	)
}

// RoundChangeJustificationWithFullData tests a prepared round change whose prepare justifications carry full data
func RoundChangeJustificationWithFullData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	justifications := make([][]byte, 0)
	for i := 1; i <= 3; i++ {
		prepare := testingutils.TestingPrepareMessage(ks.OperatorKeys[types.OperatorID(i)], types.OperatorID(i))
		prepare.FullData = testingutils.TestingQBFTFullData
		byts, err := prepare.MarshalSSZ()
		if err != nil {
			panic(err.Error())
		}
		justifications = append(justifications, byts)
	}
	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.OperatorKeys[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData, qbft.FirstRound,
		justifications)

	return tests.NewMsgSpecTest(
		"rc justification with full data",
		testdoc.MessagesRCJustificationWithFullDataDoc,
		[]*types.SignedSSVMessage{msg},
		nil,
		nil,
		types.UnexpectedJustificationErrorCode,
		ks,
	)
}
