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

// MajoritySlashable tests a slashable attestation by majority of validators
func MajoritySlashable() tests.SpecTest {
	data := &types.BeaconVote{
		BlockRoot: testingutils.TestingBlockRoot,
		Source: &spec.Checkpoint{
			Epoch: 0,
			Root:  testingutils.TestingBlockRoot,
		},
		Target: &spec.Checkpoint{
			Epoch: 1,
			Root:  testingutils.TestingBlockRoot,
		},
	}

	input, _ := data.Encode()
	validatorIndexList := testingutils.ValidatorIndexList(4)
	// Get shares
	shareMap := testingutils.TestingValidatorShareMap(validatorIndexList)

	// Make slashable map with majority
	slashableMap := make(map[string][]phase0.Slot)
	for _, share := range shareMap {
		slashableMap[hex.EncodeToString(share.SharePubKey)] = []phase0.Slot{testingutils.TestingDutySlot}
	}

	return &valcheck.MultiSpecTest{
		Name: "attestation value check with slashable majority",
		Tests: []*valcheck.SpecTest{
			{
				Name:             "attestation duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingCommitteeAttesterDuty(testingutils.TestingDutySlot, validatorIndexList),
				Input:            input,
				ValidatorsShares: shareMap,
				SlashableSlots:   slashableMap,
				ExpectedError:    "slashable attestation",
			},
			{
				Name:       "sync committee duty",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Duty: testingutils.TestingCommitteeSyncCommitteeDuty(testingutils.TestingDutySlot,
					validatorIndexList),
				Input:            input,
				ValidatorsShares: shareMap,
				SlashableSlots:   slashableMap,
			},
			{
				Name:       "attestation and sync committee duty",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleCommittee,
				Duty: testingutils.TestingCommitteeDuty(testingutils.
					TestingDutySlot, validatorIndexList, validatorIndexList),
				Input:            input,
				ValidatorsShares: shareMap,
				SlashableSlots:   slashableMap,
				ExpectedError:    "slashable attestation",
			},
		},
	}
}
