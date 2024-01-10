package testingutils

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var TestingSignedSSVMessageSignature = []byte{1, 2, 3, 4}

var TestingMessageID = types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.BNRoleAttester)

var TestingSignedSSVMessage = func(sk *bls.SecretKey, operatorID types.OperatorID) *types.SignedSSVMessage {
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

	//SignedSSVMessage
	return &types.SignedSSVMessage{
		OperatorID: operatorID,
		Signature:  TestingSignedSSVMessageSignature,
		Data:       ssvMsgByts,
	}
}
