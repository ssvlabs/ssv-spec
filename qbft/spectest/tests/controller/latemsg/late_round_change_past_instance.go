package latemsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// LateRoundChangePastInstance tests process round change msg for a previously decided instance
func LateRoundChangePastInstance() *tests.ControllerSpecTest {
	ks := testingutils.Testing4SharesSet()

	allMsgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, 5, ks)

	msgPerHeight := make(map[qbft.Height][]*qbft.SignedMessage)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]
	msgPerHeight[2] = allMsgs[14:21]
	msgPerHeight[3] = allMsgs[21:28]
	msgPerHeight[4] = allMsgs[28:35]
	msgPerHeight[5] = allMsgs[35:42]

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue:    []byte{1, 2, 3, 4},
			InputMessages: msgPerHeight[height],
			ExpectedDecidedState: tests.DecidedState{
				BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithHeight(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					height,
				),
				DecidedVal: testingutils.TestingQBFTFullData,
				DecidedCnt: 1,
			},

			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "late round change past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "24cf697092529cfab3ab06b969d8696692c8bcbb9f41a954f71dc74c3b1d7e97"),
			instanceData(1, "8aa5464b119518f178d81edf4cea1f4c918f9e084e5262a0e276d3afb00ba620"),
			instanceData(2, "1799fe0981ae08bde1eae9fef88ef8035f5952974647786287a1a8c36544a5da"),
			instanceData(3, "2ddf2b8c2f35c8115ddd68120c71e64809bfea6b023ed13e177c2474a95d137d"),
			instanceData(4, "318816cc8819ad062996704fb4b9990b8088ade0cb4c26816ea0965783bff12e"),
			instanceData(5, "fd83cdee705de628cf4a9baf9e662c424bf0942c63a680266d7db872d16e9f0a"),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingMultiSignerRoundChangeMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						4,
					),
				},
				ControllerPostRoot: "25270cd1fd958682fc7602ddcabbbd2f384a9c28fca960e214268867caa23bd7",
			},
		},
	}
}
