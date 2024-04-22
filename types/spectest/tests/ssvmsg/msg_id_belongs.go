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
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleUnknown),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleUnknown),
		},
		BelongsToValidator: true,
	}
}
