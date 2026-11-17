package beaconvote

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// NilCheckpoint tests a nil beacon vote and beacon votes with a nil source and/or target checkpoint
func NilCheckpoint() *ValidationTest {
	return NewValidationTest(
		"nil checkpoint",
		testdoc.BeaconVoteValidationNilCheckpointDoc,
		[]*types.BeaconVote{
			nil,
			{},
			{Target: &phase0.Checkpoint{Epoch: 1}},
			{Source: &phase0.Checkpoint{Epoch: 0}},
		},
		types.BeaconVoteNilCheckpointErrorCode,
	)
}
