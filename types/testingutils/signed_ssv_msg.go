package testingutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingSignedSSVMessageSignature = []byte{1, 2, 3, 4}

var TestingMessageID = types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleAttester)

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
		OperatorID: operatorID,
		Signature:  signature,
		Data:       ssvMsgByts,
	}
}

var SignedSSVMessageListF = func(ks *TestKeySet, ssvMessages []*types.SSVMessage) []*types.SignedSSVMessage {
	ret := make([]*types.SignedSSVMessage, 0)
	for _, msg := range ssvMessages {
		ret = append(ret, SignedSSVMessageF(ks, msg))
	}
	return ret
}

var SignedSSVMessageF = func(ks *TestKeySet, msg *types.SSVMessage) *types.SignedSSVMessage {

	// Discover message's signer
	signer := types.OperatorID(1)
	if msg.MsgType == types.SSVConsensusMsgType {
		signedMsg := &qbft.SignedMessage{}
		if err := signedMsg.Decode(msg.Data); err != nil {
			panic(err)
		}
		signer = signedMsg.Signers[0]
	} else if msg.MsgType == types.SSVPartialSignatureMsgType {
		signedPartial := &types.SignedPartialSignatureMessage{}
		if err := signedPartial.Decode(msg.Data); err != nil {
			panic(err)
		}
		signer = signedPartial.Signer
	} else {
		panic("unknown type")
	}

	// Convert SSVMessage to SignedSSVMessage
	return SignedSSVMessageWithSigner(signer, ks.SSVKeys[signer], msg)
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
		OperatorID: operatorID,
		Signature:  signature,
		Data:       data,
	}
}
