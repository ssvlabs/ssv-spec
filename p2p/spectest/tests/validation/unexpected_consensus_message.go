package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnexpectedConsensusMessageForValidatorRegistration tests a consensus message for the validator registration role
func UnexpectedConsensusMessageForValidatorRegistration() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleValidatorRegistration)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Round:      qbft.FirstRound,
				Height:     qbft.FirstHeight,
				Root:       testingutils.TestingQBFTRootData,
				Identifier: msgID[:],
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected consensus message for validator registration",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedConsensusMessage.Error(),
	}
}

// UnexpectedConsensusMessageForVoluntaryExit tests a consensus message for the voluntary exit role
func UnexpectedConsensusMessageForVoluntaryExit() tests.SpecTest {

	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleVoluntaryExit)

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   msgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Round:      qbft.FirstRound,
				Height:     qbft.FirstHeight,
				Root:       testingutils.TestingQBFTRootData,
				Identifier: msgID[:],
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "unexpected consensus message for voluntary exit",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnexpectedConsensusMessage.Error(),
	}
}
