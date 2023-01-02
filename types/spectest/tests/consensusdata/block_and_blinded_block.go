package consensusdata

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// BlockAndBlindedBlock tests a non nil block and blinded
func BlockAndBlindedBlock() *ValidationSpecTest {
	return &ValidationSpecTest{
		Name: "block and blinded block",
		Obj: &types.ConsensusData{
			Duty:             testingutils.TestingProposerDuty,
			BlockData:        testingutils.TestingBeaconBlock,
			BlindedBlockData: testingutils.TestingBlindedBeaconBlock,
		},
		ExpectedErr: "block and blinded block data are both != nil",
	}
}
