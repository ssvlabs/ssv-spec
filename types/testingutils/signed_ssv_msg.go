package testingutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingSignedSSVMessageSignature = []byte{1, 2, 3, 4}

var TestingMessageID = types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCommittee)

var TestingSignedSSVMessage = func(sk *bls.SecretKey, operatorID types.OperatorID, rsaSK *rsa.PrivateKey) *types.SignedSSVMessage {
	// SignedPartialSigMessage
	signedPartialSig := PreConsensusSelectionProofMsg(sk, sk, operatorID, operatorID)
	signedPartialSigByts, err := signedPartialSig.Encode()
	if err != nil {
		panic(err.Error())
	}

	// SSVMessage
	ssvMsg := types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   TestingMessageID,
		Data:    signedPartialSigByts[:],
	}
	ssvMsgByts, err := ssvMsg.Encode()
	if err != nil {
		panic(err.Error())
	}

	// Sign message
	hash := sha256.Sum256(ssvMsgByts)
	signature, err := rsa.SignPKCS1v15(nil, rsaSK, crypto.SHA256, hash[:])
	if err != nil {
		panic(err.Error())
	}

	//SignedSSVMessage
	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{operatorID},
		Signatures:  [][]byte{signature},
		SSVMessage:  &ssvMsg,
	}
}

var SignedSSVMessageListF = func(ks *TestKeySet, signers []types.OperatorID, ssvMessages []*types.SSVMessage) []*types.SignedSSVMessage {
	ret := make([]*types.SignedSSVMessage, 0)
	for i, msg := range ssvMessages {
		ret = append(ret, SignedSSVMessageWithSigner(signers[i], ks.OperatorKeys[signers[i]], msg))
	}
	return ret
}

var SignPartialSigSSVMessage = func(ks *TestKeySet, msg *types.SSVMessage) *types.SignedSSVMessage {

	// Discover message's signer
	if msg.MsgType != types.SSVPartialSignatureMsgType {
		panic("type different than SSVPartialSignatureMsgType to sign partial signature ssv message")
	}

	psigMsgs := &types.PartialSignatureMessages{}
	if err := psigMsgs.Decode(msg.Data); err != nil {
		panic(err)
	}
	var signer types.OperatorID
	if len(psigMsgs.Messages) == 0 {
		signer = 1
	} else {
		signer = psigMsgs.Messages[0].Signer
	}

	// Convert SSVMessage to SignedSSVMessage
	return SignedSSVMessageWithSigner(signer, ks.OperatorKeys[signer], msg)
}

var SignedSSVMessageWithSigner = func(operatorID types.OperatorID, rsaSK *rsa.PrivateKey, ssvMessage *types.SSVMessage) *types.SignedSSVMessage {

	data, err := ssvMessage.Encode()
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaSK, crypto.SHA256, hash[:])
	if err != nil {
		panic(err)
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{operatorID},
		Signatures:  [][]byte{signature},
		SSVMessage:  ssvMessage,
	}
}
