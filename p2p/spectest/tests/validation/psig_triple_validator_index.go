package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigTripleValidatorIndex tests a partial signature message that has the same validator index for 3 signatures
func PartialSigTripleValidatorIndex() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingCommitteeID[:], types.RoleCommittee)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodePartialSignatureMessage(&types.PartialSignatureMessages{
				Type: types.PostConsensusPartialSig,
				Slot: 0,
				Messages: []*types.PartialSignatureMessage{
					{
						PartialSignature: testingutils.MockPartialSignature[:],
						SigningRoot:      [32]byte{1},
						Signer:           1,
						ValidatorIndex:   1,
					},
					{
						PartialSignature: testingutils.MockPartialSignature[:],
						SigningRoot:      [32]byte{2},
						Signer:           1,
						ValidatorIndex:   1,
					},
					{
						PartialSignature: testingutils.MockPartialSignature[:],
						SigningRoot:      [32]byte{3},
						Signer:           1,
						ValidatorIndex:   1,
					},
				},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "partial signature with triple validator index",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrTripleValidatorIndexInPartialSignatures.Error(),
	}
}
