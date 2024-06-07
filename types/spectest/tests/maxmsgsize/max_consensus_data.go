package maxmsgsize

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
)

const (
	maxSizeDataSSZ = 4194304
)

func maxBeaconDuty() types.BeaconDuty {

	validatorSyncCommitteeIndices := [13]uint64{1}

	return types.BeaconDuty{
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

func maxConsensusData() *types.ConsensusData {

	preConsensusJustification := make([]*types.PartialSignatureMessages, 0)
	for i := 0; i < 13; i++ {
		preConsensusJustification = append(preConsensusJustification, maxPartialSignatureMessagesForPreConsensus())
	}

	return &types.ConsensusData{
		Duty:                       maxBeaconDuty(),
		Version:                    spec.DataVersionAltair,
		PreConsensusJustifications: preConsensusJustification,
		DataSSZ:                    maxDataSSZ(),
	}
}

func MaxConsensusData() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "max ConsensusData",
		Object:                maxConsensusData(),
		ExpectedEncodedLength: maxSizeFullData,
		IsMaxSizeForType:      true,
	}
}
