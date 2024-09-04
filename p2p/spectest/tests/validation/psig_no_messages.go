package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoPartialSignatureMessages tests a partial signature message with no PartialSignatureMessage
func NoPartialSignatureMessages() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodePartialSignatureMessage(&types.PartialSignatureMessages{
				Type:     types.PostConsensusPartialSig,
				Slot:     1,
				Messages: []*types.PartialSignatureMessage{}, // No messages
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "partial signature with no messages",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNoPartialSignatureMessages.Error(),
	}
}
