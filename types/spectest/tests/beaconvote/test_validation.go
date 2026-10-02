package beaconvote

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// ValidationTest runs BeaconVote.Validate on each vote, expecting the same error code for all.
type ValidationTest struct {
	Name              string
	Type              string
	Documentation     string
	BeaconVotes       []*types.BeaconVote
	ExpectedErrorCode int
}

func (test *ValidationTest) TestName() string {
	return "beacon vote validation " + test.Name
}

func (test *ValidationTest) Run(t *testing.T) {
	for _, bv := range test.BeaconVotes {
		tests.AssertErrorCode(t, test.ExpectedErrorCode, bv.Validate())
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewValidationTest(name, documentation string, beaconVotes []*types.BeaconVote, expectedErrorCode int) *ValidationTest {
	return &ValidationTest{
		Name:              name,
		Type:              testdoc.BeaconVoteValidationTestType,
		Documentation:     documentation,
		BeaconVotes:       beaconVotes,
		ExpectedErrorCode: expectedErrorCode,
	}
}
