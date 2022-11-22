package blame

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	ecies "github.com/ecies/go/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func BlameTypeInvalidScaler_HappyFlow() *frost2.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	network := testingutils.NewTestingNetwork()
	storage := testingutils.NewTestingStorage()
	keyManager := testingutils.NewTestingKeyManager()

	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
	init := &dkg.Init{
		OperatorIDs:           []types.OperatorID{1, 2, 3, 4},
		Threshold:             3,
		WithdrawalCredentials: testingutils.TestingWithdrawalCredentials,
		Fork:                  testingutils.TestingForkVersion,
	}
	initBytes, _ := init.Encode()

	return &frost2.MsgProcessingSpecTest{
		Name: "happy flow",
		Operator: &dkg.Operator{
			OperatorID:       1,
			ETHAddress:       ks.DKGOperators[1].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[1].EncryptionKey.PublicKey,
		},
		Network: network,
		Signer:  keyManager,
		Storage: storage,
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg2(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: identifier,
				Data:       initBytes,
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(2),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(3),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(4),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       makeInvalidForInvalidScalar(Round1MessageBytes(2)),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg2(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(1),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round1MessageBytes(1),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data: BlameMessageBytes(2, testingutils.SignDKGMsg2(ks.DKGOperators[2].SK, 2, &dkg.Message{
					MsgType:    dkg.ProtocolMsgType,
					Identifier: identifier,
					Data:       makeInvalidForInvalidScalar(Round1MessageBytes(2)),
				})),
			}),
		},
		Output:        map[types.OperatorID]*dkg.SignedOutput{},
		KeySet:        ks,
		ExpectedError: "",
	}
}

func round1Message(id types.OperatorID) *frost.ProtocolMsg {
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

	return &frost.ProtocolMsg{
		Round: frost.Round1,
		Round1Message: &frost.Round1Message{
			Commitment: commitments,
			ProofS:     proofS,
			ProofR:     proofR,
			Shares:     shares,
		},
	}
}

func Round1MessageBytes(id types.OperatorID) []byte {
	msg := round1Message(id)
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

func BlameMessageBytes(id types.OperatorID, blameMessage *dkg.SignedMessage) []byte {
	blameMessageBytes, _ := blameMessage.Encode()
	skBytes, _ := hex.DecodeString(testingutils.Resharing_SessionSKs[5])
	sk := ecies.NewPrivateKeyFromBytes(skBytes)
	ret, _ := (&frost.ProtocolMsg{
		Round: frost.Blame,
		BlameMessage: &frost.BlameMessage{
			Type:             frost.InvalidMessage,
			TargetOperatorID: uint32(id),
			BlameData:        [][]byte{blameMessageBytes},
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
