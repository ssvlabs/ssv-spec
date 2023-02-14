package flow

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/pkg/errors"
)

// Flow7Op tests a 7 operators flow
func Flow7Op() *tests.MsgProcessingSpecTest {

	pre := testingutils.SevenOperatorsInstanceAlea()

	N_OPERATORS := 7
	STRONG_SUPPORT := 5
	// WEAK_SUPPORT := 3
	MAIN_OPERATOR := 1

	// client requests
	signedMessages := make([]*alea.SignedMessage, 0)
	signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(tests.ProposalData1.Data),
	})
	signedMessages = append(signedMessages, signedMessage)
	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(tests.ProposalData2.Data),
	})
	signedMessages = append(signedMessages, signedMessage)
	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(tests.ProposalData3.Data),
	})
	signedMessages = append(signedMessages, signedMessage)
	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea(tests.ProposalData4.Data),
	})
	signedMessages = append(signedMessages, signedMessage)

	// msgs for VCBC agreement

	hash := tests.Hash
	priority := alea.FirstPriority
	author := types.OperatorID(1)

	readyMsgs := make([]*alea.SignedMessage, 0)
	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.VCBCReadyMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCReadyDataBytes(hash, priority, author),
	})
	readyMsgs = append(readyMsgs, signedMessage)
	for opID := 2; opID <= STRONG_SUPPORT; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(hash, priority, author),
		})
		signedMessages = append(signedMessages, signedMessage)
		readyMsgs = append(readyMsgs, signedMessage)
	}

	// get ready message and make proof
	aggregatedReadyMessages, err := alea.AggregateMsgs(readyMsgs)
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	aggregatedMsgBytes, err := aggregatedReadyMessages.Encode()
	if err != nil {
		errors.Wrap(err, "could not encode aggregated msg")
	}

	// msgs for VCBC agreement
	// init
	vote := byte(1)
	round := alea.FirstRound
	acRound := alea.FirstACRound

	for opID := 2; opID <= 2; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(vote, round, acRound),
		})
		signedMessages = append(signedMessages, signedMessage)
	}
	vote = byte(0)
	for opID := 3; opID <= N_OPERATORS; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(vote, round, acRound),
		})
		signedMessages = append(signedMessages, signedMessage)
	}

	// aux

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote, round, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// conf
	votes := []byte{0}

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// init
	vote = byte(1)
	round = alea.FirstRound
	acRound = alea.FirstACRound

	for opID := 2; opID <= 2; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(vote, round+1, acRound),
		})
		signedMessages = append(signedMessages, signedMessage)
	}
	vote = byte(0)
	for opID := 3; opID <= N_OPERATORS; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(vote, round+1, acRound),
		})
		signedMessages = append(signedMessages, signedMessage)
	}

	// aux

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote, round+1, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// conf
	votes = []byte{0}

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round+1, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// finish

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// msgs for VCBC agreement

	hash2 := tests.Hash2
	priority2 := alea.FirstPriority + 1
	author = types.OperatorID(1)

	readyMsgs = make([]*alea.SignedMessage, 0)
	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.VCBCReadyMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCReadyDataBytes(hash2, priority2, author),
	})
	readyMsgs = append(readyMsgs, signedMessage)
	for opID := 2; opID <= STRONG_SUPPORT; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(hash2, priority2, author),
		})
		signedMessages = append(signedMessages, signedMessage)
		readyMsgs = append(readyMsgs, signedMessage)
	}

	// get ready message and make proof
	aggregatedReadyMessages2, err := alea.AggregateMsgs(readyMsgs)
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	aggregatedMsgBytes2, err := aggregatedReadyMessages2.Encode()
	if err != nil {
		errors.Wrap(err, "could not encode aggregated msg")
	}

	proposals3 := tests.ProposalDataList3
	priority3 := alea.FirstPriority
	author2 := types.OperatorID(2)

	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(2)], types.OperatorID(2), &alea.Message{
		MsgType:    alea.VCBCSendMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCSendDataBytes(proposals3, priority3, author2),
	})
	signedMessages = append(signedMessages, signedMessage)

	hash3 := tests.Hash3

	readyMsgs = make([]*alea.SignedMessage, 0)
	for opID := 1; opID <= STRONG_SUPPORT; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(hash3, priority3, author2),
		})
		signedMessages = append(signedMessages, signedMessage)
		readyMsgs = append(readyMsgs, signedMessage)
	}

	// get ready message and make proof
	aggregatedReadyMessages3, err := alea.AggregateMsgs(readyMsgs)
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	aggregatedMsgBytes3, err := aggregatedReadyMessages3.Encode()
	if err != nil {
		errors.Wrap(err, "could not encode aggregated msg")
	}

	signedMessage = testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(2)], types.OperatorID(2), &alea.Message{
		MsgType:    alea.VCBCFinalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCFinalDataBytes(hash3, priority3, aggregatedMsgBytes3, author2),
	})
	signedMessages = append(signedMessages, signedMessage)

	// msgs for VCBC agreement
	// init
	vote2 := byte(1)
	round2 := alea.FirstRound
	acRound2 := alea.FirstACRound + 1

	for opID := 2; opID <= N_OPERATORS; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.ABAInitMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.ABAInitDataBytes(vote2, round2, acRound2),
		})
		signedMessages = append(signedMessages, signedMessage)
	}

	// aux

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote2, round2, acRound2),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// conf
	votes2 := []byte{1}

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes2, round2, acRound2),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// finish

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote2, acRound2),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "flow 7 operators 4 proposals",
		Pre:           pre,
		PostRoot:      "411f1b855f04715d9caf8f57498c01ce239e5241ed403e83aa215316ced886f5",
		InputMessages: signedMessages,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList, priority, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes(tests.ProposalDataList2, priority+1, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCFinalMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCFinalDataBytes(hash, priority, aggregatedMsgBytes, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote, round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(vote, round+1, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote, round+1, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round+1, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(vote, round+2, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), round, acRound+1),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCFinalMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCFinalDataBytes(hash2, priority2, aggregatedMsgBytes2, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCReadyMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCReadyDataBytes(hash3, priority3, author2),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(1), round2, acRound2),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote2, round2, acRound2),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes2, round2, acRound2),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote2, acRound2),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(1), round2+1, acRound2),
			}),
			testingutils.SignAleaMsg(testingutils.Testing7SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), round, acRound2+1),
			}),
		},
	}
}
