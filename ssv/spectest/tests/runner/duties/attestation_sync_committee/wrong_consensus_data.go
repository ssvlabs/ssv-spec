package newduty

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
)

// WrongConsensusData sends a proposal message to a cluster runner with an invalid consensus data that can't be decoded to BeaconVote
// Expects: error
func WrongConsensusData() tests.SpecTest {
	panic("implement me")
}
