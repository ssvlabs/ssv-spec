package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnexpectedPrepareJustificationInPrepareMessage tests a prepare message with a prepare justification
func UnexpectedPrepareJustificationInPrepareMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:              qbft.PrepareMsgType,
				Round:                qbft.FirstRound,
				Height:               qbft.FirstHeight,
				Root:                 testingutils.TestingQBFTRootData,
				Identifier:           msgID[:],
				PrepareJustification: [][]byte{testingutils.EncodeMessage(&types.SignedSSVMessage{})},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected prepare justification in prepare",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedPrepareJustifications.Error(),
	}
}

// UnexpectedPrepareJustificationInCommitMessage tests a commit message with a prepare justification
func UnexpectedPrepareJustificationInCommitMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:              qbft.CommitMsgType,
				Round:                qbft.FirstRound,
				Height:               qbft.FirstHeight,
				Root:                 testingutils.TestingQBFTRootData,
				Identifier:           msgID[:],
				PrepareJustification: [][]byte{testingutils.EncodeMessage(&types.SignedSSVMessage{})},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected prepare justification in commit",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedPrepareJustifications.Error(),
	}
}

// UnexpectedPrepareJustificationInRoundChangeMessage tests a round change message with a prepare justification
func UnexpectedPrepareJustificationInRoundChangeMessage() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:              qbft.RoundChangeMsgType,
				Round:                qbft.FirstRound,
				Height:               qbft.FirstHeight,
				Root:                 testingutils.TestingQBFTRootData,
				Identifier:           msgID[:],
				PrepareJustification: [][]byte{testingutils.EncodeMessage(&types.SignedSSVMessage{})},
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected prepare justification in round change",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedPrepareJustifications.Error(),
	}
}
