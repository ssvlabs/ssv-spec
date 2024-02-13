package ssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgIDDoesntBelongs tests msg id doesn't belonging to validator id
func MsgIDDoesntBelongs() *SSVMessageTest {
	return &SSVMessageTest{
		Name: "does not belong",
		MessageIDs: []types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], types.BNRoleAttester),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], types.BeaconRole(100)),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingWrongValidatorPubKey[:], types.BNRoleAttester),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingWrongValidatorPubKey[:], types.BeaconRole(100)),
		},
		BelongsToValidator: false,
	}
}
