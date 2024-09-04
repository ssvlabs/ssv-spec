package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignerNotInCommittee tests message with a signer that doesn't belong to the validator's committee
func SignerNotInCommittee() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		// By default in TestingMessageValidator, TestingValidatorPK has operators 1, 2, 3 and 4
		OperatorIDs: []types.OperatorID{5}, // So, operator 5 doesn't belong to the committee
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer),
			Data:    testingutils.EncodeQbftMessage(&qbft.Message{}),
		},
	}

	return &MessageValidationTest{
		Name:          "signer not in committee",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrSignerNotInCommittee.Error(),
	}
}
