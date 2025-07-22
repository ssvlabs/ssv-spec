package ssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgIDBelongs tests msg id belonging to validator id
func MsgIDBelongs() *SSVMessageTest {
	return NewSSVMessageTest(
		"belongs",
		testdoc.SSVMessageTestBelongsDoc,
		[]types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleUnknown),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleUnknown),
		},
		true,
	)
}
