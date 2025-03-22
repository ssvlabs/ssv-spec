package validatorconsensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// CBSigning tests an invalid consensus data for commit boost signing (has no consensus data)
func CBSigning() *ValidatorConsensusDataTest {

	dataByts, err := testingutils.TestingCBSigningRequest.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	cd := types.ValidatorConsensusData{
		Duty:    testingutils.TestingCBSigningDuty.Duty,
		Version: spec.DataVersionCapella,
		DataSSZ: dataByts,
	}

	return &ValidatorConsensusDataTest{
		Name:          "commit boost signing",
		ConsensusData: cd,
		ExpectedError: "commit-boost signing has no consensus data",
	}
}
