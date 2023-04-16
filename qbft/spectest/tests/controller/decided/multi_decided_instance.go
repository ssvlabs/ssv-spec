package decided

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// MultiDecidedInstances tests deciding multiple instances
func MultiDecidedInstances() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue: []byte{1, 2, 3, 4},
			InputMessages: []*qbft.SignedMessage{
				testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, height),
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
			instanceData(1, "e8da0b4b1afeee3b611e2c91feb6a66b48dc6287728ab6211c40bac76282ebc4"),
			instanceData(2, "1fcc752b2ec7645dd4f46c50c7e0260699aceebfe03c8812cfd5aaaa059a9659"),
			instanceData(3, "64dd8825012aae061b0d2cdd36e6379f1b30c1fb8f7d8db2c800161565e9bee1"),
			instanceData(4, "dbafdd3370a7cb493058842a2036b9a95a6aa0040300f1cf96ab74f086eb7749"),
			instanceData(5, "562a4c43ebed97b9f6f24fbb2fa5c9850f6ed921187e70c667143028250fb301"),
			instanceData(6, "489ca498ca6fb5e38e6b28a2bb54fb5b8ac57892647972e4c91ce48aa8d17950"),
			instanceData(7, "c83067c025a4222fe887fecf099f0a8a6c4372bafce05a588bbb11f640700910"),
			instanceData(8, "fd9008ae027eff439bc2acfeeaab368ecd93d7964b1e1ddd7c760bdebfde0a92"),
			instanceData(9, "4bc1dca94e40137aa7e000ad13b9b6eca02f7a6b31f3730971dfc162f5c941f5"),
			instanceData(10, "a76efa101019ba414067c5918e419a9dcc31a348035c7357fa073852d46d98fc"),
		},
	}
}
