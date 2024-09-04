package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigMultipleSigners tests a partial signature message with multiple signers
func PartialSigMultipleSigners() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1, 2}, // Multiple signers
		Signatures:  [][]byte{testingutils.MockRSASignature[:], testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodePartialSignatureMessage(&types.PartialSignatureMessages{
				Type: types.PostConsensusPartialSig,
				Slot: 1,
				Messages: []*types.PartialSignatureMessage{
					{
						PartialSignature: testingutils.MockPartialSignature[:],
						SigningRoot:      [32]byte{1},
						Signer:           1,
						ValidatorIndex:   1,
					},
				},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "partial signature with multiple signers",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrPartialSigMultipleSigners.Error(),
	}
}
