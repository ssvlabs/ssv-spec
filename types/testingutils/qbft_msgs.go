package testingutils

import (
	"crypto/sha256"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingIdentifier = []byte{1, 2, 3, 4}

var DifferentFullData = append(TestingQBFTFullData, []byte("different")...)
var DifferentRoot = func() [32]byte {
	return sha256.Sum256(DifferentFullData)
}()

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
	return TestingProposalMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, nil, nil)
}
var TestingProposalMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingProposalMessageWithParams(sk, id, qbft.FirstRound, height, TestingQBFTRootData, nil, nil)
}
var TestingProposalMessageDifferentRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
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
var TestingProposalMessageWithRoundAndRC = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round, roundChangeJustification [][]byte) *qbft.SignedMessage {
	return TestingProposalMessageWithParams(sk, id, round, qbft.FirstHeight, TestingQBFTRootData, roundChangeJustification, nil)
}
var TestingProposalMessageWithIdentifierAndFullData = func(sk *bls.SecretKey, id types.OperatorID, identifier, fullData []byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier,
		Root:       sha256.Sum256(fullData),
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}
var TestingProposalMessageWithParams = func(
	sk *bls.SecretKey,
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
		Identifier:               TestingIdentifier,
		Root:                     root,
		RoundChangeJustification: roundChangeJustification,
		PrepareJustification:     prepareJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
var TestingMultiSignerProposalMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingMultiSignerProposalMessageWithParams(sks, ids, qbft.FirstRound, qbft.FirstHeight, TestingIdentifier, TestingQBFTFullData, TestingQBFTRootData)
}
var TestingMultiSignerProposalMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingMultiSignerProposalMessageWithParams(sks, ids, qbft.FirstRound, height, TestingIdentifier, TestingQBFTFullData, TestingQBFTRootData)
}
var TestingMultiSignerProposalMessageWithParams = func(
	sk []*bls.SecretKey,
	id []types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier,
	fullData []byte,
	root [32]byte,
) *qbft.SignedMessage {
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
var TestingPrepareMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingPrepareMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingPrepareMessageWithParams(sk, id, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingPrepareMessageWithParams(sk, id, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingIdentifier, DifferentRoot)
}
var TestingPrepareMessageWithFullData = func(sk *bls.SecretKey, id types.OperatorID, fullData []byte) *qbft.SignedMessage {
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
var TestingPrepareMessageWithParams = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier []byte,
	root [32]byte,
) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     height,
		Round:      round,
		Identifier: identifier,
		Root:       root,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
}
var TestingPrepareMultiSignerMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithRound(sks, ids, qbft.FirstRound)
}
var TestingPrepareMultiSignerMessageWithRound = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithHeightAndIdentifier = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height, identifier []byte) *qbft.SignedMessage {
	return TestingPrepareMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, identifier, TestingQBFTRootData)
}
var TestingPrepareMultiSignerMessageWithParams = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round, height qbft.Height, identifier []byte, root [32]byte) *qbft.SignedMessage {
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
var TestingCommitMessage = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithRound(sk, id, qbft.FirstRound)
}
var TestingCommitMessageWithRound = func(sk *bls.SecretKey, id types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData)
}
var TestingCommitMessageWithHeight = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData)
}
var TestingCommitMessageWrongRoot = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, TestingIdentifier, DifferentRoot)
}
var TestingCommitMessageWrongHeight = func(sk *bls.SecretKey, id types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMessageWithParams(sk, id, qbft.FirstRound, 10, TestingIdentifier, DifferentRoot)
}
var TestingCommitMessageWithIdentifierAndFullData = func(sk *bls.SecretKey, id types.OperatorID, identifier, fullData []byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: identifier,
		Root:       sha256.Sum256(fullData),
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = fullData
	return ret
}
var TestingCommitMessageWithParams = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier []byte,
	root [32]byte,
) *qbft.SignedMessage {
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
var TestingCommitMultiSignerMessage = func(sks []*bls.SecretKey, ids []types.OperatorID) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithRound(sks, ids, qbft.FirstRound)
}
var TestingCommitMultiSignerMessageWithRound = func(sks []*bls.SecretKey, ids []types.OperatorID, round qbft.Round) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, round, qbft.FirstHeight, TestingIdentifier, TestingQBFTRootData, TestingQBFTFullData)
}
var TestingCommitMultiSignerMessageWithHeight = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, TestingIdentifier, TestingQBFTRootData, TestingQBFTFullData)
}
var TestingCommitMultiSignerMessageWithHeightAndIdentifier = func(sks []*bls.SecretKey, ids []types.OperatorID, height qbft.Height, identifier []byte) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, height, identifier, TestingQBFTRootData, TestingQBFTFullData)
}
var TestingCommitMultiSignerMessageWithIdentifierAndFullData = func(sks []*bls.SecretKey, ids []types.OperatorID, identifier, fullData []byte) *qbft.SignedMessage {
	return TestingCommitMultiSignerMessageWithParams(sks, ids, qbft.FirstRound, qbft.FirstHeight, identifier, sha256.Sum256(fullData), fullData)
}
var TestingCommitMultiSignerMessageWithParams = func(
	sks []*bls.SecretKey,
	ids []types.OperatorID,
	round qbft.Round,
	height qbft.Height,
	identifier []byte,
	root [32]byte,
	fullData []byte,
) *qbft.SignedMessage {
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
var TestingRoundChangeMessageWithHeightAndIdentifier = func(sk *bls.SecretKey, id types.OperatorID, height qbft.Height, identifier []byte) *qbft.SignedMessage {
	msg := &qbft.Message{
		MsgType:    qbft.RoundChangeMsgType,
		Height:     height,
		Round:      qbft.FirstRound,
		Identifier: TestingIdentifier,
		Root:       TestingQBFTRootData,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
	return ret
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
		Identifier:               TestingIdentifier,
		Root:                     root,
		DataRound:                dataRound,
		RoundChangeJustification: roundChangeJustification,
	}
	ret := SignQBFTMsg(sk, id, msg)
	ret.FullData = TestingQBFTFullData
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
	return TestingRoundChangeMessageWithParams(sk, id, qbft.FirstRound, qbft.FirstHeight, DifferentRoot, qbft.NoRound, nil)
}
var TestingMultiSignerRoundChangeMessageWithParams = func(sk []*bls.SecretKey, id []types.OperatorID, round qbft.Round, height qbft.Height, root [32]byte) *qbft.SignedMessage {
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
