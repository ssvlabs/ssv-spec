package tests

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/dkg/keygen"
	"github.com/bloxapp/ssv-spec/dkg/testutils"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type WipMsgProcessingSpecTest struct {
	Name          string
	Pre           dkg.Protocol
	Messages      []*keygen.ParsedMessage
	Output        *keygen.LocalKeyShare
	KeySet        *testingutils.TestKeySet
	ExpectedError string
}

func (test *WipMsgProcessingSpecTest) TestName() string {
	return test.Name
}

func (test *WipMsgProcessingSpecTest) Run(t *testing.T) {
	pre := testutils.BaseInstance()


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

	lks := keygen.LocalKeyShare{}
	err = json.Unmarshal(output, &lks)

	if err != nil {
		lastErr = err
	}

	if len(test.ExpectedError) > 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
		require.Equal(t, lks, *test.Output)
	}
}
