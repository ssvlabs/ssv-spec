package keysign

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/keysign"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

func HappyFlow() *tests.MsgProcessingSpecTest {
	ks := testingutils.TestingKeygenKeySet()
	network := testingutils.NewTestingNetwork()
	storage := testingutils.NewTestingStorage()
	keyManager := testingutils.NewTestingKeyManager()

	storage.SaveKeyGenOutput(&dkg.KeyGenOutput{
		Share:       ks.Shares[1],
		ValidatorPK: ks.ValidatorPK.Serialize(),
		Threshold:   ks.Threshold,
		OperatorPubKeys: map[types.OperatorID]*bls.PublicKey{
			1: ks.Shares[1].GetPublicKey(),
			2: ks.Shares[2].GetPublicKey(),
			3: ks.Shares[3].GetPublicKey(),
			4: ks.Shares[4].GetPublicKey(),
		},
	})

	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)

	signingRoot := testingutils.DespositDataSigningRoot(ks, testingutils.InitMessageData(
		[]types.OperatorID{1, 2, 3, 4},
		uint16(ks.Threshold),
		testingutils.TestingWithdrawalCredentials,
		testingutils.TestingForkVersion,
	))

	init := &dkg.KeySign{
		ValidatorPK: ks.ValidatorPK.Serialize(),
		SigningRoot: signingRoot,
	}
	initBytes, _ := init.Encode()

	testingNode := dkg.NewNode(
		&dkg.Operator{
			OperatorID:       1,
			ETHAddress:       ks.DKGOperators[1].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[1].EncryptionKey.PublicKey,
		},
		&dkg.Config{
			KeygenProtocol:      frost.New,
			ReshareProtocol:     frost.NewResharing,
			KeySign:             keysign.NewSignature,
			Network:             network,
			Storage:             storage,
			SignatureDomainType: types.PrimusTestnet,
			Signer:              keyManager,
		},
	)

	return &tests.MsgProcessingSpecTest{
		Name:        "keysign/happy flow",
		TestingNode: testingNode,
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.KeySignMsgType,
				Identifier: identifier,
				Data:       initBytes,
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_PreparationMessageBytes(2, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_PreparationMessageBytes(3, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_PreparationMessageBytes(4, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_Round1MessageBytes(2, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_Round1MessageBytes(3, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_Round1MessageBytes(4, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedKeySignOutputBytes(identifier, 2, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedKeySignOutputBytes(identifier, 3, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedKeySignOutputBytes(identifier, 4, signingRoot),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_PreparationMessageBytes(1, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       keysign.Testing_Round1MessageBytes(1, ks, signingRoot),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedKeySignOutputBytes(identifier, 1, signingRoot),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			1: ks.SignedKeySignOutputObject(identifier, 1, signingRoot),
			2: ks.SignedKeySignOutputObject(identifier, 2, signingRoot),
			3: ks.SignedKeySignOutputObject(identifier, 3, signingRoot),
			4: ks.SignedKeySignOutputObject(identifier, 4, signingRoot),
		},
		KeySet:        ks,
		ExpectedError: "",
	}
}
