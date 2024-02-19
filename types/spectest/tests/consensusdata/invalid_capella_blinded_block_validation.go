package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidCapellaBlindedBlockValidation tests an invalid consensus data with capella blinded block
func InvalidCapellaBlindedBlockValidation() *ConsensusDataTest {

	version := spec.DataVersionCapella

	cd := &types.ConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return &ConsensusDataTest{
		Name:          "invalid capella blinded block",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
