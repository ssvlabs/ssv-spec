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
			instanceData(1, "dddf0ce6fdcfafaed531289af68c47b7e362609336943f55124ef3fcde7d6c99"),
			instanceData(2, "c9da6816ee973d84e0927f8ee4b754b66ee64517f40b3475a89a8476c2cae587"),
			instanceData(3, "8573f8edb90224af61512c2dec830e7036d4a1cb4e8f2b6861fd16966472a593"),
			instanceData(4, "91ebf1bea18bc017d5e7690258c5e74ce1988a2de02587528ba2224b62fff3ec"),
			instanceData(5, "67e49a573df14c606ffcad73692a7e9b2e85e5571d1237f426f2e504e317b888"),
			instanceData(6, "442b28e2561527b0be0672e523bda4ee87d3142858b99af9989e08afe8d3fb99"),
			instanceData(7, "0339bc66256af2769607959a8bccd54af89a12d36849e389262ccd657ea18439"),
			instanceData(8, "019bcea934afcfe34641fe0414638120599c8f3708d5ef29b9a9ed1cdba89505"),
			instanceData(9, "d8c8c67ecac985437772630107a96a00bf76014f8f2ab8a30533f4a64f03239e"),
			instanceData(10, "2e0b9406ce11f5174ae5be2de092290fecc07d933d0927a4a626142960168c97"),
		},
	}
}
