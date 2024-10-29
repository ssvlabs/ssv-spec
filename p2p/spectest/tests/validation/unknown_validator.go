package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownValidator tests message with an unknown validator
func UnknownValidator() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingNonExistentValidator[:], types.RoleProposer),
			Data:    testingutils.EncodeQbftMessage(&qbft.Message{}),
		},
	}

	return &MessageValidationTest{
		Name:          "unknown validator",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrUnknownValidator.Error(),
	}
}
