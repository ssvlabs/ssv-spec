package tests

import (
	"github.com/bloxapp/ssv-spec/dkg/types"
	types2 "github.com/bloxapp/ssv-spec/gg20/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgProcessingSpecTest struct {
	Name          string
	Pre           types.Protocol
	Messages      []*types2.ParsedMessage
	Output        *types.LocalKeyShare
	KeySet        *testingutils.TestKeySet
	ExpectedError string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	pre := test.Pre

	var lastErr, err error
	_, err = pre.Start()
	if err != nil {
		lastErr = err
	}
	for _, msg := range test.Messages {

		if baseMsg, err := msg.ToBase(); err == nil {
			_, err = pre.ProcessMsg(baseMsg)
		}

		if err != nil {
			lastErr = err
		}
	}

	output, err := pre.Output()
	if err != nil {
		lastErr = err
	}

	expected, err := test.Output.Encode()
	require.NoError(t, err)

	if err != nil {
		lastErr = err
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
		require.Equal(t, expected, output)
	}
}
