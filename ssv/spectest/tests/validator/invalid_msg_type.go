package validator

import (
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func InvalidType() tests.SpecTest {

	msgs := []*types.SSVMessage{
		{
			MsgType: types.MsgType(100),
			MsgID:   types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester),
			Data:    []byte{1},
		},
	}

	return &ValidatorTest{
		Name:                   "invalid type",
		Messages:               msgs,
		OutputMessages:         []*types.SSVMessage{},
		BeaconBroadcastedRoots: []string{},
		ExpectedError:          "unknown msg",
	}
}
