package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigWithFullData tests a partial signaturr message with a non-empty full data field
func PartialSigWithFullData() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
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
						Signer:           1,
						ValidatorIndex:   1,
					},
				},
			}),
		},
		FullData: []byte{1, 2, 3, 4}, // Non-empty full data
	}

	return &MessageValidationTest{
		Name:          "partial signature with full data",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrFullDataNotInConsensusMessage.Error(),
	}
}
