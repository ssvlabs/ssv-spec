package tests

import (
	"encoding/hex"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type MsgProcessingSpecTest struct {
	Name                    string
	Runner                  *ssv.Runner
	Duty                    *types.Duty
	Messages                []*types.SSVMessage
	PostDutyRunnerStateRoot string
	ExpectedError           string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return "msg processing " + test.Name
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	v := testingutils.BaseValidator(keySetForShare(test.Runner.Share))
	v.DutyRunners[test.Runner.BeaconRoleType] = test.Runner

	lastErr := v.StartDuty(test.Duty)
	for _, msg := range test.Messages {
		err := v.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	postRoot, err := test.Runner.State.GetRoot()
	require.NoError(t, err)

	require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(postRoot))
}

func keySetForShare(share *types.Share) *testingutils.TestKeySet {
	if share.Quorum == 5 {
		return testingutils.Testing7SharesSet()
	}
	if share.Quorum == 7 {
		return testingutils.Testing10SharesSet()
	}
	if share.Quorum == 9 {
		return testingutils.Testing13SharesSet()
	}
	return testingutils.Testing4SharesSet()
}
