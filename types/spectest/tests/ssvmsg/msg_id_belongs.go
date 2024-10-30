package ssvmsg

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MsgIDBelongs tests msg id belonging to validator id
func MsgIDBelongs() *SSVMessageTest {
	return &SSVMessageTest{
		Name: "belongs",
		MessageIDs: []types.MessageID{
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleProposer),
			types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPubKey[:], types.RoleAggregator),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleProposer),
			types.NewMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, testingutils.TestingValidatorPubKey[:], types.RoleAggregator),
		},
		BelongsToValidator: true,
	}
}
