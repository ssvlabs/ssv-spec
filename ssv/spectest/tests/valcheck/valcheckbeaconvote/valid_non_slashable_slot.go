package valcheckbeaconvote

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
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
	validatorIndices := testingutils.ValidatorIndexList(1)
	shareMap := testingutils.TestingValidatorShareMap(validatorIndices)
	sharePKBytes := shareMap[1].SharePubKey
	shareString := hex.EncodeToString(sharePKBytes)

	return &valcheck.MultiSpecTest{
		Name: "beacon vote valid with non slashable slot",
		Tests: []*valcheck.SpecTest{
			{
				Name:       "attestation duty",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Duty:       testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot+1, validatorIndices),
				Input:      input,
				SlashableSlots: map[string][]phase0.Slot{
					shareString: {
						testingutils.TestingDutySlot,
					},
				},
				ValidatorsShares: shareMap,
			},
			{
				Name:       "sync committee duty",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Duty:       testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot+1, validatorIndices),
				Input:      input,
				SlashableSlots: map[string][]phase0.Slot{
					shareString: {
						testingutils.TestingDutySlot,
					},
				},
				ValidatorsShares: shareMap,
			},
			{
				Name:       "attestation and sync committee duty",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Duty:       testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot+1, validatorIndices, validatorIndices),
				Input:      input,
				SlashableSlots: map[string][]phase0.Slot{
					shareString: {
						testingutils.TestingDutySlot,
					},
				},
				ValidatorsShares: shareMap,
			},
		},
	}
}
