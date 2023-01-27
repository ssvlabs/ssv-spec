package tests

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/pkg/errors"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	pre := testingutils.BaseInstanceAlea()
	// client requests
	signedProposal1 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ProposalMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytesAlea([]byte{1, 2, 3, 4}),
	})
	signedProposal2 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
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

	// msgs for VCBC agreement
	vcbcReady1 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.VCBCReadyMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCReadyDataBytes(hash, priority, author),
	})
	vcbcReady2 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
		MsgType:    alea.VCBCReadyMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCReadyDataBytes(hash, priority, author),
	})
	vcbcReady3 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
		MsgType:    alea.VCBCReadyMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.VCBCReadyDataBytes(hash, priority, author),
	})

	// msgs for VCBC agreement
	// init
	vote := byte(1)
	round := alea.FirstRound
	abainit1 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ABAInitMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAInitDataBytes(vote, round),
	})
	abainit4 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
		MsgType:    alea.ABAInitMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAInitDataBytes(vote, round),
	})
	abainit2 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAInitMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAInitDataBytes(vote, round),
	})
	abainit3 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
		MsgType:    alea.ABAInitMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAInitDataBytes(vote, round),
	})

	// aux
	abaaux1 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ABAAuxMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAAuxDataBytes(vote, round),
	})
	abaaux4 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
		MsgType:    alea.ABAAuxMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAAuxDataBytes(vote, round),
	})
	abaaux2 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAAuxMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAAuxDataBytes(vote, round),
	})
	abaaux3 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
		MsgType:    alea.ABAAuxMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAAuxDataBytes(vote, round),
	})

	// conf
	votes := []byte{1}

	abaconf1 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ABAConfMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAConfDataBytes(votes, round),
	})
	abaconf4 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
		MsgType:    alea.ABAConfMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAConfDataBytes(votes, round),
	})
	abaconf2 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAConfMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAConfDataBytes(votes, round),
	})
	abaconf3 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
		MsgType:    alea.ABAConfMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAConfDataBytes(votes, round),
	})

	// finish
	abafinish1 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(vote),
	})
	abafinish4 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[4], types.OperatorID(4), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(vote),
	})
	abafinish2 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[2], types.OperatorID(2), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(vote),
	})
	abafinish3 := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &alea.Message{
		MsgType:    alea.ABAFinishMsgType,
		Height:     alea.FirstHeight,
		Round:      alea.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ABAFinishDataBytes(vote),
	})

	aggregatedReadyMessages, err := alea.AggregateMsgs([]*alea.SignedMessage{vcbcReady1, vcbcReady2, vcbcReady3})
	if err != nil {
		errors.Wrap(err, "could not aggregate vcbcready messages in happy flow")
	}
	proof := aggregatedReadyMessages.Signature

	msgs := []*alea.SignedMessage{signedProposal1, signedProposal2, vcbcReady1, vcbcReady2, vcbcReady3, abainit1, abainit2, abainit3, abainit4, abaaux1, abaaux2, abaaux3, abaaux4, abaconf1, abaconf2, abaconf3, abaconf4, abafinish1, abafinish2, abafinish3, abafinish4}
	return &MsgProcessingSpecTest{
		Name:          "happy flow",
		Pre:           pre,
		PostRoot:      "65c738989062ab5b80a80e9a207e107772ca56e7fe9897d6981e98b7b60737e6",
		InputMessages: msgs,
		OutputMessages: []*alea.SignedMessage{
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCSendMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCSendDataBytes([]*alea.ProposalData{proposalData1, proposalData2}, priority, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.VCBCFinalMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.VCBCFinalDataBytes(hash, priority, proof, author),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(vote, round),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAAuxMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAAuxDataBytes(vote, round),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAConfMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAConfDataBytes(votes, round),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAFinishMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAFinishDataBytes(vote),
			}),
			testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
				MsgType:    alea.ABAInitMsgType,
				Height:     alea.FirstHeight,
				Round:      alea.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ABAInitDataBytes(vote, round+1),
			}),
		},
	}
}
