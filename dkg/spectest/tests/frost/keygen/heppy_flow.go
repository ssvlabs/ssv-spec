package keygen

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func HappyFlow() *tests.MsgProcessingSpecTest {
	ks := testingutils.TestingKeygenKeySet()
	network := testingutils.NewTestingNetwork()
	storage := testingutils.NewTestingStorage()
	keyManager := testingutils.NewTestingKeyManager()

	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
	init := testingutils.InitMessageData(
		[]types.OperatorID{1, 2, 3, 4},
		uint16(ks.Threshold),
		testingutils.TestingWithdrawalCredentials,
		testingutils.TestingForkVersion,
	)
	initBytes, _ := init.Encode()
	root := testingutils.DespositDataSigningRoot(ks, init)

	testingNode := dkg.NewNode(
		&dkg.Operator{
			OperatorID:       1,
			ETHAddress:       ks.DKGOperators[1].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[1].EncryptionKey.PublicKey,
		},
		&dkg.Config{
			KeygenProtocol:  frost.New,
			ReshareProtocol: frost.NewResharing,
			Network:         network,
			Storage:         storage,
			// SignatureDomainType: sigDomainType,
			Signer: keyManager,
		},
	)

	return &tests.MsgProcessingSpecTest{
		Name:        "keygen/happy flow",
		TestingNode: testingNode,
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: identifier,
				Data:       initBytes,
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_PreparationMessageBytes(2, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_PreparationMessageBytes(3, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_PreparationMessageBytes(4, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round1MessageBytes(2, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round1MessageBytes(3, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round1MessageBytes(4, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round2MessageBytes(2, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round2MessageBytes(3, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round2MessageBytes(4, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(2, root, ks.Shares[2]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(3, root, ks.Shares[3]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(4, root, ks.Shares[4]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 2, root),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 3, root),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 4, root),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_PreparationMessageBytes(1, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round1MessageBytes(1, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       frost.Testing_Round2MessageBytes(1, testingutils.KeygenMsgStore),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.DepositDataMsgType,
				Identifier: identifier,
				Data:       testingutils.PartialDepositDataBytes(1, root, ks.Shares[1]),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 1, root),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			1: ks.SignedOutputObject(identifier, 1, root),
			2: ks.SignedOutputObject(identifier, 2, root),
			3: ks.SignedOutputObject(identifier, 3, root),
			4: ks.SignedOutputObject(identifier, 4, root),
		},
		KeySet:        ks,
		ExpectedError: "",
	}
}
