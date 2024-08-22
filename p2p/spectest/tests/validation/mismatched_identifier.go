package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MismatchedIdentifier tests a consensus message with a mismatched identifier
func MismatchedIdentifier() tests.SpecTest {

	signature := [validation.RsaSignatureSize]byte{1, 2, 3, 4}

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{signature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer), // First identifier
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Round:      qbft.FirstRound,
				Identifier: []byte{1}, // Different identifier
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "mismatched identifier",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrMismatchedIdentifier.Error(),
	}
}
