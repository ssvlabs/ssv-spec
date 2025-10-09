package share

import (
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

type ShareTest struct {
	Name                  string
	Type                  string
	Documentation         string
	Share                 types.Share
	Message               types.SignedSSVMessage
	ExpectedHasQuorum     bool
	ExpectedFullCommittee bool
	ExpectedErrorCode     int
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
	tests.AssertErrorCode(t, test.ExpectedErrorCode, err)

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}

func NewShareTest(name, documentation string, share types.Share, message types.SignedSSVMessage, expectedHasQuorum bool, expectedFullCommittee bool, expectedErrorCode int) *ShareTest {
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
