package ssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgIDBelongs tests msg id belonging to validator id
func MsgIDBelongs() *SSVMessageTest {
	return NewSSVMessageTest(
		"belongs",
		"Test that message IDs with matching validator public key belong to the validator",
		[]types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleUnknown),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleUnknown),
		},
		true,
	)
}
