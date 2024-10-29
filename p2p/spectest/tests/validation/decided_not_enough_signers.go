package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DecidedNotEnoughSigners tests a decided message with not enough signers
func DecidedNotEnoughSigners() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1, 2}, // Not enough signers (for TestingValidatorPK considering TestingMessageValidator().NetworkDataFetcher)
		Signatures:  [][]byte{testingutils.MockRSASignature[:], testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer),
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType: qbft.CommitMsgType,
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "decided without enough signers",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrDecidedNotEnoughSigners.Error(),
	}
}
