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

// Slashable tests a slashable AttestationData
func Slashable() tests.SpecTest {
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
	shareMap := testingutils.TestingValidatorShareMap(testingutils.ValidatorIndexList(1))
	sharePKBytes := shareMap[1].SharePubKey
	shareString := hex.EncodeToString(sharePKBytes)

	return &valcheck.MultiSpecTest{
		Name: "attestation value check slashable",
		Tests: []*valcheck.SpecTest{
			{
				Name:             "attestation duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingAttesterDuty,
				Input:            input,
				ExpectedError:    "slashable attestation",
				SlashableSlots:   map[string][]phase0.Slot{shareString: {testingutils.TestingDutySlot}},
				ValidatorsShares: shareMap,
			},
			{
				Name:             "sync committee duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingSyncCommitteeDuty,
				Input:            input,
				SlashableSlots:   map[string][]phase0.Slot{shareString: {testingutils.TestingDutySlot}},
				ValidatorsShares: shareMap,
			},
			{
				Name:             "attestation and sync committee duty",
				Network:          types.BeaconTestNetwork,
				RunnerRole:       types.RoleCommittee,
				Duty:             testingutils.TestingCommitteeDuty(testingutils.TestingDutySlot, testingutils.ValidatorIndexList(1), testingutils.ValidatorIndexList(1)),
				Input:            input,
				ExpectedError:    "slashable attestation",
				SlashableSlots:   map[string][]phase0.Slot{shareString: {testingutils.TestingDutySlot}},
				ValidatorsShares: shareMap,
			},
		},
	}
}
