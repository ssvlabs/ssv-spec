package resharing

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/spectest/tests/frost2"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
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

	return &frost2.MsgProcessingSpecTest{
		Name: "resharing/happy flow",
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
