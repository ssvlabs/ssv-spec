package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// PrepareWithFullData tests a prepare message with a non-empty full data field
func PrepareWithFullData() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer),
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType: qbft.PrepareMsgType,
			}),
		},
		FullData: []byte{1, 2, 3, 4}, // Non-empty full data field
	}

	return &MessageValidationTest{
		Name:          "prepare with full data",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrPrepareOrCommitWithFullData.Error(),
	}
}
