package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// InvalidIndex tests that an attestation data index above 1 is rejected (SIP #94 §2: only 0 = payload absent
// and 1 = payload present are valid).
func InvalidIndex() tests.SpecTest {
	data := testingutils.TestBeaconVote
	data.AttestationDataIndex = 2
	input, err := data.Encode()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewSpecTest(
		"attestation value check invalid index",
		testdoc.ValCheckAttestationInvalidIndexDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlot,
		input,
		*data.Source,
		*data.Target,
		map[string][]phase0.Slot{},
		[]types.ShareValidatorPK{},
		types.BeaconVoteInvalidIndexErrorCode,
		false,
	)
}
