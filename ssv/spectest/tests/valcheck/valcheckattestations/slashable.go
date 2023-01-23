package valcheckattestations

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// Slashable tests a slashable AttestationData
func Slashable() *valcheck.SpecTest {
	data := &types.ConsensusData{
		Duty: &types.Duty{
			Type:                    types.BNRoleAttester,
			PubKey:                  testingutils.TestingValidatorPubKey,
			Slot:                    testingutils.TestingDutySlot,
			ValidatorIndex:          testingutils.TestingValidatorIndex,
			CommitteeIndex:          50,
			CommitteesAtSlot:        36,
			CommitteeLength:         128,
			ValidatorCommitteeIndex: 11,
		},
		AttestationData: &spec.AttestationData{
			Slot:            testingutils.TestingDutySlot,
			Index:           50,
			BeaconBlockRoot: spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			Source: &spec.Checkpoint{
				Epoch: 0,
				Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			},
			Target: &spec.Checkpoint{
				Epoch: 1,
				Root:  spec.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			},
		},
	}

	r, _ := data.AttestationData.HashTreeRoot()

	input, _ := data.Encode()

	return &valcheck.SpecTest{
		Name:       "attestation value check slashable",
		Network:    types.BeaconTestNetwork,
		BeaconRole: types.BNRoleAttester,
		Input:      input,
		AnyError:   true,
		SlashableDataRoots: [][]byte{
			r[:],
		},
	}
}
