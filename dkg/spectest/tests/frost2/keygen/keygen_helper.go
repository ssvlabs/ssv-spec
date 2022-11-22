package keygen

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func Round1MessageBytes(id types.OperatorID) []byte {
	commitments := [][]byte{}
	for _, commitment := range testingutils.Round1[id].Commitments {
		cbytes, _ := hex.DecodeString(commitment)
		commitments = append(commitments, cbytes)
	}
	proofS, _ := hex.DecodeString(testingutils.Round1[id].ProofS)
	proofR, _ := hex.DecodeString(testingutils.Round1[id].ProofR)
	shares := map[uint32][]byte{}
	for peerID, share := range testingutils.Round1[id].Shares {
		shareBytes, _ := hex.DecodeString(share)
		shares[peerID] = shareBytes
	}

	msg := frost.ProtocolMsg{
		Round: frost.Round1,
		Round1Message: &frost.Round1Message{
			Commitment: commitments,
			ProofS:     proofS,
			ProofR:     proofR,
			Shares:     shares,
		},
	}
	byts, _ := msg.Encode()
	return byts
}

func PreparationMessageBytes(id types.OperatorID) []byte {
	pk, _ := hex.DecodeString(testingutils.SessionPKs[id])
	msg := &frost.ProtocolMsg{
		Round: frost.Preparation,
		PreparationMessage: &frost.PreparationMessage{
			SessionPk: pk,
		},
	}
	byts, _ := msg.Encode()
	return byts
}

func Round2MessageBytes(id types.OperatorID) []byte {
	vk, _ := hex.DecodeString(testingutils.Round2[id].Vk)
	vkshare, _ := hex.DecodeString(testingutils.Round2[id].VkShare)
	msg := frost.ProtocolMsg{
		Round: frost.Round2,
		Round2Message: &frost.Round2Message{
			Vk:      vk,
			VkShare: vkshare,
		},
	}
	byts, _ := msg.Encode()
	return byts
}

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
	share := skFromHex(testingutils.Round2[opId].SkShare)
	validatorPublicKey, _ := hex.DecodeString(testingutils.Round2[1].Vk)

	sk := &bls.SecretKey{}
	_ = sk.DeserializeHexStr(testingutils.Round2[opId].Sk)

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
