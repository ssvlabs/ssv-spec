package maxmsgsize

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizeDataSSZ = 4194304
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
	dataSSZ := [maxSizeDataSSZ]byte{1}
	return dataSSZ[:]
}

func maxConsensusData() *types.ValidatorConsensusData {

	return &types.ValidatorConsensusData{
		Duty:    maxValidatorDuty(),
		Version: spec.DataVersionAltair,
		DataSSZ: maxDataSSZ(),
	}
}

func MaxConsensusData() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "max ValidatorConsensusData",
		Object:                maxConsensusData(),
		ExpectedEncodedLength: maxSizeFullData,
		IsMaxSize:             true,
	}
}
