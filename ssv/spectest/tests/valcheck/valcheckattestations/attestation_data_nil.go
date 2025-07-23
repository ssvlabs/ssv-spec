package valcheckattestations

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BeaconVoteDataNil tests consensus data != nil
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
		nil,
		nil,
		"attestation data source >= target",
		false,
	)
}
