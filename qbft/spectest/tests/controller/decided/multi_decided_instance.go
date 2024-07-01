package decided

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// MultiDecidedInstances tests deciding multiple instances
func MultiDecidedInstances() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue: []byte{1, 2, 3, 4},
			InputMessages: []*types.SignedSSVMessage{
				testingutils.TestingCommitMultiSignerMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]}, []types.OperatorID{1, 2, 3}, height),
			},
			ExpectedDecidedState: tests.DecidedState{
				DecidedCnt: 1,
				DecidedVal: testingutils.TestingQBFTFullData,
			},
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "multi decide instances",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "f552f5aedb2e0d7933e77c4297c69e761000e88f78ae02e0afd4d053847b8d5c"),
			instanceData(1, "d0e04e5bce1d0e75def07c8b1917981b86fa25e0d488b5ed365be477ee6a6298"),
			instanceData(2, "248ad9f2454d5db1f7060c9755374e4ceab8fdb51c1030f3dc6d1b9492155c85"),
			instanceData(3, "de6d0efcb3e55e1e33e850acf87933b71fa93f2e7f541b6e3bf139d2dd9740eb"),
			instanceData(4, "7f595bfb4a7e5d6a1cfdd2fed35a55ed57d23de63bbf3fd94f9392829600e0c3"),
			instanceData(5, "f9d585fd5629b1e619704c66eb2f4b5bb248504d504f830cb63adcad37bcf46e"),
			instanceData(6, "53e58c1a850e79e9a21d11d39a68ed86a89f5fe84a625d5c6f29064b70838d0a"),
			instanceData(7, "1bd069bcea9ee80e4b39603d636346e72dd32b599f7e14eaf846a7fc1e438b39"),
			instanceData(8, "7a7b12431f8880f3f09b7a40f2d88928596e1b692a505813483450e5c637a7dd"),
			instanceData(9, "0cdc302065fe16d722a3a1db04c07edb8139f7c53aafb847e1aeb89797eb4282"),
			instanceData(10, "d1383c9a269c18df528f41f2563322a0b3bf490b4625f7fdc9277aca53dacc03"),
		},
	}
}
