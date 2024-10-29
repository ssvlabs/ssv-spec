package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigInconsistentSigners tests a partial signature message with inconsistent signers
func PartialSigInconsistentSigners() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1}, // Main signer
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
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
						Signer:           2, // Inconsistent signer
						ValidatorIndex:   1,
					},
				},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "partial signature with inconsistent signers",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrInconsistentSigners.Error(),
	}
}
