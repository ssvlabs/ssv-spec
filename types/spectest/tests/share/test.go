package share

import (
	"testing"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
)

type ShareTest struct {
	Name                     string
	Share                    types.Share
	Message                  qbft.SignedMessage
	ExpectedHasPartialQuorum bool
	ExpectedHasQuorum        bool
	ExpectedFullCommittee    bool
	ExpectedError            string
}

func (test *ShareTest) TestName() string {
	return "share " + test.Name
}

// Returns the number of unique signers in the message signers list
func (test *ShareTest) GetUniqueMessageSignersCount() int {
	uniqueSigners := make(map[uint64]bool)

	for _, element := range test.Message.Signers {
		uniqueSigners[element] = true
	}

	return len(uniqueSigners)
}

func (test *ShareTest) Run(t *testing.T) {

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
	require.Equal(t, test.ExpectedHasPartialQuorum, test.Share.HasPartialQuorum(numSigners))
	require.Equal(t, test.ExpectedHasQuorum, test.Share.HasQuorum(numSigners))
	require.Equal(t, test.ExpectedFullCommittee, (len(test.Share.Committee) == numSigners))
}
