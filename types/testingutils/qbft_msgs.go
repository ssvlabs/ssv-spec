package testingutils

import (
	"crypto/rsa"
	"crypto/sha256"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingIdentifier = []byte{1, 2, 3, 4}

var DifferentFullData = append(TestingQBFTFullData, []byte("different")...)
var DifferentRoot = func() [32]byte {
	return sha256.Sum256(DifferentFullData)
}()

var MarshalJustifications = func(msgs []*types.SignedSSVMessage) [][]byte {
	bytes, err := qbft.MarshalJustifications(msgs)
	if err != nil {
		panic(err)
	}

	return bytes
}

var MultiSignQBFTMsg = func(sks []*rsa.PrivateKey, ids []types.OperatorID, msg *qbft.Message) *types.SignedSSVMessage {
	if len(sks) == 0 || len(ids) != len(sks) {
		panic("sks != ids")
	}
	var signed *types.SignedSSVMessage
	for i, sk := range sks {
		if signed == nil {
			signed = SignQBFTMsg(sk, ids[i], msg)
		} else {
			if err := signed.Aggregate(SignQBFTMsg(sk, ids[i], msg)); err != nil {
				panic(err.Error())
			}
		}
	}

	return signed
}

var SignQBFTMsg = func(sk *rsa.PrivateKey, id types.OperatorID, msg *qbft.Message) *types.SignedSSVMessage {
	encodedMsg, err := msg.Encode()
	if err != nil {
		panic(err)
	}

	msgID := [56]byte{}
	copy(msgID[:], msg.Identifier)

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	signature, err := SignSSVMessage(sk, ssvMsg)
	if err != nil {
		panic(err)
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{id},
		Signatures:  [][]byte{signature},
		SSVMessage:  ssvMsg,
	}
}

var TestingInvalidMessage = func(sk *rsa.PrivateKey, id types.OperatorID, msgType qbft.MessageType) *types.SignedSSVMessage {
	return TestingMultiSignerInvalidMessage([]*rsa.PrivateKey{sk}, []types.OperatorID{id}, msgType)
}
var TestingMultiSignerInvalidMessage = func(sks []*rsa.PrivateKey, ids []types.OperatorID, msgType qbft.MessageType) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    msgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{}, // invalid
		Root:       TestingQBFTRootData,
	}
	ret := MultiSignQBFTMsg(sks, ids, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}

var ToProcessingMessage = func(msg *types.SignedSSVMessage) *qbft.ProcessingMessage {
	pm, _ := qbft.NewProcessingMessage(msg)
	return pm
}

var ToProcessingMessages = func(msgs []*types.SignedSSVMessage) []*qbft.ProcessingMessage {
	ret := make([]*qbft.ProcessingMessage, 0)
	for _, msg := range msgs {
		ret = append(ret, ToProcessingMessage(msg))
	}
	return ret
}

/*
*
Proposal messages
*/
var TestingProposalMessage = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingProposalMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingProposalMessageWithID = func(sk *rsa.PrivateKey, id types.OperatorID, msgID types.MessageID) *types.SignedSSVMessage {
	ret := TestingProposalMessageWithRound(sk, id, qbft.FirstRound)

	qbftMsg, err := qbft.DecodeMessage(ret.SSVMessage.Data)
	if err != nil {
		panic(err)
	}
	qbftMsg.Identifier = msgID[:]

	return SignQBFTMsg(sk, id, qbftMsg)
}

var TestingProposalMessageWithRound = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingProposalMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, nil, nil)
}
var TestingProposalMessageWithHeight = func(sk *rsa.PrivateKey, id types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingProposalMessageWithParams(sk, id, qbft.FirstRound, height, TestingQBFTRootData, nil, nil)
}
var TestingProposalMessageDifferentRoot = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: TestingIdentifier,
		Root:       DifferentRoot,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = DifferentFullData
	return ret
}
var TestingProposalMessageWithRoundAndRC = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round, roundChangeJustification [][]byte) *types.SignedSSVMessage {
	return TestingProposalMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, roundChangeJustification, nil)
}

func TestingProposalMessageWithIdentifierAndFullData(sk *rsa.PrivateKey, id types.OperatorID, identifier, fullData []byte, height qbft.Height) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     height,
		Round:      qbft.FirstRound,
		Identifier: identifier,
		Root:       sha256.Sum256(fullData),

		RoundChangeJustification: [][]byte{},
		PrepareJustification:     [][]byte{},
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}

