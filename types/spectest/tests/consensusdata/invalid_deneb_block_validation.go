package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidDenebBlockValidation tests an invalid consensus data with deneb block
func InvalidDenebBlockValidation() *ConsensusDataTest {
	version := spec.DataVersionDeneb
	cd := &types.ConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return &ConsensusDataTest{
		Name:          "invalid deneb block",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
