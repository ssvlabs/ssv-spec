package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidHash tests a consensus message with an invalid hash
func InvalidHash() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer),
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType: qbft.ProposalMsgType,
				Root:    [32]byte{1}, // Invalid hash for FullData
			}),
		},
		FullData: testingutils.TestingQBFTFullData,
	}

	return &MessageValidationTest{
		Name:          "invalid hash",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrInvalidHash.Error(),
	}
}
