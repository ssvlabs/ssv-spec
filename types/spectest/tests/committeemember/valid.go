package committeemember

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid committee member for every known domain type
func Valid() *ValidationTest {
	ks := testingutils.Testing4SharesSet()

	committeeMembers := make([]*types.CommitteeMember, 0, len(testingutils.TestingKnownDomainTypes))
	for _, domain := range testingutils.TestingKnownDomainTypes {
		cm := testingutils.TestingCommitteeMember(ks)
		cm.DomainType = domain
		committeeMembers = append(committeeMembers, cm)
	}

	return NewValidationTest(
		"valid",
		testdoc.CommitteeMemberValidationValidDoc,
		committeeMembers,
		0,
	)
}
