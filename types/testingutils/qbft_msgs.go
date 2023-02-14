package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var DefaultIdentifier = []byte{1, 2, 3, 4}
var WrongRoot = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 5, 6, 7, 8, 9}

var MarshalJustifications = func(msgs []*qbft.SignedMessage) [][]byte {
	bytes, err := qbft.MarshalJustifications(msgs)
	if err != nil {
		panic(err)
	}

	return bytes
}

var MultiSignQBFTMsg = func(sks []*bls.SecretKey, ids []types.OperatorID, msg *qbft.Message) *qbft.SignedMessage {
	if len(sks) == 0 || len(ids) != len(sks) {
		panic("sks != ids")
	}
	var signed *qbft.SignedMessage
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

var SignQBFTMsg = func(sk *bls.SecretKey, id types.OperatorID, msg *qbft.Message) *qbft.SignedMessage {
	domain := TestingSSVDomainType
	sigType := types.QBFTSignatureType

	r, _ := types.ComputeSigningRoot(msg, types.ComputeSignatureDomain(domain, sigType))
	sig := sk.SignByte(r)

	return &qbft.SignedMessage{
		Message:   *msg,
		Signers:   []types.OperatorID{id},
		Signature: sig.Serialize(),
	}
}

var TestingInvalidMessage = func(sk *bls.SecretKey, id types.OperatorID, msgType qbft.MessageType) *qbft.SignedMessage {
	return TestingMultiSignerInvalidMessage([]*bls.SecretKey{sk}, []types.OperatorID{id}, msgType)
}
var TestingMultiSignerInvalidMessage = func(sks []*bls.SecretKey, ids []types.OperatorID, msgType qbft.MessageType) *qbft.SignedMessage {
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

/*
*
Proposal messages
*/
var TestingProposalMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingProposalMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingProposalMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingProposalMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingQBFTRootData, nil, nil)
}
var TestingProposalMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingProposalMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingQBFTRootData, nil, nil)
}
var TestingProposalMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingProposalMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, WrongRoot, nil, nil)
}
var TestingProposalMessageWithRoundAndRC = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, roundChangeJustification [][]byte) *qbft.SignedMessage {
	return TestingProposalMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingQBFTRootData, roundChangeJustification, nil)
}
var TestingProposalMessageWithParams = func(sk *bls.SecretKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	root [32]byte,
	roundChangeJustification [][]byte,
	prepareJustification [][]byte,
) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:                  qbft.ProposalMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               []byte{1, 2, 3, 4},
		Root:                     root,
		RoundChangeJustification: roundChangeJustification,
		PrepareJustification:     prepareJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
var TestingMultiSignerProposalMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingMultiSignerProposalMessageWithParams(sks, ids, qbft.FirstRound, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingMultiSignerProposalMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingMultiSignerProposalMessageWithParams(sks, ids, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingMultiSignerProposalMessageWithParams = func(sk []*bls.SecretKey, id []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}

/*
*
Prepare messages
*/
var TestingPrepareMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingPrepareMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingPrepareMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingPrepareMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingPrepareMessageWithParams(sk, id, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingPrepareMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, WrongRoot)
}
var TestingPrepareMessageWithParams = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
var TestingPrepareMultiSignerMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithRound(sks, ids, qbft.FirstRound)
}
var TestingPrepareMultiSignerMessageWithRound = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithParams = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sks, ids, msg)
	return ret
}

/*
*
Commit messages
*/
var TestingCommitMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingCommitMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingCommitMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingCommitMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, WrongRoot)
}
var TestingCommitMessageWrongHeight = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, 10, WrongRoot)
}
var TestingCommitMessageWithParams = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
var TestingCommitMultiSignerMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithRound(sks, ids, qbft.FirstRound)
}
var TestingCommitMultiSignerMessageWithRound = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, DefaultIdentifier, TestingQBFTRootData)
}
var TestingCommitMultiSignerMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, DefaultIdentifier, TestingQBFTRootData)
}
var TestingCommitMultiSignerMessageWithHeightAndIdentifier = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height, identifier []byte) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, identifier, TestingQBFTRootData)
}
var TestingCommitMultiSignerMessageWithParams = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round, height qbft.Height, identifier []byte, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sks, ids, msg)
	return ret
}
var TestingInvalidCommitMultiSignerMessageWithParams = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sks, ids, msg)
	return ret
}

/*
*
Round Change messages
*/
var TestingRoundChangeMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingRoundChangeMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, qbft.NoRound, nil)
}
var TestingRoundChangeMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithParams(sk, id, qbft.FirstRound, height, TestingQBFTRootData, qbft.NoRound, nil)
}
var TestingRoundChangeMessageWithRoundAndHeight = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, height qbft.Height) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithParams(sk, id, round, height, TestingQBFTRootData, qbft.NoRound, nil)
}
var TestingRoundChangeMessageWithRoundAndRC = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, roundChangeJustification [][]byte) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, qbft.FirstRound, roundChangeJustification)
}
var TestingRoundChangeMessageWithParams = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	root [32]byte,
	dataRound qbft.Round,
	roundChangeJustification [][]byte,
) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               []byte{1, 2, 3, 4},
		Root:                     root,
		DataRound:                dataRound,
		RoundChangeJustification: roundChangeJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
var TestingMultiSignerRoundChangeMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingMultiSignerRoundChangeMessageWithParams(sks, ids, qbft.FirstRound, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingMultiSignerRoundChangeMessageWithRound = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingMultiSignerRoundChangeMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingMultiSignerRoundChangeMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingMultiSignerRoundChangeMessageWithParams(sks, ids, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingRoundChangeMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, WrongRoot, qbft.NoRound, nil)
}
var TestingMultiSignerRoundChangeMessageWithParams = func(sk []*bls.SecretKey, id []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := MultiSignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
