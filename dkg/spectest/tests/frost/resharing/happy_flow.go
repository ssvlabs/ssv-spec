package resharing

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/frost"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func HappyFlow() *tests.MsgProcessingSpecTest {
	ks := testingutils.TestingResharingKeySet()
	network := testingutils.NewTestingNetwork()
	storage := testingutils.NewTestingStorage()
	keyManager := testingutils.NewTestingKeyManager()

	identifier := dkg.NewRequestID(ks.DKGOperators[5].ETHAddress, 5)
	reshareBytes := testingutils.ReshareMessageDataBytes(
		[]types.OperatorID{5, 6, 7, 8}, // new committee
		uint16(ks.Threshold),
		types.ValidatorPK(ks.ValidatorPK.Serialize()),
	)

	testingNode := dkg.NewResharingNode(
		&dkg.Operator{
			OperatorID:       5,
			ETHAddress:       ks.DKGOperators[5].ETHAddress,
			EncryptionPubKey: &ks.DKGOperators[5].EncryptionKey.PublicKey,
		},
		[]types.OperatorID{1, 2, 3}, // old committee
		&dkg.Config{
			KeygenProtocol:  frost.New,
			ReshareProtocol: frost.NewResharing,
			Network:         network,
			Storage:         storage,
			Signer:          keyManager,
			// SignatureDomainType: sigDomainType,
		},
	)

	return &tests.MsgProcessingSpecTest{
		Name:        "resharing/happy flow",
		TestingNode: testingNode,
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.ReshareMsgType,
				Identifier: identifier,
				Data:       reshareBytes,
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[6].SK, 6, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.PreparationMessageBytes(6),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[7].SK, 7, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.PreparationMessageBytes(7),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[8].SK, 8, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.PreparationMessageBytes(8),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round1MessageBytes(1),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round1MessageBytes(2),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round1MessageBytes(3),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[6].SK, 6, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round2MessageBytes(6),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[7].SK, 7, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round2MessageBytes(7),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[8].SK, 8, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round2MessageBytes(8),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[6].SK, 6, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 6, nil),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[7].SK, 7, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 7, nil),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[8].SK, 8, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 8, nil),
			}),
		},
		OutputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.PreparationMessageBytes(5),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ResharingMsgStore.Round2MessageBytes(5),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[5].SK, 5, &dkg.Message{
				MsgType:    dkg.OutputMsgType,
				Identifier: identifier,
				Data:       ks.SignedOutputBytes(identifier, 5, nil),
			}),
		},
		Output: map[types.OperatorID]*dkg.SignedOutput{
			5: ks.SignedOutputObject(identifier, 5, nil),
			6: ks.SignedOutputObject(identifier, 6, nil),
			7: ks.SignedOutputObject(identifier, 7, nil),
			8: ks.SignedOutputObject(identifier, 8, nil),
		},
		KeySet:        ks,
		ExpectedError: "",
	}
}
