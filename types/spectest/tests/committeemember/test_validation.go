package committeemember

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// ValidationTest runs CommitteeMember.Validate on each committee member, expecting the same error code for all.
type ValidationTest struct {
	Name              string
	Type              string
	Documentation     string
	CommitteeMembers  []*types.CommitteeMember
	ExpectedErrorCode int
}

func (test *ValidationTest) TestName() string {
	return "committee member validation " + test.Name
}

func (test *ValidationTest) Run(t *testing.T) {
	for _, cm := range test.CommitteeMembers {
		tests.AssertErrorCode(t, test.ExpectedErrorCode, cm.Validate())
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewValidationTest(name, documentation string, committeeMembers []*types.CommitteeMember, expectedErrorCode int) *ValidationTest {
	return &ValidationTest{
		Name:              name,
		Type:              testdoc.CommitteeMemberValidationTestType,
		Documentation:     documentation,
		CommitteeMembers:  committeeMembers,
		ExpectedErrorCode: expectedErrorCode,
	}
}
