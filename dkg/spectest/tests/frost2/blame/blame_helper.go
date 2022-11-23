package blame

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/frost/frostutils"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/coinbase/kryptology/pkg/core/curves"
	ecies "github.com/ecies/go/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

var skFromHex = func(str string) *bls.SecretKey {
	types.InitBLS()
	ret := &bls.SecretKey{}
	if err := ret.DeserializeHexStr(str); err != nil {
		panic(err.Error())
	}
	return ret
}

func SignedOutputBytes(requestID dkg.RequestID, opId types.OperatorID, sk *ecdsa.PrivateKey, pk *rsa.PublicKey, root []byte) []byte {
	d := SignedOutputObject(requestID, opId, sk, pk, root)
	ret, _ := d.Encode()
	return ret
}

func SignedOutputObject(requestID dkg.RequestID, opId types.OperatorID, opSK *ecdsa.PrivateKey, opPK *rsa.PublicKey, root []byte) *dkg.SignedOutput {
	share := skFromHex(frostutils.KeygenMsgStore.Round2[opId].SkShare)
	validatorPublicKey, _ := hex.DecodeString(frostutils.KeygenMsgStore.Round2[1].Vk)

	sk := &bls.SecretKey{}
	_ = sk.DeserializeHexStr(frostutils.KeygenMsgStore.Round2[opId].Sk)

	o := &dkg.Output{
		RequestID:            requestID,
		EncryptedShare:       testingutils.TestingEncryption(opPK, share.Serialize()),
		SharePubKey:          share.GetPublicKey().Serialize(),
		ValidatorPubKey:      validatorPublicKey,
		DepositDataSignature: sk.SignByte(root).Serialize(),
	}

	root1, _ := o.GetRoot()

	sig, _ := crypto.Sign(root1, opSK)

	ret := &dkg.SignedOutput{
		Data:      o,
		Signer:    opId,
		Signature: sig,
	}
	return ret
}

func BlameMessageBytes(id types.OperatorID, blameType frost.BlameType, blameMessages []*dkg.SignedMessage) []byte {
	blameData := make([][]byte, 0)
	for _, blameMessage := range blameMessages {
		byts, _ := blameMessage.Encode()
		blameData = append(blameData, byts)
	}

	skBytes, _ := hex.DecodeString(frostutils.KeygenMsgStore.SessionSKs[1])
	sk := ecies.NewPrivateKeyFromBytes(skBytes)

	ret, _ := (&frost.ProtocolMsg{
		Round: frost.Blame,
		BlameMessage: &frost.BlameMessage{
			Type:             blameType,
			TargetOperatorID: uint32(id),
			BlameData:        blameData,
			BlamerSessionSk:  sk.Bytes(),
		},
	}).Encode()
	return ret
}

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
