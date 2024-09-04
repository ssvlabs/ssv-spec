package messageratetest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/p2p/messagerate"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

const (
	allowedError = 1e-3
)

type CommitteesWithValidators struct {
	NumCommittees int
	NumValidators int
}

type TestCase struct {
	// Input: list of CommitteesWithValidators(NumCommittees, NumValidators) objects.
	// Each element can be understood as: for this test case, there are [NumCommittees] committees each with a total of [NumValidators] validators.
	// Note: we use this compression, instead of naturally asking for a list of messagerate.Committees,
	// due to the big size of the JSON tests.
	CommitteesConfig []CommitteesWithValidators

	// Output: expected message rate for each test case
	ExpectedMessageRate float64
}

// SpecTest to test the message rate estimation
type MessageRateTest struct {
	Name string

	// List of test cases
	TestCases []TestCase
}

func (test *MessageRateTest) TestName() string {
	return "messagerate " + test.Name
}

func (test *MessageRateTest) Run(t *testing.T) {
	for _, testCase := range test.TestCases {

		// Build the list of committees
		committees := make([]*messagerate.Committee, 0)
		for _, committeesWithValidators := range testCase.CommitteesConfig {
			committees = append(committees, testingutils.TestingDisjointCommittees(committeesWithValidators.NumCommittees, committeesWithValidators.NumValidators)...)
		}

		// Check result
		expectedValue := testCase.ExpectedMessageRate
		result := messagerate.EstimateMessageRateForTopic(committees)
		require.InDelta(t, expectedValue, result, allowedError)
	}
}

func (test *MessageRateTest) GetPostState() (interface{}, error) {
	return nil, nil
}
