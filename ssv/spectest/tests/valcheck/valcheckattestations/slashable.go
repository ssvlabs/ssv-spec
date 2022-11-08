package valcheckattestations

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// Slashable tests a slashable AttestationData
func Slashable() *valcheck.SpecTest {
	cd := &types.ConsensusData{
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
		AttestationData: &phase0.AttestationData{
			Slot:            testingutils.TestingDutySlot,
			Index:           50,
			BeaconBlockRoot: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			Source: &phase0.Checkpoint{
				Epoch: 0,
				Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			},
			Target: &phase0.Checkpoint{
				Epoch: 1,
				Root:  phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
			},
		},
	}

	r, _ := cd.AttestationData.HashTreeRoot()

	source, _ := cd.MarshalSSZ()
	root, _ := cd.HashTreeRoot()
	input := &qbft.Data{
		Root:   root,
		Source: source,
	}

	return &valcheck.SpecTest{
		Name:       "attestation value check slashable",
		Network:    types.NowTestNetwork,
		BeaconRole: types.BNRoleAttester,
		Input:      input,
		AnyError:   true,
		SlashableDataRoots: [][]byte{
			r[:],
		},
	}
}
