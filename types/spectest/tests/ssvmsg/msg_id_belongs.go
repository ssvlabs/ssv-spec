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
			types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingValidatorPubKey), types.RoleProposer),
			types.NewValidatorMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, types.ValidatorPK(testingutils.TestingValidatorPubKey), types.RoleProposer),
			types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingValidatorPubKey), types.RoleValidatorRegistration),
			types.NewValidatorMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, types.ValidatorPK(testingutils.TestingValidatorPubKey), types.RoleValidatorRegistration),
		},
		true,
	)
}
