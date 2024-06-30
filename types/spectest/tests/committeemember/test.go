package committeemember

import (
	reflect2 "reflect"
	"testing"

	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type CommitteeMemberTest struct {
	Name                  string
	CommitteeMember       types.CommitteeMember
	Message               types.SignedSSVMessage
	ExpectedHasQuorum     bool
	ExpectedFullCommittee bool
	ExpectedError         string
}

func (test *CommitteeMemberTest) TestName() string {
	return "committee member " + test.Name
}

// Returns the number of unique signers in the message signers list
func (test *CommitteeMemberTest) GetUniqueMessageSignersCount() int {
	uniqueSigners := make(map[uint64]bool)

	for _, element := range test.Message.OperatorIDs {
		uniqueSigners[element] = true
	}

	return len(uniqueSigners)
}

func (test *CommitteeMemberTest) Run(t *testing.T) {

	// Validate message
	err := test.Message.Validate()
	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	// Get unique signers
	numSigners := test.GetUniqueMessageSignersCount()

	// Test expected thresholds results
	require.Equal(t, test.ExpectedHasQuorum, test.CommitteeMember.HasQuorum(numSigners))
	require.Equal(t, test.ExpectedFullCommittee, (len(test.CommitteeMember.Committee) == numSigners))

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}
