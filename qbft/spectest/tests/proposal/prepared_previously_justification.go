package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreparedPreviouslyJustification tests a proposal for > 1 round, prepared previously with quorum of round change msgs justification
func PreparedPreviouslyJustification() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	ks10 := testingutils.Testing10SharesSet() // TODO: should be 4?
	pre := testingutils.BaseInstance()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
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
		testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
	}
	msgs = append(msgs, prepareMsgs...)
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	)
	return &tests.MsgProcessingSpecTest{
		Name:          "previously prepared proposal",
		Pre:           pre,
		PostRoot:      "6eeb1befcd4883bfa86293c940680546911f244fda81a10140e283586365f955",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks10.Shares[1], types.OperatorID(1)),
			testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1)),
			testingutils.TestingRoundChangeMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
				testingutils.TestingQBFTRootData, qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs)),
			testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
				testingutils.TestingQBFTRootData,
				testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs)),
			testingutils.TestingPrepareMessageWithRound(ks10.Shares[1], types.OperatorID(1), 2),
		},
	}
}

func PreparedPreviouslyJustificationDebug() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Root:       testingutils.TestingQBFTRootData,
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Root:       testingutils.TestingQBFTRootData,
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Root:       testingutils.TestingQBFTRootData,
		}),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:                  qbft.RoundChangeMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    2,
			Identifier:               []byte{1, 2, 3, 4},
			DataRound:                qbft.FirstRound,
			RoundChangeJustification: testingutils.MarshalJustifications(prepareMsgs),
			Root:                     testingutils.TestingQBFTRootData,
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:                  qbft.RoundChangeMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    2,
			Identifier:               []byte{1, 2, 3, 4},
			RoundChangeJustification: testingutils.MarshalJustifications(prepareMsgs),
			Root:                     testingutils.TestingQBFTRootData,
		}),
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:                  qbft.RoundChangeMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    2,
			Identifier:               []byte{1, 2, 3, 4},
			RoundChangeJustification: testingutils.MarshalJustifications(prepareMsgs),
			Root:                     testingutils.TestingQBFTRootData,
		}),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Root:       testingutils.TestingQBFTRootData,
		}),
	}
	msgs = append(msgs, prepareMsgs...)
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:                  qbft.ProposalMsgType,
			Height:                   qbft.FirstHeight,
			Round:                    2,
			Identifier:               []byte{1, 2, 3, 4},
			RoundChangeJustification: testingutils.MarshalJustifications(rcMsgs),
			PrepareJustification:     testingutils.MarshalJustifications(prepareMsgs),
			Root:                     testingutils.TestingQBFTRootData,
		}),
	)
	return &tests.MsgProcessingSpecTest{
		Name:          "previously prepared proposal",
		Pre:           pre,
		PostRoot:      "146c12b2ad626200f2bb8f933ec259b92e9128a3010b6a55c8de95225199c45a",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Root:       testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Root:       testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.RoundChangeMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    2,
				Identifier:               []byte{1, 2, 3, 4},
				RoundChangeJustification: testingutils.MarshalJustifications(prepareMsgs),
				Root:                     testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:                  qbft.ProposalMsgType,
				Height:                   qbft.FirstHeight,
				Round:                    2,
				Identifier:               []byte{1, 2, 3, 4},
				RoundChangeJustification: testingutils.MarshalJustifications(rcMsgs),
				PrepareJustification:     testingutils.MarshalJustifications(prepareMsgs),
				Root:                     testingutils.TestingQBFTRootData,
			}),
			testingutils.SignQBFTMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      2,
				Identifier: []byte{1, 2, 3, 4},
				Root:       testingutils.TestingQBFTRootData,
			}),
		},
	}
}
