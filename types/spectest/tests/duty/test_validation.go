package duty

import (
	reflect2 "reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// ValidationTest runs ValidatorDuty.Validate on each duty, expecting the same error code for all.
type ValidationTest struct {
	Name              string
	Type              string
	Documentation     string
	ValidatorDuties   []*types.ValidatorDuty
	ExpectedErrorCode int
}

func (test *ValidationTest) TestName() string {
	return "validator duty validation " + test.Name
}

func (test *ValidationTest) Run(t *testing.T) {
	for _, duty := range test.ValidatorDuties {
		tests.AssertErrorCode(t, test.ExpectedErrorCode, duty.Validate())
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewValidationTest(name, documentation string, validatorDuties []*types.ValidatorDuty, expectedErrorCode int) *ValidationTest {
	return &ValidationTest{
		Name:              name,
		Type:              testdoc.ValidatorDutyValidationTestType,
		Documentation:     documentation,
		ValidatorDuties:   validatorDuties,
		ExpectedErrorCode: expectedErrorCode,
	}
}

func testingAttesterDuty() *types.ValidatorDuty {
	return testingutils.TestingAttesterDuty(spec.DataVersionElectra).ValidatorDuties[0]
}
