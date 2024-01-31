package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidSyncCommitteeBlockValidation tests an invalid consensus data with sync committee block data.
func InvalidSyncCommitteeBlockValidation() *ConsensusDataTest {

	cd := types.ConsensusData{
		Duty:    testingutils.TestingSyncCommitteeDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: []byte{1},
	}

	return &ConsensusDataTest{
		Name:          "invalid sync committee",
		ConsensusData: cd,
		ExpectedError: "could not unmarshal ssz: expected buffer of length 32 receiced 1",
	}
}
