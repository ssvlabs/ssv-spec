package tests

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/pkg/errors"
)

var SignedProposal1 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{1, 2, 3, 4}),
})
var SignedProposal2 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{5, 6, 7, 8}),
})
var SignedProposal3 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{9, 10, 11, 12}),
})
var SignedProposal4 = testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &alea.Message{
	MsgType:    alea.ProposalMsgType,
	Height:     alea.FirstHeight,
	Round:      alea.FirstRound,
	Identifier: []byte{1, 2, 3, 4},
	Data:       testingutils.ProposalDataBytesAlea([]byte{13, 14, 15, 16}),
})

var ProposalData1, _ = SignedProposal1.Message.GetProposalData()
var ProposalData2, _ = SignedProposal2.Message.GetProposalData()
var ProposalData3, _ = SignedProposal3.Message.GetProposalData()
var ProposalData4, _ = SignedProposal4.Message.GetProposalData()

var ProposalDataList = []*alea.ProposalData{ProposalData1, ProposalData2}
var ProposalDataList2 = []*alea.ProposalData{ProposalData1, ProposalData2}

var Entries = [][]*alea.ProposalData{{ProposalData1, ProposalData2}, {ProposalData3, ProposalData4}}
var Priorities = []alea.Priority{alea.FirstPriority, alea.FirstPriority + 1}

var Hash, _ = alea.GetProposalsHash([]*alea.ProposalData{ProposalData1, ProposalData2})
var Hash2, _ = alea.GetProposalsHash([]*alea.ProposalData{ProposalData3, ProposalData4})

var AggregatedMsgBytes = func() []byte {

	readyMsgs := make([]*alea.SignedMessage, 0)
	for opID := 1; opID <= 4; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(Hash, alea.FirstPriority, types.OperatorID(1)),
		})
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
	return aggregatedMsgBytes
}()

var AggregatedMsgBytes2 = func() []byte {

	readyMsgs := make([]*alea.SignedMessage, 0)
	for opID := 1; opID <= 4; opID++ {
		signedMessage := testingutils.SignAleaMsg(testingutils.Testing4SharesSet().Shares[types.OperatorID(opID)], types.OperatorID(opID), &alea.Message{
			MsgType:    alea.VCBCReadyMsgType,
			Height:     alea.FirstHeight,
			Round:      alea.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.VCBCReadyDataBytes(Hash2, alea.FirstPriority, types.OperatorID(1)),
		})
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
	return aggregatedMsgBytes
}()

var AggregatedMsgBytesList = [][]byte{AggregatedMsgBytes, AggregatedMsgBytes2}