var TestingProposalMessageWithParams = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	root [32]byte,
	roundChangeJustification [][]byte,
	prepareJustification [][]byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:                  qbft.ProposalMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               TestingIdentifier,
		Root:                     root,
		RoundChangeJustification: roundChangeJustification,
		PrepareJustification:     prepareJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
var TestingMultiSignerProposalMessage = func(sks []*rsa.PrivateKey, ids []types.OperatorID) *types.SignedSSVMessage {
	return TestingMultiSignerProposalMessageWithParams(sks, ids, qbft.FirstRound, qbft.FirstHeight, TestingIdentifier, TestingQBFTFullData, TestingQBFTRootData)
}
var TestingMultiSignerProposalMessageWithHeight = func(sks []*rsa.PrivateKey, ids []types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingMultiSignerProposalMessageWithParams(sks, ids, qbft.FirstRound, height, TestingIdentifier, TestingQBFTFullData, TestingQBFTRootData)
}
var TestingMultiSignerProposalMessageWithParams = func(
	sk []*rsa.PrivateKey,
	id []types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier,
	fullData []byte,
	root [32]byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     height,
		Round:      round,
		Identifier: identifier,
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}

/*
*
Prepare messages
*/
var TestingPrepareMessage = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingPrepareMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingPrepareMessageWithRound = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingPrepareMessageWithParams(sk, id, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMessageWithHeight = func(sk *rsa.PrivateKey, id types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingPrepareMessageWithParams(sk, id, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMessageWrongRoot = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingPrepareMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingIdentifier, DifferentRoot)
}
var TestingPrepareMessageWithFullData = func(sk *rsa.PrivateKey, id types.OperatorID, fullData []byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: TestingIdentifier,
		Root:       sha256.Sum256(fullData),
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}
var TestingPrepareMessageWithIdentifierAndRoot = func(sk *rsa.PrivateKey, id types.OperatorID, identifier []byte, root [32]byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier,
		Root:       root,

		RoundChangeJustification: [][]byte{},
		PrepareJustification:     [][]byte{},
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = []byte{}
	return ret
}
var TestingPrepareMessageWithRoundAndFullData = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	fullData []byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     qbft.FirstHeight,
		Round:      round,
		Identifier: TestingIdentifier,
		Root:       sha256.Sum256(fullData),
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}
var TestingPrepareMessageWithParams = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier []byte,
	root [32]byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     height,
		Round:      round,
		Identifier: identifier,
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
var TestingPrepareMultiSignerMessage = func(sks []*rsa.PrivateKey, ids []types.OperatorID) *types.SignedSSVMessage {
	return TestingPrepareMultiSignerMessageWithRound(sks, ids, qbft.FirstRound)
}
var TestingPrepareMultiSignerMessageWithRound = func(sks []*rsa.PrivateKey, ids []types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithHeight = func(sks []*rsa.PrivateKey, ids []types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithHeightAndIdentifier = func(sks []*rsa.PrivateKey, ids []types.OperatorID, height qbft.Height, identifier []byte) *types.SignedSSVMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, identifier, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithParams = func(sks []*rsa.PrivateKey, ids []types.OperatorID, round qbft.Round, height qbft.Height, identifier []byte, root [32]byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     height,
		Round:      round,
		Identifier: identifier,
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sks, ids, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}

/*
*
Commit messages
*/
var TestingCommitMessage = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingCommitMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingCommitMessageWithRound = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingCommitMessageWithParams(sk, id, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData)
}
var TestingCommitMessageWithHeight = func(sk *rsa.PrivateKey, id types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData)
}
var TestingCommitMessageWrongRoot = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingIdentifier, DifferentRoot)
}
var TestingCommitMessageWrongHeight = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, 10, TestingIdentifier, DifferentRoot)
}

func TestingCommitMessageWithHeightIdentifierAndFullData(sk *rsa.PrivateKey, id types.OperatorID, height qbft.Height, identifier, fullData []byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      qbft.FirstRound,
		Identifier: identifier,
		Root:       sha256.Sum256(fullData),
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}

