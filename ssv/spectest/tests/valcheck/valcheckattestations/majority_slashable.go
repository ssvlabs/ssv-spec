package valcheckattestations

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

	// Get shares
	keySet := testingutils.Testing4SharesSet()
	sharesPKBytes := make([]types.ShareValidatorPK, 0)
	sharesPKString := make([]string, 0)
	for _, shareKey := range keySet.Shares {
		shareBytes := shareKey.Serialize()
		sharesPKBytes = append(sharesPKBytes, shareBytes)
		sharesPKString = append(sharesPKString, hex.EncodeToString(shareBytes))
	}

	// Make slashable map with majority
	slashableMap := make(map[string][]phase0.Slot)
	for i := uint64(0); i < keySet.Threshold; i++ {
		slashableMap[sharesPKString[i]] = []phase0.Slot{testingutils.TestingDutySlot}
	}

	return &valcheck.SpecTest{
		Name:              "attestation value check with slashable majority",
		Network:           types.BeaconTestNetwork,
		RunnerRole:        types.RoleCommittee,
		DutySlot:          testingutils.TestingDutySlot,
		Input:             input,
		ExpectedError:     "slashable attestation",
		SlashableSlots:    slashableMap,
		ShareValidatorsPK: sharesPKBytes,
	}
}
