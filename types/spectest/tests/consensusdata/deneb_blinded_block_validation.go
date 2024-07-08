package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlindedBlockValidation tests a valid consensus data with deneb blinded block
func DenebBlindedBlockValidation() *ConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ConsensusDataTest{
		Name:          "valid deneb blinded block",
		ConsensusData: *testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb),
	}
}
