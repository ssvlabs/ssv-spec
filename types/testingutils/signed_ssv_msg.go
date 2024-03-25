package testingutils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"

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
		OperatorID: []types.OperatorID{operatorID},
		Signature:  [][]byte{signature},
		SSVMessage: &ssvMsg,
	}
}

var SignedSSVMessageF = func(ssvMessage *types.SSVMessage, operatorID []types.OperatorID, sk []*rsa.PrivateKey) *types.SignedSSVMessage {

	if ssvMessage == nil {
		panic("Can't create a SignedSSVMessage with a nil SSVMessage")
	}

	data, err := ssvMessage.Encode()
	if err != nil {
		panic(err)
	}

	signatures := make([][]byte, len(operatorID))
	km := NewTestingKeyManager()
	for i, sk_i := range sk {
		pkByts, err := types.GetPublicKeyPem(sk_i)
		if err != nil {
			panic(err)
		}
		signature, err := km.SignNetworkData(data, pkByts)
		if err != nil {
			panic(err)
		}
		signatures[i] = signature
	}

	return &types.SignedSSVMessage{
		OperatorID: operatorID,
		Signature:  signatures,
		SSVMessage: ssvMessage,
	}
}
