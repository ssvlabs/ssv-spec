package ssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgIDDoesntBelongs tests msg id doesn't belonging to validator id
func MsgIDDoesntBelongs() *SSVMessageTest {
	return NewSSVMessageTest(
		"does not belong",
		testdoc.SSVMessageTestDoesNotBelongDoc,
		[]types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingWrongValidatorPubKey[:], types.RoleUnknown),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingWrongValidatorPubKey[:], types.RoleCommittee),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingWrongValidatorPubKey[:], types.RoleUnknown),
		},
		false,
	)
}
