package valcheckproposer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// GloasBlocks covers the Gloas (ePBS, SIP #94 §4) proposer value-check rules. At a Gloas duty slot the
// consensus value carries an opaque bid-only gloas.BeaconBlock in DataSSZ, so the value check decodes
// it directly (never routing through the pre-Gloas Validate()/GetBlockData()), pins the block's slot to
// the duty's, and rejects a leader-stamped Version that disagrees with the duty slot's fork in either
// direction.
func GloasBlocks() tests.SpecTest {
	gloasDuty := testingutils.TestingProposerDutyV(gloas.DataVersionGloas)
	electraDuty := testingutils.TestingProposerDutyV(spec.DataVersionElectra)

	encode := func(cd *types.ProposerConsensusData) []byte {
		byts, err := cd.Encode()
		if err != nil {
			panic(err.Error())
		}
		return byts
	}
	gloasBlockBytes := testingutils.TestingBeaconBlockBytesV(gloas.DataVersionGloas)

	slotMismatchBlock, err := gloas.TestingBeaconBlock(gloasDuty.Slot + 100).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	electraSlotBlock, err := gloas.TestingBeaconBlock(electraDuty.Slot).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return valcheck.NewMultiSpecTest(
		"gloas blocks",
		testdoc.ValCheckProposerGloasBlocksDoc,
		[]*valcheck.SpecTest{
			{
				Name:       "valid gloas block",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleProposer,
				Input:      testingutils.TestProposerConsensusDataBytsV(gloas.DataVersionGloas),
			},
			{
				Name:              "undecodable gloas block",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *gloasDuty, Version: gloas.DataVersionGloas, DataSSZ: []byte("garbage")}),
				ExpectedErrorCode: types.UnmarshalSSZErrorCode,
			},
			{
				Name:              "block slot does not match duty slot",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *gloasDuty, Version: gloas.DataVersionGloas, DataSSZ: slotMismatchBlock}),
				ExpectedErrorCode: types.ProposerBlockSlotMismatchErrorCode,
			},
			{
				// The leader-stamped Version is attacker-controlled, so on a Gloas slot it is pinned to
				// the slot's fork; without this rule a mixed cluster could split on the same value.
				Name:              "pre-gloas version on a gloas slot",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *gloasDuty, Version: spec.DataVersionElectra, DataSSZ: gloasBlockBytes}),
				ExpectedErrorCode: types.QBFTValueInvalidErrorCode,
			},
			{
				// The reverse mismatch takes the pre-Gloas branch, where Validate() rejects the Gloas
				// version as unknown — both mismatch directions fail.
				Name:              "gloas version on a pre-gloas slot",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *electraDuty, Version: gloas.DataVersionGloas, DataSSZ: electraSlotBlock}),
				ExpectedErrorCode: types.QBFTValueInvalidErrorCode,
			},
			{
				Name:       "slashable proposal slot",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleProposer,
				Input:      testingutils.TestProposerConsensusDataBytsV(gloas.DataVersionGloas),
				SlashableSlots: map[string][]phase0.Slot{
					hex.EncodeToString(testingutils.Testing4SharesSet().Shares[1].Serialize()): {
						gloasDuty.Slot,
					},
				},
				ExpectedErrorCode: types.SlashableProposalErrorCode,
			},
		},
	)
}
