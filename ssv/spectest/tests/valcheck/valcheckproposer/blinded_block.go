package valcheckproposer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// BlindedBlock tests if blinded blocks pass validation according to configuration
func BlindedBlock() tests.SpecTest {
	return &valcheck.MultiSpecTest{
		Name: "blinded blocks",
		Tests: []*valcheck.SpecTest{
			{
				Name:            "blinded blocks not allowed",
				Network:         types.TestNetwork,
				BeaconRole:      types.BNRoleProposer,
				Input:           testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionBellatrix),
				SupportsBlinded: false,
				ExpectedError:   "blinded blocks are not supported",
			},
			{
				Name:            "blinded blocks allowed",
				Network:         types.TestNetwork,
				BeaconRole:      types.BNRoleProposer,
				Input:           testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionBellatrix),
				SupportsBlinded: true,
				AnyError:        false,
			},
		},
	}
}
