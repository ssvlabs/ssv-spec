package timeout

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

type SpecTest struct {
	Name               string
	Pre                *qbft.Instance
	PostRoot           string
	PostState          types.Root `json:"-"` // Field is ignored by encoding/json
	OutputMessages     []*types.SignedSSVMessage
	ExpectedTimerState *testingutils.TimerState
	ExpectedError      string
}

func (test *SpecTest) TestName() string {
	return "qbft timeout " + test.Name
}

func (test *SpecTest) Run(t *testing.T) {
	err := test.Pre.UponRoundTimeout()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	// test calling timeout
	timer, ok := test.Pre.GetConfig().GetTimer().(*testingutils.TestQBFTTimer)
	require.True(t, ok)
	require.Equal(t, test.ExpectedTimerState.Timeouts, timer.State.Timeouts)
	require.Equal(t, test.ExpectedTimerState.Round, timer.State.Round)

	// test output message
	broadcastedSignedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	require.NoError(t, testingutils.VerifyListOfSignedSSVMessages(broadcastedSignedMsgs, test.Pre.State.CommitteeMember.Committee))
	if len(test.OutputMessages) > 0 || len(broadcastedSignedMsgs) > 0 {
		require.Len(t, broadcastedSignedMsgs, len(test.OutputMessages))

		for i, msg := range test.OutputMessages {
			r1, _ := msg.GetRoot()

			msg2 := broadcastedSignedMsgs[i]
			r2, _ := msg2.GetRoot()

			require.EqualValuesf(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
		}
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)
	if test.PostRoot != hex.EncodeToString(postRoot[:]) {
		diff := typescomparable.PrintDiff(test.Pre.State, test.PostState)
		require.Fail(t, "post state not equal", diff)
	}
}

func (test *SpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}
