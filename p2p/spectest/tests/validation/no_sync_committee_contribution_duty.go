package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoSyncCommitteeContributionDutyConsensusMessage tests a consensus message with a sync committee contribution role but the validator doesn't have the duty
// according to TestingMessageValidator().DutyFetcher
func NoSyncCommitteeContributionDutyConsensusMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPKWithoutSyncCommitteeContribution[:], types.RoleSyncCommitteeContribution)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Round:      qbft.FirstRound,
				Identifier: msgID[:],
				Root:       [32]byte{1},
				Height:     qbft.FirstHeight,
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "no sync committee contribution duty consensus message",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNoDuty.Error(),
	}
}

// NoSyncCommitteeContributionDutyPartialSignatureMessage tests a partial signature message with a sync committee contribution role but the validator doesn't have the duty
// according to TestingMessageValidator().DutyFetcher
func NoSyncCommitteeContributionDutyPartialSignatureMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPKWithoutSyncCommitteeContribution[:], types.RoleSyncCommitteeContribution)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodePartialSignatureMessage(&types.PartialSignatureMessages{
				Type: types.PostConsensusPartialSig,
				Slot: 1,
				Messages: []*types.PartialSignatureMessage{
					{
						PartialSignature: testingutils.MockPartialSignature[:],
						SigningRoot:      [32]byte{1},
						Signer:           1,
						ValidatorIndex:   testingutils.ValidatorIndexWithoutSyncCommitteeContributionDuty, // ValidatorIndex that doesn't have the duty
					},
				},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "no sync committee contribution duty partial signature message",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNoDuty.Error(),
	}
}
