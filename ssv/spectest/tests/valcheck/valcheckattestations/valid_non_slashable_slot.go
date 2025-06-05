package valcheckattestations

import (
	"encoding/hex"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ValidNonSlashableSlot tests a valid AttestationData with a slot that is not slashable
func ValidNonSlashableSlot() tests.SpecTest {
	data := &types.BeaconVote{
		BlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		Source: &spec.Checkpoint{
			Epoch: 0,
			Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
		Target: &spec.Checkpoint{
			Epoch: 1,
			Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
		},
	}

	input, _ := data.Encode()

	keySet := testingutils.Testing4SharesSet()
	sharePKBytes := keySet.Shares[1].Serialize()
	shareString := hex.EncodeToString(sharePKBytes)

	return &valcheck.SpecTest{
		Name:       "attestation valid with non slashable slot",
		Network:    types.BeaconTestNetwork,
		RunnerRole: types.RoleCommittee,
		DutySlot:   testingutils.TestingDutySlot + 1,
		Input:      input,
		SlashableSlots: map[string][]spec.Slot{
			shareString: {
				testingutils.TestingDutySlot,
			},
		},
		ShareValidatorsPK: []types.ShareValidatorPK{sharePKBytes},
	}
}
