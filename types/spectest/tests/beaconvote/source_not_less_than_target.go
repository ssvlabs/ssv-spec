package beaconvote

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

// SourceNotLessThanTarget tests beacon votes whose source epoch is greater than or equal to the target epoch
func SourceNotLessThanTarget() *ValidationTest {
	return NewValidationTest(
		"source not less than target",
		testdoc.BeaconVoteValidationSourceNotLessThanTargetDoc,
		[]*types.BeaconVote{
			{Source: &phase0.Checkpoint{Epoch: 2}, Target: &phase0.Checkpoint{Epoch: 1}},
			{Source: &phase0.Checkpoint{Epoch: 1}, Target: &phase0.Checkpoint{Epoch: 1}},
		},
		types.AttestationSourceNotLessThanTargetErrorCode,
	)
}
