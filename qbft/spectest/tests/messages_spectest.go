package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgSpecTest struct {
	Name          string
	Messages      []*qbft.SignedMessage
	ExpectedError string
}

func (test *MsgSpecTest) Run(t *testing.T) {
	var lastErr error
	for _, msg := range test.Messages {
		if err := msg.Validate(); err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *MsgSpecTest) TestName() string {
	return test.Name
}
