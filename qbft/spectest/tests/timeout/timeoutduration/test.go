package timeoutduration

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

// TimeoutDurationTest tests the expected duration of the timer
type TimeoutDurationTest struct {
	Name             string
	Role             types.BeaconRole
	Height           qbft.Height
	Round            qbft.Round
	Network          types.BeaconNetwork
	CurrentTime      int64
	ExpectedDuration int64
}

func (test *TimeoutDurationTest) GetPostState() (interface{}, error) {
	return nil, nil
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
