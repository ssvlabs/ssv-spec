package resharing

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/hex"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func HappyFlow() *frost2.MsgProcessingSpecTest {
	ks := testingutils.TestingResharingKeySet()
	network := testingutils.NewTestingNetwork()
	storage := testingutils.NewTestingStorage()
	keyManager := testingutils.NewTestingKeyManager()

	identifier := dkg.NewRequestID(ks.DKGOperators[5].ETHAddress, 5)
	reshare := &dkg.Reshare{
		ValidatorPK: types.ValidatorPK(ks.ValidatorPK.Serialize()),
		OperatorIDs: []types.OperatorID{5, 6, 7, 8},
		Threshold:   3,
	}
	reshareBytes, _ := reshare.Encode()

	// vk := ks.ValidatorPK.Serialize()
	// root := func(validatorPK []byte, initMsg *dkg.Reshare) []byte {
	// 	root, _, _ := types.GenerateETHDepositData(
	// 		validatorPK,
	// 		types.DomainDeposit,
	// 	)
	// 	return root
	// }(vk, init)

	return &frost2.MsgProcessingSpecTest{
		Name: "happy flow",
		Operator: &dkg.Operator{
			OperatorID:       5,
			ETHAddress:       ks.DKGOperators[5].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[5].EncryptionKey.PublicKey,
		},
		IsResharing:  true,
		OperatorsOld: []types.OperatorID{1, 2, 3},
		Network:      network,
		Signer:       keyManager,
		Storage:      storage,
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg2(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.ReshareMsgType,
				Identifier: identifier,
				Data:       reshareBytes,
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[6].SK, 6, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(6),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[7].SK, 7, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(7),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[8].SK, 8, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(8),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round1MessageBytes(1),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round1MessageBytes(2),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round1MessageBytes(3),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[6].SK, 6, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(6),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[7].SK, 7, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(7),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[8].SK, 8, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(8),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[6].SK, 6, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 6, ks.DKGOperators[6].SK, &ks.DKGOperators[6].EncryptionKey.PublicKey, nil),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[7].SK, 7, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 7, ks.DKGOperators[7].SK, &ks.DKGOperators[7].EncryptionKey.PublicKey, nil),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[8].SK, 8, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 8, ks.DKGOperators[8].SK, &ks.DKGOperators[8].EncryptionKey.PublicKey, nil),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg2(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       PreparationMessageBytes(5),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(5),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 5, ks.DKGOperators[5].SK, &ks.DKGOperators[5].EncryptionKey.PublicKey, nil),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			5: SignedOutputObject(identifier, 5, ks.DKGOperators[5].SK, &ks.DKGOperators[5].EncryptionKey.PublicKey, nil),
			6: SignedOutputObject(identifier, 6, ks.DKGOperators[6].SK, &ks.DKGOperators[6].EncryptionKey.PublicKey, nil),
			7: SignedOutputObject(identifier, 7, ks.DKGOperators[7].SK, &ks.DKGOperators[7].EncryptionKey.PublicKey, nil),
			8: SignedOutputObject(identifier, 8, ks.DKGOperators[8].SK, &ks.DKGOperators[8].EncryptionKey.PublicKey, nil),
		},
		KeySet:        ks,
		ExpectedError: "",
	}
}

func Round1MessageBytes(id types.OperatorID) []byte {
	commitments := [][]byte{}
	for _, commitment := range testingutils.Resharing_Round1[id].Commitments {
		cbytes, _ := hex.DecodeString(commitment)
		commitments = append(commitments, cbytes)
	}
	proofS, _ := hex.DecodeString(testingutils.Resharing_Round1[id].ProofS)
	proofR, _ := hex.DecodeString(testingutils.Resharing_Round1[id].ProofR)
	shares := map[uint32][]byte{}
	for peerID, share := range testingutils.Resharing_Round1[id].Shares {
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
	pk, _ := hex.DecodeString(testingutils.Resharing_SessionPKs[id])
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
	vk, _ := hex.DecodeString(testingutils.Resharing_Round2[id].Vk)
	vkshare, _ := hex.DecodeString(testingutils.Resharing_Round2[id].VkShare)
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
	share := skFromHex(testingutils.Resharing_Round2[opId].SkShare)
	validatorPublicKey, _ := hex.DecodeString(testingutils.Resharing_Round2[5].Vk)

	o := &dkg.Output{
		RequestID:            requestID,
		EncryptedShare:       testingutils.TestingEncryption(opPK, share.Serialize()),
		SharePubKey:          share.GetPublicKey().Serialize(),
		ValidatorPubKey:      validatorPublicKey,
		DepositDataSignature: nil,
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
