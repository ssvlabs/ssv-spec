package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// LateRoundChangePastInstance tests process round change msg for a previously decided instance
func LateRoundChangePastInstance() tests.SpecTest {
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

	instanceData := func(height qbft.Height) *tests.RunInstanceData {
		sc := lateRoundChangePastInstanceStateComparison(height, nil)
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
			ControllerPostRoot:  sc.Root(),
			ControllerPostState: sc.ExpectedState,
		}
	}

	lateMsg := testingutils.TestingMultiSignerRoundChangeMessageWithHeight([]*bls.SecretKey{ks.Shares[4]}, []types.OperatorID{4}, 4)
	sc := lateRoundChangePastInstanceStateComparison(6, lateMsg)

	return &tests.ControllerSpecTest{
		Name: "late round change past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight),
			instanceData(1),
			instanceData(2),
			instanceData(3),
			instanceData(4),
			instanceData(5),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingMultiSignerRoundChangeMessageWithHeight(
						[]*bls.SecretKey{ks.Shares[4]},
						[]types.OperatorID{4},
						4,
					),
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
	}
}

func lateRoundChangePastInstanceStateComparison(height qbft.Height, lateMsg *qbft.SignedMessage) *qbftcomparable.StateComparison {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)

	ks := testingutils.Testing4SharesSet()
	allMsgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, 5, ks)

	for i := 0; i <= int(height); i++ {
		contr.Height = qbft.Height(i)
		_ = contr.StartNewInstance([]byte{1, 2, 3, 4})

		offset := 7 * i
		msgs := allMsgs[offset : offset+7]

		state := testingutils.BaseInstance().State
		state.Height = qbft.Height(i)

		// last height
		if lateMsg != nil && i == int(height) {
			contr.StoredInstances[0].State = state
			break
		}

		state.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.Height(i), testingutils.TestingQBFTRootData, nil, nil)
		state.ProposalAcceptedForCurrentRound.Message.Height = qbft.Height(i)
		state.LastPreparedRound = 1
		state.LastPreparedValue = testingutils.TestingQBFTFullData
		state.Decided = true
		state.DecidedValue = testingutils.TestingQBFTFullData
		state.ProposeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
			qbft.FirstRound: {
				msgs[0],
			},
		}}
		state.PrepareContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
			qbft.FirstRound: {
				msgs[1],
				msgs[2],
				msgs[3],
			},
		}}
		state.CommitContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
			qbft.FirstRound: {
				msgs[4],
				msgs[5],
				msgs[6],
			},
		}}

		if lateMsg != nil && qbft.Height(i) == lateMsg.Message.Height {
			state.RoundChangeContainer = &qbft.MsgContainer{Msgs: map[qbft.Round][]*qbft.SignedMessage{
				qbft.FirstRound: {
					lateMsg,
				},
			}}
		}

		contr.StoredInstances[0].State = state
	}

	return &qbftcomparable.StateComparison{ExpectedState: contr}
}
