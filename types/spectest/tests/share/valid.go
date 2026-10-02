package share

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid share for every known domain type
func Valid() *ValidationTest {
	ks := testingutils.Testing4SharesSet()

	shares := make([]*types.Share, 0, len(testingutils.TestingKnownDomainTypes))
	for _, domain := range testingutils.TestingKnownDomainTypes {
		share := testingutils.TestingShare(ks, testingutils.TestingValidatorIndex)
		share.DomainType = domain
		shares = append(shares, share)
	}

	return NewValidationTest(
		"valid",
		testdoc.ShareValidationValidDoc,
		shares,
		0,
	)
}
