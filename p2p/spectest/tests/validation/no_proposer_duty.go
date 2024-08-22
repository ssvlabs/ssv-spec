package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NoProposerDutyConsensusMessage tests a consensus message with a proposer role but the validator doesn't have the duty
// according to TestingMessageValidator().DutyFetcher
func NoProposerDutyConsensusMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPKWithoutProposerDuty[:], types.RoleProposer)

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
		Name:          "no proposer duty consensus message",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNoDuty.Error(),
	}
}

// NoProposerDutyPartialSignatureMessage tests a partial signature message with a proposer role but the validator doesn't have the duty
// according to TestingMessageValidator().DutyFetcher
func NoProposerDutyPartialSignatureMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPKWithoutProposerDuty[:], types.RoleProposer)

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
						ValidatorIndex:   testingutils.ValidatorIndexWithoutProposerDuty,
					},
				},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "no proposer duty partial signature message",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNoDuty.Error(),
	}
}
