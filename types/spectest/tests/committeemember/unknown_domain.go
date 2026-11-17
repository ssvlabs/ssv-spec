package committeemember

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownDomain tests a committee member with a domain type not defined in the spec
func UnknownDomain() *ValidationTest {
	cm := testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet())
	cm.DomainType = testingutils.TestingUnknownDomainType

	return NewValidationTest(
		"unknown domain",
		testdoc.CommitteeMemberValidationUnknownDomainDoc,
		[]*types.CommitteeMember{cm},
		types.InvalidCommitteeMemberErrorCode,
	)
}
