package keygen

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

	vk, _ := hex.DecodeString(testingutils.Round2[1].Vk)
	root := func(validatorPK []byte, initMsg *dkg.Init) []byte {
		root, _, _ := types.GenerateETHDepositData(
			validatorPK,
			initMsg.WithdrawalCredentials,
			initMsg.Fork,
			types.DomainDeposit,
		)
		return root
	}(vk, init)

	return &frost2.MsgProcessingSpecTest{
		Name: "happy flow",
		Operator: &dkg.Operator{
			OperatorID:       1,
			ETHAddress:       ks.DKGOperators[1].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[1].EncryptionKey.PublicKey,
		},
		NodeConfig: &dkg.Config{
			KeygenProtocol:  frost.New,
			ReshareProtocol: frost.NewResharing,
			Network:         network,
			Signer:          keyManager,
			Storage:         storage,
		},
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
				Data:       Round1MessageBytes(2),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round1MessageBytes(3),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round1MessageBytes(4),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(2),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(3),
			}),
			testingutils.SignDKGMsg2(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       Round2MessageBytes(4),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(2, root, skFromHex(testingutils.Round2[2].SkShare)),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(3, root, skFromHex(testingutils.Round2[3].SkShare)),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(4, root, skFromHex(testingutils.Round2[4].SkShare)),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 2, ks.DKGOperators[2].SK, &ks.DKGOperators[2].EncryptionKey.PublicKey, root),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 3, ks.DKGOperators[3].SK, &ks.DKGOperators[3].EncryptionKey.PublicKey, root),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 4, ks.DKGOperators[4].SK, &ks.DKGOperators[4].EncryptionKey.PublicKey, root),
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
				Data:       Round2MessageBytes(1),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(1, root, skFromHex(testingutils.Round2[1].SkShare)),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       SignedOutputBytes(identifier, 1, ks.DKGOperators[1].SK, &ks.DKGOperators[1].EncryptionKey.PublicKey, root),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			1: SignedOutputObject(identifier, 1, ks.DKGOperators[1].SK, &ks.DKGOperators[1].EncryptionKey.PublicKey, root),
			2: SignedOutputObject(identifier, 2, ks.DKGOperators[2].SK, &ks.DKGOperators[2].EncryptionKey.PublicKey, root),
			3: SignedOutputObject(identifier, 3, ks.DKGOperators[3].SK, &ks.DKGOperators[3].EncryptionKey.PublicKey, root),
			4: SignedOutputObject(identifier, 4, ks.DKGOperators[4].SK, &ks.DKGOperators[4].EncryptionKey.PublicKey, root),
		},
		KeySet:        ks,
		ExpectedError: "",
	}
}

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
