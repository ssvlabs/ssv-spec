package valcheckattestations

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidNonSlashableSlot tests a valid AttestationData with a slot that is not slashable
func ValidNonSlashableSlot() tests.SpecTest {
	data := &types.BeaconVote{
		BlockRoot: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		Source: &phase0.Checkpoint{
			Epoch: 0,
			Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Target: &phase0.Checkpoint{
			Epoch: 1,
			Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
	}

	input, _ := data.Encode()

	keySet := testingutils.Testing4SharesSet()
	sharePKBytes := keySet.Shares[1].Serialize()
	shareString := hex.EncodeToString(sharePKBytes)

	return valcheck.NewSpecTest(
		"attestation valid with non slashable slot",
		testdoc.ValCheckAttestationValidNonSlashableSlotDoc,
		types.BeaconTestNetwork,
		types.RoleCommittee,
		testingutils.TestingDutySlot+1,
		input,
		0,
		1,
		map[string][]phase0.Slot{
			shareString: {
				testingutils.TestingDutySlot,
			},
		},
		[]types.ShareValidatorPK{sharePKBytes},
		0,
		false,
	)
}
