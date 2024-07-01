package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// DenebBlockValidation tests a valid consensus data with deneb block
func DenebBlockValidation() *ValidatorConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ValidatorConsensusDataTest{
		Name:          "valid deneb block",
		ConsensusData: *testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb),
	}
}
