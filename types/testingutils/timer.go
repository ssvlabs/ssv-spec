package testingutils

import "github.com/bloxapp/ssv-spec/qbft"

type TestQBFTTimer struct {
}

func NewTestingTimer() qbft.Timer {
	return &TestQBFTTimer{}
}

func (t *TestQBFTTimer) TimeoutForRound(round qbft.Round) {

}
