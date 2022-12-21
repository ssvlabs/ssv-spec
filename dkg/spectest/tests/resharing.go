package tests

import (
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ResharingHappyFlow tests a simple (dummy) resharing flow, the difference between this and keygen happy flow is
// resharing doesn't sign deposit data
func ResharingHappyFlow() *MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	identifier := dkg.NewRequestID(ks.DKGOperators[1].ETHAddress, 1)
	reshare := &dkg.Reshare{
		ValidatorPK: make([]byte, 48),
		OperatorIDs: []types.OperatorID{1, 2, 3, 4},
		Threshold:   3,
	}
	reshareBytes, _ := reshare.Encode()
	var root []byte

	return &MsgProcessingSpecTest{
		Name: "resharing happy flow",
		InputMessages: []*dkg.SignedMessage{
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ReshareMsgType,
				Identifier: identifier,
				Data:       reshareBytes,
			}),
			testingutils.SignDKGMsg(ks.DKGOperators[1].SK, 1, &dkg.Message{
				MsgType:    dkg.ProtocolMsgType,
				Identifier: identifier,
				Data:       nil, // GLNOTE: Dummy message simulating the Protocol to complete
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
