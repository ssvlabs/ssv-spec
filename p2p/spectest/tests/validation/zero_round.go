package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ZeroRound tests a consensus message with round 0
func ZeroRound() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer),
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType: qbft.PrepareMsgType,
				Round:   0,
			}),
		},
	}

	return &MessageValidationTest{
		Name:          "zero round",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrZeroRound.Error(),
	}
}
