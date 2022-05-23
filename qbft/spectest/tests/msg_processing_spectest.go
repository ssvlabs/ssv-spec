package tests

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgProcessingSpecTest struct {
	Name          string
	Pre           *qbft.Instance
	PostRoot      string
	Messages      []*qbft.SignedMessage
	ExpectedError string
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	var lastErr error
	for _, msg := range test.Messages {
		_, _, _, err := test.Pre.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)

	require.EqualValues(t, test.PostRoot, hex.EncodeToString(postRoot), "post root not valid")
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}