var TestingCommitMessageWithIdentifierAndRoot = func(sk *rsa.PrivateKey, id types.OperatorID, identifier []byte, root [32]byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier,
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = []byte{}
	return ret
}
var TestingCommitMessageWithParams = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier []byte,
	root [32]byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      round,
		Identifier: identifier,
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
var TestingCommitMultiSignerMessage = func(sks []*rsa.PrivateKey, ids []types.OperatorID) *types.SignedSSVMessage {
	return TestingCommitMultiSignerMessageWithRound(sks, ids, qbft.FirstRound)
}
var TestingCommitMultiSignerMessageWithRound = func(sks []*rsa.PrivateKey, ids []types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData, TestingQBFTFullData)
}
var TestingCommitMultiSignerMessageWithHeight = func(sks []*rsa.PrivateKey, ids []types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData, TestingQBFTFullData)
}
var TestingCommitMultiSignerMessageWithHeightAndIdentifier = func(sks []*rsa.PrivateKey, ids []types.OperatorID, height qbft.Height, identifier []byte) *types.SignedSSVMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, identifier, TestingQBFTRootData, TestingQBFTFullData)
}

func TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
	sks []*rsa.PrivateKey,
	ids []types.OperatorID,
	height qbft.Height,
	identifier, fullData []byte,
) *types.SignedSSVMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, identifier, sha256.Sum256(fullData), fullData)
}

var TestingCommitMultiSignerMessageWithParams = func(
	sks []*rsa.PrivateKey,
	ids []types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier []byte,
	root [32]byte,
	fullData []byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      round,
		Identifier: identifier,
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sks, ids, msg)
	ret.FullData = fullData
	return ret
}

/*
*
Round Change messages
*/
var TestingRoundChangeMessage = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingRoundChangeMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingRoundChangeMessageWithRound = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingRoundChangeMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, qbft.NoRound, nil)
}
var TestingRoundChangeMessageWithHeight = func(sk *rsa.PrivateKey, id types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingRoundChangeMessageWithParams(sk, id, qbft.FirstRound, height, TestingQBFTRootData, qbft.NoRound, nil)
}
var TestingRoundChangeMessageWithRoundAndHeight = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round, height qbft.Height) *types.SignedSSVMessage {
	return TestingRoundChangeMessageWithParams(sk, id, round, height, TestingQBFTRootData, qbft.NoRound, nil)
}
var TestingRoundChangeMessageWithRoundAndRC = func(sk *rsa.PrivateKey, id types.OperatorID, round qbft.Round, roundChangeJustification [][]byte) *types.SignedSSVMessage {
	ret := TestingRoundChangeMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, qbft.FirstRound, roundChangeJustification)
	ret.FullData = TestingQBFTFullData
	return ret
}
var TestingRoundChangeMessageWithHeightAndIdentifier = func(sk *rsa.PrivateKey, id types.OperatorID, height qbft.Height, identifier []byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     height,
		Round:      qbft.FirstRound,
		Identifier: TestingIdentifier,
		Root:       TestingQBFTRootData,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
var TestingRoundChangeMessageWithRoundAndFullData = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	fullData []byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     qbft.FirstHeight,
		Round:      round,
		Identifier: TestingIdentifier,
		Root:       sha256.Sum256(fullData),
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}
var TestingRoundChangeMessageWithParams = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	root [32]byte,
	dataRound qbft.Round,
	roundChangeJustification [][]byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               TestingIdentifier,
		Root:                     root,
		DataRound:                dataRound,
		RoundChangeJustification: roundChangeJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}

var TestingRoundChangeMessageWithParamsAndFullData = func(
	sk *rsa.PrivateKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	root [32]byte,
	dataRound qbft.Round,
	roundChangeJustification [][]byte,
	fullData []byte,
) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               TestingIdentifier,
		Root:                     root,
		DataRound:                dataRound,
		RoundChangeJustification: roundChangeJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}

var TestingMultiSignerRoundChangeMessage = func(sks []*rsa.PrivateKey, ids []types.OperatorID) *types.SignedSSVMessage {
	return TestingMultiSignerRoundChangeMessageWithParams(sks, ids, qbft.FirstRound, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingMultiSignerRoundChangeMessageWithRound = func(sks []*rsa.PrivateKey, ids []types.OperatorID, round qbft.Round) *types.SignedSSVMessage {
	return TestingMultiSignerRoundChangeMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingMultiSignerRoundChangeMessageWithHeight = func(sks []*rsa.PrivateKey, ids []types.OperatorID, height qbft.Height) *types.SignedSSVMessage {
	return TestingMultiSignerRoundChangeMessageWithParams(sks, ids, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingRoundChangeMessageWrongRoot = func(sk *rsa.PrivateKey, id types.OperatorID) *types.SignedSSVMessage {
	return TestingRoundChangeMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, DifferentRoot, qbft.NoRound, nil)
}
var TestingMultiSignerRoundChangeMessageWithParams = func(sk []*rsa.PrivateKey, id []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *types.SignedSSVMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     height,
		Round:      round,
		Identifier: TestingIdentifier,
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
