package share

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// UnknownDomain tests a share with a domain type not defined in the spec
func UnknownDomain() *ValidationTest {
	share := testingutils.TestingShare(testingutils.Testing4SharesSet(), testingutils.TestingValidatorIndex)
	share.DomainType = testingutils.TestingUnknownDomainType

	return NewValidationTest(
		"unknown domain",
		testdoc.ShareValidationUnknownDomainDoc,
		[]*types.Share{share},
		types.InvalidShareErrorCode,
	)
}
