package attester

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
)

// ValidDecided tests a decided msg received
func ValidDecided() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	dr := testingutils.AttesterRunner(ks)

	msgs := []*types.SSVMessage{
		testingutils.SSVMsgAttester(testingutils.SignQBFTMsg(ks.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.AttesterMsgID,
			Data:       testingutils.ProposalDataBytes(testingutils.TestAttesterConsensusDataByts, nil, nil),
		}), nil),
		testingutils.SSVMsgAttester(testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{
				ks.Shares[1],
				ks.Shares[2],
				ks.Shares[3],
			},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: testingutils.AttesterMsgID,
				Data:       testingutils.CommitDataBytes(testingutils.TestAttesterConsensusDataByts),
			}), nil),

		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)),
		testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[3], 3, qbft.FirstHeight)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:                    "attester valid decided",
		Runner:                  dr,
		Duty:                    testingutils.TestAttesterConsensusData.Duty,
		Messages:                msgs,
		PostDutyRunnerStateRoot: "d4d043fe512268cd3bdda612c5414323cf599c1f480adda63ce76892b834ed83",
		OutputMessages: []*ssv.SignedPartialSignatureMessage{
			testingutils.PostConsensusAttestationMsg(testingutils.Testing4SharesSet().Shares[1], 1, qbft.FirstHeight),
		},
	}
}
