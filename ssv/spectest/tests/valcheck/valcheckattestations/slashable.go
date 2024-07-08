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

	attestationData := &spec.AttestationData{
		Slot:            testingutils.TestingDutySlot,
		Index:           math.MaxUint64,
		BeaconBlockRoot: data.BlockRoot,
		Source:          data.Source,
		Target:          data.Target,
	}

	r, _ := attestationData.HashTreeRoot()

	input, _ := data.Encode()

	keySet := testingutils.Testing4SharesSet()
	sharePKBytes := keySet.Shares[1].Serialize()
	shareString := hex.EncodeToString(sharePKBytes)

	return &valcheck.SpecTest{
		Name:          "attestation value check slashable",
		Network:       types.BeaconTestNetwork,
		RunnerRole:    types.RoleCommittee,
		Input:         input,
		ExpectedError: "slashable attestation",
		SlashableDataRoots: map[string][][]byte{
			shareString: {
				r[:],
			},
		},
		ShareValidatorsPK: []types.ShareValidatorPK{sharePKBytes},
	}
}
