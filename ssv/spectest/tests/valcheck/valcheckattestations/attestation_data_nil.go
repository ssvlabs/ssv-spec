package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BeaconVoteDataNil tests value-check rejection of a BeaconVote with nil Source/Target (which encode to
// zero checkpoints, so source is not less than target).
func BeaconVoteDataNil() tests.SpecTest {
	consensusData := &types.BeaconVote{
		Source: nil,
		Target: nil,
	}
	input, _ := consensusData.Encode()

	return valcheck.NewSpecTest(
		"consensus data value check nil",
		testdoc.ValCheckAttestationBeaconVoteDataNilDoc,
		types.PraterNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlot,
		input,
		phase0.Checkpoint{},
		phase0.Checkpoint{},
		nil,
		nil,
		types.AttestationSourceNotLessThanTargetErrorCode,
		false,
	)
}
