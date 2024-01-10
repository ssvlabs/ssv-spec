package signedssvmsg

import (
	"testing"

	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
)

type SignedSSVMessageTest struct {
	Name          string
	Messages      []*types.SignedSSVMessage
	ExpectedError string
}

func (test *SignedSSVMessageTest) TestName() string {
	return "signedssvmsg " + test.Name
}

func (test *SignedSSVMessageTest) Run(t *testing.T) {

	for _, msg := range test.Messages {

		// test validation
		err := msg.Validate()

		// decode Data if there is no error
		if err == nil {
			_, err = msg.GetSSVMessageFromData()
		}

		if len(test.ExpectedError) != 0 {
			require.EqualError(t, err, test.ExpectedError)
		} else {
			require.NoError(t, err)
		}
	}
}
