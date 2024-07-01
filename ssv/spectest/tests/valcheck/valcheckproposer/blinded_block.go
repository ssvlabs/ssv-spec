package valcheckproposer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// BlindedBlock tests if blinded blocks pass validation according to configuration
func BlindedBlock() tests.SpecTest {
	return &valcheck.MultiSpecTest{
		Name: "blinded blocks",
		Tests: []*valcheck.SpecTest{
			{
				Name:       "blinded blocks accepted",
				Network:    types.BeaconTestNetwork,
				RunnerRole: types.RoleProposer,
				Input:      testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
				AnyError:   false,
			},
		},
	}
}
