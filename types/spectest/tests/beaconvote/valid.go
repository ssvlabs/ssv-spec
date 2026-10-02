package beaconvote

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests a valid beacon vote
func Valid() *ValidationTest {
	bv := testingutils.TestBeaconVote

	return NewValidationTest(
		"valid",
		testdoc.BeaconVoteValidationValidDoc,
		[]*types.BeaconVote{&bv},
		0,
	)
}
