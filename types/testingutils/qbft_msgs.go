package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

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

/**
Proposal messages
*/
var TestingProposalMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingProposalMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingProposalMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       TestingQBFTRootData,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
var TestingMultiSignerProposalMessage = func(sk []*bls.SecretKey, id []types.OperatorID) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Root:       TestingQBFTRootData,
	}
	ret := MultiSignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}

/**
Prepare messages
*/
var TestingPrepareMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingPrepareMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     qbft.FirstHeight,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       TestingQBFTRootData,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}

/**
Commit messages
*/
var TestingCommitMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingCommitMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingCommitMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 5, 6, 7, 8, 9})
}
var TestingCommitMessageWrongHeight = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, 10, [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 5, 6, 7, 8, 9})
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
	return TestingCommitMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingCommitMultiSignerMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, TestingQBFTRootData)
}
var TestingCommitMultiSignerMessageWithParams = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
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

/**
Round Change messages
*/
var TestingRoundChangeMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingRoundChangeMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData)
}
var TestingRoundChangeMessageWithParams = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     height,
		Round:      round,
		Identifier: []byte{1, 2, 3, 4},
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	return ret
}
