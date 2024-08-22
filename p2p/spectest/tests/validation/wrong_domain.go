package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongDomain tests message with a wrong domain
func WrongDomain() tests.SpecTest {

	wrongDomain := types.DomainType([4]byte{9, 9, 9, 9}) // bad domain

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(wrongDomain, testingutils.TestingValidatorPK[:], types.RoleProposer),
			Data:    testingutils.EncodeQbftMessage(&qbft.Message{}),
		},
	}

	return &MessageValidationTest{
		Name:          "wrong domain",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrWrongDomain.Error(),
	}
}
