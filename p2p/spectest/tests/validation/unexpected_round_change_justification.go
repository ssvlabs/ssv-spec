package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnexpectedRoundChangeJustificationInPrepareMessage tests a prepare message with a round change justification
func UnexpectedRoundChangeJustificationInPrepareMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:                  qbft.PrepareMsgType,
				Round:                    qbft.FirstRound,
				Height:                   qbft.FirstHeight,
				Root:                     testingutils.TestingQBFTRootData,
				Identifier:               msgID[:],
				RoundChangeJustification: [][]byte{testingutils.EncodeMessage(&types.SignedSSVMessage{})},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected round change justification in prepare",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedRoundChangeJustifications.Error(),
	}
}

// UnexpectedRoundChangeJustificationInCommitMessage tests a commit message with a round change justification
func UnexpectedRoundChangeJustificationInCommitMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:                  qbft.CommitMsgType,
				Round:                    qbft.FirstRound,
				Height:                   qbft.FirstHeight,
				Root:                     testingutils.TestingQBFTRootData,
				Identifier:               msgID[:],
				RoundChangeJustification: [][]byte{testingutils.EncodeMessage(&types.SignedSSVMessage{})},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected round change justification in commit",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedRoundChangeJustifications.Error(),
	}
}
