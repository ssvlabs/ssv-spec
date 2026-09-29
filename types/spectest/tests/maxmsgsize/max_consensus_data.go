package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	MaxSizeDataSSZ = 8388608
)

func maxValidatorDuty() types.ValidatorDuty {

	validatorSyncCommitteeIndices := [13]uint64{1}

	return types.ValidatorDuty{
		Type:                          types.BNRoleAttester,
		PubKey:                        [48]byte{1},
		Slot:                          1,
		ValidatorIndex:                1,
		CommitteeIndex:                1,
		CommitteeLength:               1,
		CommitteesAtSlot:              2,
		ValidatorCommitteeIndex:       2,
		ValidatorSyncCommitteeIndices: validatorSyncCommitteeIndices[:],
	}
}

func maxDataSSZ() []byte {
	dataSSZ := [MaxSizeDataSSZ]byte{1}
	return dataSSZ[:]
}

func maxConsensusData() *types.ProposerConsensusData {

	return &types.ProposerConsensusData{
		Duty:    maxValidatorDuty(),
		Version: gloas.DataVersionGloas,
		DataSSZ: maxDataSSZ(),
	}
}

func MaxConsensusData() *StructureSizeTest {
	return NewStructureSizeTest(
		"max ProposerConsensusData",
		testdoc.StructureSizeTestMaxConsensusDataDoc,
		maxConsensusData(),
		MaxSizeFullConsensusData,
		true,
	)
}
