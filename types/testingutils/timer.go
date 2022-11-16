package testingutils

import "github.com/bloxapp/ssv-spec/qbft"

type TestQBFTTimer struct {
	Timeouts int
	Round qbft.Round
}

func NewTestingTimer() qbft.Timer {
	return &TestQBFTTimer{}
}

func (t *TestQBFTTimer) TimeoutForRound(round qbft.Round) {
	t.Timeouts++
	t.Round = round
}
