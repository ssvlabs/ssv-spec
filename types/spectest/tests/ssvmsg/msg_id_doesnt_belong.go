package ssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgIDDoesntBelongs tests msg id doesn't belonging to validator id
func MsgIDDoesntBelongs() *SSVMessageTest {
	return &SSVMessageTest{
		Name: "does not belong",
		MessageIDs: []types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], types.RoleUnknown),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingWrongValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingWrongValidatorPubKey[:], types.RoleUnknown),
		},
		BelongsToValidator: false,
	}
}
