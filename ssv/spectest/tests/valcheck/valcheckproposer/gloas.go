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
// consensus value carries an opaque {block, payload_root} GloasProposalData in DataSSZ, so the value check
// decodes it directly (never routing through the pre-Gloas Validate()/GetBlockData()), pins the block's
// slot to the duty's and its proposer index to the duty's validator, pins the duty's slot to the running
// duty's, requires payload_root to be zero iff the bid is not self-build, and rejects a leader-stamped
// Version that disagrees with the duty slot's fork in either direction.
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

	slotMismatchBlock := testingutils.TestingGloasProposalDataBytes(gloasDuty.Slot + 100)
	// A self-build value (the fixture bid is self-build) with a zero payload_root trips the §4 presence rule.
	selfBuildZeroPayloadRoot, err := (&gloas.GloasProposalData{Block: gloas.TestingBeaconBlock(gloasDuty.Slot)}).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	electraSlotBlock, err := gloas.TestingBeaconBlock(electraDuty.Slot).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	// Correct slot but a proposer index that is not the duty's validator — trips the §4 proposer-index pin.
	wrongProposerBlock := gloas.TestingBeaconBlock(gloasDuty.Slot)
	wrongProposerBlock.ProposerIndex = gloasDuty.ValidatorIndex + 1
	wrongProposerIndexBytes, err := (&gloas.GloasProposalData{Block: wrongProposerBlock, PayloadRoot: testingutils.TestingGloasPayloadRoot}).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	// An external-bid value (not self-build) must carry a zero payload_root; a non-zero one trips the §4 rule.
	externalNonZeroPayload, err := (&gloas.GloasProposalData{Block: gloas.TestingBeaconBlockExternalBuild(gloasDuty.Slot), PayloadRoot: testingutils.TestingGloasPayloadRoot}).MarshalSSZ()
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
				// SIP #94 §4: the block's proposer index must equal the duty's validator (whose key the
				// cluster signs the block root with); a mismatch would have the cluster sign a block the
				// beacon node rejects, losing the slot.
				Name:              "block proposer index does not match duty validator index",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *gloasDuty, Version: gloas.DataVersionGloas, DataSSZ: wrongProposerIndexBytes}),
				ExpectedErrorCode: types.ProposerBlockProposerIndexMismatchErrorCode,
			},
			{
				// SIP #94 §4: a value for a slot other than the running duty's is rejected before consensus
				// can commit it — else an operator would decide the instance on a slot it is not proposing
				// and brick the duty. DutySlot drives the running-slot provider; here it is one past the
				// value's own slot. (Isolated cases above leave DutySlot 0, which skips this bind.)
				Name:              "duty slot does not match running slot",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				DutySlot:          gloasDuty.Slot + 1,
				Input:             testingutils.TestProposerConsensusDataBytsV(gloas.DataVersionGloas),
				ExpectedErrorCode: types.ProposerDutySlotMismatchErrorCode,
			},
			{
				// payload_root MUST be non-zero on the self-build path (SIP #94 §4); zero trips the rule.
				Name:              "self-build with zero payload_root",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *gloasDuty, Version: gloas.DataVersionGloas, DataSSZ: selfBuildZeroPayloadRoot}),
				ExpectedErrorCode: types.QBFTValueInvalidErrorCode,
			},
			{
				// The mirror image: an external-bid value must carry a zero payload_root (SIP #94 §4); a
				// non-zero root trips the same presence rule.
				Name:              "external bid with non-zero payload_root",
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             encode(&types.ProposerConsensusData{Duty: *gloasDuty, Version: gloas.DataVersionGloas, DataSSZ: externalNonZeroPayload}),
				ExpectedErrorCode: types.QBFTValueInvalidErrorCode,
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
