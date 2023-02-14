package tests

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/pkg/errors"
)

// TenOperators tests a simple full happy flow until decided
func TenOperators() *MsgProcessingSpecTest {

	pre := testingutils.TenOperatorsInstanceAlea()

	N_OPERATORS := 10
	STRONG_SUPPORT := 7
	// WEAK_SUPPORT := 4
	MAIN_OPERATOR := 1

	// client requests
	signedProposal1 := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea([]byte{1, 2, 3, 4}),
	})
	signedProposal2 := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea([]byte{5, 6, 7, 8}),
	})

	proposalData1, err := signedProposal1.Message.GetProposalData()
	if err != nil {
		errors.Wrap(err, "could not get proposal data 1 in happy flow")
	}

	proposalData2, err := signedProposal2.Message.GetProposalData()
	if err != nil {
		errors.Wrap(err, "could not get proposal data 2 in happy flow")
	}

	hash, err := alea.GetProposalsHash([]*alea.ProposalData{proposalData1, proposalData2})
	if err != nil {
		errors.Wrap(err, "could not produce hash in happy flow")
	}
	priority := alea.FirstPriority
	author := types.OperatorID(1)

	signedMessages := []*alea.SignedMessage{signedProposal1, signedProposal2}

	// msgs for VCBC agreement

	readyMsgs := make([]*alea.SignedMessage, 0)
	signedMessage := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[types.OperatorID(1)], types.OperatorID(1), &alea.Message{
		MsgType:    alea.VCBCReadyMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCReadyDataBytes(hash, priority, author),
	})
	readyMsgs = append(readyMsgs, signedMessage)
	for opID := 2; opID <= STRONG_SUPPORT; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
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

	for opID := 2; opID <= N_OPERATORS; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
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
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
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
	votes := []byte{1}

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	// finish

	for opID := 2; opID <= N_OPERATORS; opID++ {
		if opID != MAIN_OPERATOR {
			signedMessage := testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote, acRound),
			})
			signedMessages = append(signedMessages, signedMessage)
		}
	}

	return &MsgProcessingSpecTest{
		Name:          "happy flow ten operators",
		Pre:           pre,
		PostRoot:      "e6444688949542a61c0f29979933b211e9bc686f88edc3588c393113d1fc0718",
		InputMessages: signedMessages,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes([]*alea.ProposalData{proposalData1, proposalData2}, priority, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCFinalMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCFinalDataBytes(hash, priority, aggregatedMsgBytes, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(vote, round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote, round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(vote, round+1, acRound),
			}),
			testingutils.SignAleaMsg(testingutils.Testing10SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(byte(0), round, acRound+1),
			}),
		},
	}
}
