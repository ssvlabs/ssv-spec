package share

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// ValidationTest runs Share.Validate on each share, expecting the same error code for all.
type ValidationTest struct {
	Name              string
	Type              string
	Documentation     string
	Shares            []*types.Share
	ExpectedErrorCode int
}

func (test *ValidationTest) TestName() string {
	return "share validation " + test.Name
}

func (test *ValidationTest) Run(t *testing.T) {
	for _, share := range test.Shares {
		tests.AssertErrorCode(t, test.ExpectedErrorCode, share.Validate())
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewValidationTest(name, documentation string, shares []*types.Share, expectedErrorCode int) *ValidationTest {
	return &ValidationTest{
		Name:              name,
		Type:              testdoc.ShareValidationTestType,
		Documentation:     documentation,
		Shares:            shares,
		ExpectedErrorCode: expectedErrorCode,
	}
}
