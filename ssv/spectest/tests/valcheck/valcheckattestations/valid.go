package valcheckattestations

import (
	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Valid tests valid data
func Valid() tests.SpecTest {
	return valcheck.NewSpecTest(
		"attestation value check valid",
		testdoc.ValCheckAttestationValidDoc,
		types.PraterNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlot,
		testingutils.TestBeaconVoteByts,
		nil,
		nil,
		"",
		false,
	)
}
