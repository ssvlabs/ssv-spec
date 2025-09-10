package share

import (
	reflect2 "reflect"
	"testing"

	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/errcodes"
	"github.com/stretchr/testify/require"
)

type ShareTest struct {
	Name                  string
	Type                  string
	Documentation         string
	Share                 types.Share
	Message               types.SignedSSVMessage
	ExpectedHasQuorum     bool
	ExpectedFullCommittee bool
	ExpectedErrorCode     errcodes.Code
}

func (test *ShareTest) TestName() string {
	return "share " + test.Name
}

// Returns the number of unique signers in the message signers list
func (test *ShareTest) GetUniqueMessageSignersCount() int {
	uniqueSigners := make(map[uint64]bool)

	for _, element := range test.Message.OperatorIDs {
		uniqueSigners[element] = true
	}

	return len(uniqueSigners)
}

func (test *ShareTest) Run(t *testing.T) {

	// Validate message
	err := test.Message.Validate()
	if test.ExpectedErrorCode != 0 {
		require.Equal(t, test.ExpectedErrorCode, errcodes.FromError(err))
	} else {
		require.NoError(t, err)
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewShareTest(name, documentation string, share types.Share, message types.SignedSSVMessage, expectedHasQuorum bool, expectedFullCommittee bool, expectedErrorCode errcodes.Code) *ShareTest {
	return &ShareTest{
		Name:                  name,
		Type:                  testdoc.ShareTestType,
		Documentation:         documentation,
		Share:                 share,
		Message:               message,
		ExpectedHasQuorum:     expectedHasQuorum,
		ExpectedFullCommittee: expectedFullCommittee,
		ExpectedErrorCode:     expectedErrorCode,
	}
}
