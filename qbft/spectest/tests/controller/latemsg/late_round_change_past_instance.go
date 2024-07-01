package latemsg

import (
	"crypto/rsa"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// LateRoundChangePastInstance tests process round change msg for a previously decided instance
func LateRoundChangePastInstance() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	allMsgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, 1, ks)

	msgPerHeight := make(map[qbft.Height][]*types.SignedSSVMessage)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]

	instanceData := func(height qbft.Height) *tests.RunInstanceData {
		sc := lateRoundChangePastInstanceStateComparison(height, nil)
		return &tests.RunInstanceData{
			InputValue:    []byte{1, 2, 3, 4},
			InputMessages: msgPerHeight[height],
			ExpectedDecidedState: tests.DecidedState{
				BroadcastedDecided: testingutils.TestingCommitMultiSignerMessageWithHeight(
					[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
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

	lateMsg := testingutils.TestingMultiSignerRoundChangeMessageWithHeight([]*rsa.PrivateKey{ks.OperatorKeys[4]}, []types.OperatorID{4}, qbft.FirstHeight)
	sc := lateRoundChangePastInstanceStateComparison(2, lateMsg)

	return &tests.ControllerSpecTest{
		Name: "late round change past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight),
			instanceData(1),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*types.SignedSSVMessage{
					testingutils.TestingMultiSignerRoundChangeMessageWithHeight(
						[]*rsa.PrivateKey{ks.OperatorKeys[4]},
						[]types.OperatorID{4},
						qbft.FirstHeight,
					),
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
		ExpectedError: "could not process msg: instance stopped processing messages",
	}
}

// lateRoundChangePastInstanceStateComparison returns a comparable.StateComparison for controller running up to the given height.
// lateMsg will be added to the rc container of the instance at the proper height.
func lateRoundChangePastInstanceStateComparison(height qbft.Height, lateMsg *types.SignedSSVMessage) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	allMsgs := testingutils.ExpectedDecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData, testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, 1, ks)
	offset := 7 // 7 messages per height (1 propose + 3 prepare + 3 commit)

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	for i := 0; i <= int(height); i++ {
		contr.Height = qbft.Height(i)

		instance := &qbft.Instance{
			StartValue: []byte{1, 2, 3, 4},
			State: &qbft.State{
				CommitteeMember: testingutils.TestingCommitteeMember(testingutils.Testing4SharesSet()),
				ID:              testingutils.TestingIdentifier,
				Round:           qbft.FirstRound,
				Height:          qbft.Height(i),
			},
		}

		// last height
		if lateMsg != nil && qbft.Height(i) == height {
			comparable.InitContainers(instance)
			contr.StoredInstances = append([]*qbft.Instance{instance}, contr.StoredInstances...)
			break
		}

		instance.State.ProposalAcceptedForCurrentRound = testingutils.ToProcessingMessage(testingutils.TestingProposalMessageWithParams(ks.OperatorKeys[1], types.OperatorID(1), qbft.FirstRound, qbft.Height(i), testingutils.TestingQBFTRootData, nil, nil))
		instance.State.LastPreparedRound = qbft.FirstRound
		instance.State.LastPreparedValue = testingutils.TestingQBFTFullData
		instance.State.Decided = true
		instance.State.DecidedValue = testingutils.TestingQBFTFullData
		if qbft.Height(i) != height {
			instance.ForceStop()
		}
		msgs := allMsgs[offset*i : offset*(i+1)]
		comparable.SetSignedMessages(instance, msgs)

		contr.StoredInstances = append([]*qbft.Instance{instance}, contr.StoredInstances...)
	}

	return &comparable.StateComparison{ExpectedState: contr}
}
