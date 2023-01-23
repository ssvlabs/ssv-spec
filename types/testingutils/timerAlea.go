package testingutils


import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
)

type TimerStateAlea struct {
	Timeouts int
	Round    alea.Round
}

type TestQBFTTimerAlea struct {
	State TimerStateAlea
}

func NewTestingTimerAlea() alea.Timer {
	return &TestQBFTTimerAlea{
		State: TimerStateAlea{},
	}
}

func (t *TestQBFTTimerAlea) TimeoutForRound(round alea.Round) {
	t.State.Timeouts++
	t.State.Round = round
}
