package testingutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/bloxapp/ssv-spec/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SignedSSVMessage
	OperatorID      types.OperatorID
	OperatorSK      *rsa.PrivateKey
}

func NewTestingNetwork(operatorID types.OperatorID, sk *rsa.PrivateKey) *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SignedSSVMessage, 0),
		OperatorID:      operatorID,
		OperatorSK:      sk,
	}
}

func (net *TestingNetwork) Broadcast(message *types.SSVMessage) error {

	encodedMessage, err := message.Encode()
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(encodedMessage)

	signature, err := rsa.SignPKCS1v15(rand.Reader, net.OperatorSK, crypto.SHA256, hash[:])
	if err != nil {
		panic(err)
	}

	signedMessage := &types.SignedSSVMessage{
		OperatorID: net.OperatorID,
		Signature:  signature,
		Data:       encodedMessage,
	}

	net.BroadcastedMsgs = append(net.BroadcastedMsgs, signedMessage)
	return nil
}

func ConvertBroadcastedMessagesToSSVMessages(signedMessages []*types.SignedSSVMessage) []*types.SSVMessage {
	ret := make([]*types.SSVMessage, 0)
	for _, msg := range signedMessages {
		ssvMsg := &types.SSVMessage{}
		err := ssvMsg.Decode(msg.Data)
		if err != nil {
			panic(err)
		}
		ret = append(ret, ssvMsg)
	}
	return ret
}
