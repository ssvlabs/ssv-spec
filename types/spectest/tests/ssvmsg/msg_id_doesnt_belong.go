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
			types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingWrongValidatorPubKey), types.RoleProposer),
			types.NewValidatorMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, types.ValidatorPK(testingutils.TestingWrongValidatorPubKey), types.RoleProposer),
			types.NewValidatorMsgID(testingutils.TestingSSVDomainType, types.ValidatorPK(testingutils.TestingWrongValidatorPubKey), types.RoleValidatorRegistration),
			types.NewValidatorMsgID(types.DomainType{0x99, 0x99, 0x99, 0x99}, types.ValidatorPK(testingutils.TestingWrongValidatorPubKey), types.RoleValidatorRegistration),
		},
		false,
	)
}
