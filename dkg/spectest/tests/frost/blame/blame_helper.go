package blame

import (
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/coinbase/kryptology/pkg/core/curves"
)

func makeInvalidForInvalidScalar(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	protocolMessage.Round1Message.ProofR = []byte("rubbish-value")
	d, _ := protocolMessage.Encode()
	return d
}

func makeInvalidForInvalidCommitment(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	protocolMessage.Round1Message.Commitment[len(protocolMessage.Round1Message.Commitment)-1] = []byte("rubbish-value")
	d, _ := protocolMessage.Encode()
	return d
}

func makeInvalidForInconsistentMessage(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	thisCurve := curves.BLS12381G1()
	lastCommitment, _ := thisCurve.NewIdentityPoint().FromAffineCompressed(protocolMessage.Round1Message.Commitment[len(protocolMessage.Round1Message.Commitment)-1])
	lastCommitment = lastCommitment.Double()
	protocolMessage.Round1Message.Commitment[len(protocolMessage.Round1Message.Commitment)-1] = lastCommitment.ToAffineCompressed()

	d, _ := protocolMessage.Encode()
	return d
}

func makeInvalidForInvalidShare(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	encryptedValue, _ := hex.DecodeString("0496849c3097292ab5d7326907216193660ea71d9bba9c58bcb17207579cc470f703c14933946f5cb540c3ecfffb533456ea516171e1f098614bc984e5bc13062b106711593e44458d86aa86800f6f5d81d645f2bc2c5ddae055b36040fae47b59a5cfe2881540d64a7d43")
	protocolMessage.Round1Message.Shares[1] = encryptedValue

	d, _ := protocolMessage.Encode()
	return d
}

func makeInvalidForInvalidShare_FailedDecrypt(data []byte) []byte {
	protocolMessage := &frost.ProtocolMsg{}
	_ = protocolMessage.Decode(data)

	protocolMessage.Round1Message.Shares[1] = []byte("rubbish value")

	d, _ := protocolMessage.Encode()
	return d
}
