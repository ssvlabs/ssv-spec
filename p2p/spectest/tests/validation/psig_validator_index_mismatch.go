package validation

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PartialSigValidatorIndexMismatch tests a partial signature message with a ValidatorIndex that doesn't match the ValidatorPublicKey
func PartialSigValidatorIndexMismatch() tests.SpecTest {

	validatorPK := testingutils.TestingValidatorPK
	validatorIndex := phase0.ValidatorIndex(100) // By TestingMessageValidator().NetworkDataFetcher, 100 doesn't match with testingutils.TestingValidatorPK

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, validatorPK[:], types.RoleProposer)
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
						ValidatorIndex:   validatorIndex,
					},
				},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "partial signature with validator index mismatch",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrValidatorIndexMismatch.Error(),
	}
}

// PartialSigValidatorIndexMismatchForCommitteeRole tests a partial signature message with a ValidatorIndex that doesn't match the ValidatorPublicKey
// but for the committe role. For this role, no error is expected since we don't assume synchorny on the committees' validators sets
func PartialSigValidatorIndexMismatchForCommitteeRole() tests.SpecTest {

	// By TestingMessageValidator().NetworkDataFetcher, 100 doesn't belong to testingutils.TestingCommitteeID
	validatorIndex := phase0.ValidatorIndex(100)

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingCommitteeID[:], types.RoleCommittee)

	ssvMsg := &types.SSVMessage{
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
					ValidatorIndex:   validatorIndex,
				},
			},
		}),
	}

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.SignSSVMessage(testingutils.DefaultKeySet.OperatorKeys[1], ssvMsg)},
		SSVMessage:  ssvMsg,
	}

	return &MessageValidationTest{
		Name:          "partial signature with validator index mismatch for committee role",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: "", // No error is expected
	}
}
