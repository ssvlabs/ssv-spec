package validator

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidDataPartialSig tests a validator that raises an error due to a partial signature message with invalid data
func InvalidDataPartialSig() tests.SpecTest {

	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester),
			Data:    []byte{1, 2, 3, 4},
		},
	}

	return &ValidatorTest{
		Name:                   "invalid data partial sig",
		Messages:               msgs,
		OutputMessages:         []*types.SSVMessage{},
		BeaconBroadcastedRoots: []string{},
		ExpectedError:          "could not get post consensus Message from network Message: incorrect size",
	}
}
