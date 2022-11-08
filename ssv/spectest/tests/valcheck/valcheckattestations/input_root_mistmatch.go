package valcheckattestations

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InputRootMismatch tests ConsensusData.Root != computed root
func InputRootMismatch() *valcheck.SpecTest {
	return &valcheck.SpecTest{
		Name:       "attestation value check root mismatch",
		Network:    types.PraterNetwork,
		BeaconRole: types.BNRoleAttester,
		Input: &qbft.Data{
			Root:   [32]byte{1, 2, 3, 4},
			Source: testingutils.TestAttesterConsensusDataByts,
		},
		ExpectedError: "invalid input data: msg root data != calculated root data",
	}
}
