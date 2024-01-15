package ssvmsg

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgIDBelongs tests msg id belonging to validator id
func MsgIDBelongs() *SSVMessageTest {
	return &SSVMessageTest{
		Name: "belongs",
		MessageIDs: []types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.BeaconRole(100)),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.BNRoleAttester),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.BeaconRole(100)),
		},
		BelongsToValidator: true,
	}
}
