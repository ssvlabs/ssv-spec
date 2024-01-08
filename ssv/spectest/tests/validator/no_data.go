package validator

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoData tests a validator that raises an error due to a message with no data
func NoData() tests.SpecTest {

	msgs := []*types.SSVMessage{
		{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester),
			Data:    []byte{},
		},
	}

	return &ValidatorTest{
		Name:                   "no data",
		Messages:               msgs,
		OutputMessages:         []*types.SSVMessage{},
		BeaconBroadcastedRoots: []string{},
		ExpectedError:          "Message invalid: msg data is invalid",
	}
}
