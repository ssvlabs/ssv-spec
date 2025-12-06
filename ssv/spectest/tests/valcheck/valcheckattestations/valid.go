package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
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
		*testingutils.TestBeaconVote.Source,
		*testingutils.TestBeaconVote.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		0,
		false,
	)
}