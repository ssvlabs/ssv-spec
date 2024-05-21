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

// MinoritySlashable tests a slashable attestation by majority of validators
func MinoritySlashable() tests.SpecTest {
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

	// Make slashable map with minority
	slashableMap := map[string][][]byte{
		sharesPKString[0]: {
			r[:],
		},
	}

	return &valcheck.SpecTest{
		Name:               "attestation value check with slashable minority",
		Network:            types.BeaconTestNetwork,
		RunnerRole:         types.RoleCommittee,
		Input:              input,
		ExpectedError:      "slashable attestation",
		SlashableDataRoots: slashableMap,
		ShareValidatorsPK:  sharesPKBytes,
	}
}
