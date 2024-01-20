package timeout

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	typescomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

// TimeoutDurationTest tests the expected duration of the timer
type TimeoutDurationTest struct {
	Name             string
	Role             types.BeaconRole
	Height           qbft.Height
	Round            qbft.Round
	Network          types.BeaconNetwork
	CurrentTime      uint64
	ExpectedDuration uint64
}

func (test *TimeoutDurationTest) GetPostState() (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (test *TimeoutDurationTest) TestName() string {
	return "qbft timeout duration " + test.Name
}

func (test *TimeoutDurationTest) Run(t *testing.T) {
	timer := qbft.RoundTimer{
		Role:        test.Role,
		Height:      test.Height,
		Network:     test.Network,
		CurrentTime: test.CurrentTime,
	}

	require.Equal(t, test.ExpectedDuration, timer.TimeoutForRound(test.Round), "timeout duration is not as expected")
}

type UponTimeoutTest struct {
	Name               string
	Pre                *qbft.Instance
	PostRoot           string
	PostState          types.Root `json:"-"` // Field is ignored by encoding/json
	OutputMessages     []*qbft.SignedMessage
	ExpectedTimerState *testingutils.TimerState
	ExpectedError      string
}

func (test *UponTimeoutTest) TestName() string {
	return "qbft upon timeout " + test.Name
}

func (test *UponTimeoutTest) Run(t *testing.T) {
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
	broadcastedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	if len(test.OutputMessages) > 0 || len(broadcastedMsgs) > 0 {
		require.Len(t, broadcastedMsgs, len(test.OutputMessages))

		for i, msg := range test.OutputMessages {
			r1, _ := msg.GetRoot()

			msg2 := &qbft.SignedMessage{}
			require.NoError(t, msg2.Decode(broadcastedMsgs[i].Data))
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

func (test *UponTimeoutTest) GetPostState() (interface{}, error) {
	return nil, nil
}
