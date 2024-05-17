package valcheckattestations

import (
	"encoding/hex"
	"math"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MajoritySlashable tests a slashable attestation by majority of validators
func MajoritySlashable() tests.SpecTest {
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

	attestationData := &spec.AttestationData{
		Slot:            testingutils.TestingDutySlot,
		Index:           math.MaxUint64,
		BeaconBlockRoot: data.BlockRoot,
		Source:          data.Source,
		Target:          data.Target,
	}

	r, _ := attestationData.HashTreeRoot()

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
	slashableMap := make(map[string][][]byte)
	for i := 0; i < int(keySet.Threshold); i++ {
		slashableMap[sharesPKString[i]] = make([][]byte, 0)
		slashableMap[sharesPKString[i]] = append(slashableMap[sharesPKString[i]], r[:])
	}

	return &valcheck.SpecTest{
		Name:               "attestation value check with slashable majority",
		Network:            types.BeaconTestNetwork,
		RunnerRole:         types.RoleCommittee,
		Input:              input,
		ExpectedError:      "slashable attestation",
		SlashableDataRoots: slashableMap,
		ShareValidatorsPK:  sharesPKBytes,
	}
}
