package tests

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/stubdkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// HappyFlow tests a simple full happy flow until decided
func HappyFlow() *MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
	return &MsgProcessingSpecTest{
		Name:   "happy flow",
		KeySet: testingutils.Testing4SharesSet(),
		Messages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.InitMsgType,
				Identifier: identifier,
				Data:       testingutils.InitMessageDataBytes([]types.OperatorID{1, 2, 3, 4}, 3, testingutils.TestingWithdrawalCredentials),
			}),

			// stage 1
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage1),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage1),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage1),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage1),
			}),

			// stage 2
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage2),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage2),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage2),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage2),
			}),

			// stage 3
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage3),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[2].SK, 2, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage3),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[3].SK, 3, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage3),
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[4].SK, 4, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       testingutils.ProtocolMsgDataBytes(stubdkg.StubStage3),
			}),
		},
	}
}
