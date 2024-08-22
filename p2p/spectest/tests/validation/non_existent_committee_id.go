package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// NonExistentCommitteeID tests message with a non-existent validator
func NonExistentCommitteeID() tests.SpecTest {

	// Non-existent validator according to TestingMessageValidator().NetworkDataFetcher
	nonExistentValidator := testingutils.TestingNonExistentCommitteeID

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, nonExistentValidator[:], types.RoleCommittee),
			Data:    testingutils.EncodeQbftMessage(&qbft.Message{}),
		},
	}

	return &MessageValidationTest{
		Name:          "non existent committeeid",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrNonExistentCommitteeID.Error(),
	}
}
