package valcheckattestations

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// Slashable tests a slashable AttestationData
func Slashable() tests.SpecTest {
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

	return &valcheck.SpecTest{
		Name:                "attestation value check slashable",
		Network:             types.BeaconTestNetwork,
		RunnerRole:          types.RoleCommittee,
		DutySlot:            testingutils.TestingDutySlot,
		Input:               input,
		ExpectedSourceEpoch: 0,
		ExpectedTargetEpoch: 1,
		ExpectedError:       "slashable attestation",
		SlashableSlots: map[string][]phase0.Slot{
			shareString: {
				testingutils.TestingDutySlot,
			},
		},
		ShareValidatorsPK: []types.ShareValidatorPK{sharePKBytes},
	}
}
