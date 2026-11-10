package valcheckduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ProposerBlockSlotMismatch tests a proposer value whose block slot differs from the duty slot
func ProposerBlockSlotMismatch() tests.SpecTest {
	return proposerBlockMismatch(
		"proposer block slot mismatch",
		testdoc.ValCheckDutyProposerBlockSlotMismatchDoc,
		func(slot *phase0.Slot, _ *phase0.ValidatorIndex) { *slot++ },
		types.ProposerBlockSlotMismatchErrorCode,
	)
}

// ProposerBlockProposerIndexMismatch tests a proposer value whose block proposer index differs from the duty
// validator index
func ProposerBlockProposerIndexMismatch() tests.SpecTest {
	return proposerBlockMismatch(
		"proposer block proposer index mismatch",
		testdoc.ValCheckDutyProposerBlockProposerIndexMismatchDoc,
		func(_ *phase0.Slot, index *phase0.ValidatorIndex) { *index++ },
		types.ProposerBlockProposerIndexMismatchErrorCode,
	)
}

func proposerBlockMismatch(
	name, doc string,
	mutate func(*phase0.Slot, *phase0.ValidatorIndex),
	expectedErrCode int,
) tests.SpecTest {
	specTests := make([]*valcheck.SpecTest, 0)
	for _, version := range testingutils.SupportedBlockVersions {
		for _, blinded := range []bool{false, true} {
			cd := testingutils.TestProposerConsensusDataV(version)
			if blinded {
				cd = testingutils.TestProposerBlindedBlockConsensusDataV(version)
			}
			specTests = append(specTests, &valcheck.SpecTest{
				Name:              fmt.Sprintf("%s blinded=%t", version.String(), blinded),
				Network:           types.BeaconTestNetwork,
				RunnerRole:        types.RoleProposer,
				Input:             mutateBlock(cd, mutate),
				ExpectedErrorCode: expectedErrCode,
			})
		}
	}
	return valcheck.NewMultiSpecTest(name, doc, specTests)
}

// mutateBlock returns the encoded consensus data with the block's slot and proposer index changed by mutate
func mutateBlock(cd *types.ProposerConsensusData, mutate func(*phase0.Slot, *phase0.ValidatorIndex)) []byte {
	blk, _, err := cd.GetBlockData() // decodes a fresh copy
	if err != nil {
		panic(err.Error())
	}

	var obj ssz.Marshaler
	switch {
	case blk.Version == spec.DataVersionCapella && blk.Blinded:
		mutate(&blk.CapellaBlinded.Slot, &blk.CapellaBlinded.ProposerIndex)
		obj = blk.CapellaBlinded
	case blk.Version == spec.DataVersionCapella:
		mutate(&blk.Capella.Slot, &blk.Capella.ProposerIndex)
		obj = blk.Capella
	case blk.Version == spec.DataVersionDeneb && blk.Blinded:
		mutate(&blk.DenebBlinded.Slot, &blk.DenebBlinded.ProposerIndex)
		obj = blk.DenebBlinded
	case blk.Version == spec.DataVersionDeneb:
		mutate(&blk.Deneb.Block.Slot, &blk.Deneb.Block.ProposerIndex)
		obj = blk.Deneb
	case blk.Version == spec.DataVersionElectra && blk.Blinded:
		mutate(&blk.ElectraBlinded.Slot, &blk.ElectraBlinded.ProposerIndex)
		obj = blk.ElectraBlinded
	case blk.Version == spec.DataVersionElectra:
		mutate(&blk.Electra.Block.Slot, &blk.Electra.Block.ProposerIndex)
		obj = blk.Electra
	case blk.Version == spec.DataVersionFulu && blk.Blinded:
		mutate(&blk.FuluBlinded.Slot, &blk.FuluBlinded.ProposerIndex)
		obj = blk.FuluBlinded
	case blk.Version == spec.DataVersionFulu:
		mutate(&blk.Fulu.Block.Slot, &blk.Fulu.Block.ProposerIndex)
		obj = blk.Fulu
	default:
		panic("unsupported version")
	}

	dataSSZ, err := obj.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	ret, err := (&types.ProposerConsensusData{Duty: cd.Duty, Version: cd.Version, DataSSZ: dataSSZ}).Encode()
	if err != nil {
		panic(err.Error())
	}
	return ret
}
