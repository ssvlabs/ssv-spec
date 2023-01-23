package valcheckattestations

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
)

// SlotMismatch tests Duty.Slot != AttestationData.Slot
func SlotMismatch() *valcheck.SpecTest {
	data := &types.ConsensusData{
		Duty: &types.Duty{
			Type:                    types.BNRoleAttester,
			PubKey:                  testingutils.TestingValidatorPubKey,
			Slot:                    1,
			ValidatorIndex:          testingutils.TestingValidatorIndex,
			CommitteeIndex:          3,
			CommitteesAtSlot:        36,
			CommitteeLength:         128,
			ValidatorCommitteeIndex: 11,
		},
		AttestationData: &spec.AttestationData{
			Slot:            2,
			Index:           3,
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

	input, _ := data.Encode()

	return &valcheck.SpecTest{
		Name:          "attestation value check slot mismatch",
		Network:       types.PraterNetwork,
		BeaconRole:    types.BNRoleAttester,
		Input:         input,
		ExpectedError: "attestation data slot != duty slot",
	}
}
