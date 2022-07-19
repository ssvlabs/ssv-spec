package tests

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/dkg/types"
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgProcessingSpecTest struct {
	Name     string
	Pre      types.Protocol
	Messages []*keygen.ParsedMessage
	Output   *types.LocalKeyShare
	KeySet   *testingutils.TestKeySet
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

	lks := types.LocalKeyShare{}
	err = json.Unmarshal(output, &lks)

	if err != nil {
		lastErr = err
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
		require.Equal(t, *test.Output, lks)
	}
}
