package tests

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []struct {
		InputValue    []byte
		InputMessages []*qbft.SignedMessage
		Decided       bool
		DecidedVal    []byte
	}
	ValCheck       qbft.ProposedValueCheckF
	OutputMessages []*qbft.SignedMessage
	ExpectedError  string
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		test.ValCheck,
		func(state *qbft.State, round qbft.Round) types.OperatorID {
			return 1
		},
	)

	var lastErr error
	for _, runData := range test.RunInstanceData {
		startedInstance := false
		err := contr.StartNewInstance(runData.InputValue)
		if err != nil {
			lastErr = err
		} else {
			startedInstance = true
		}

		if !startedInstance {
			continue
		}

		for _, msg := range runData.InputMessages {
			_, _, err := contr.ProcessMsg(msg)
			if err != nil {
				lastErr = err
			}
		}

		isDecided, decidedVal := contr.InstanceForHeight(contr.Height).IsDecided()
		require.EqualValues(t, runData.Decided, isDecided)
		require.EqualValues(t, runData.DecidedVal, decidedVal)
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	//var lastErr error
	//for _, msg := range test.InputMessages {
	//	_, _, _, err := test.Pre.ProcessMsg(msg)
	//	if err != nil {
	//		lastErr = err
	//	}
	//}
	//
	//if len(test.ExpectedError) != 0 {
	//	require.EqualError(t, lastErr, test.ExpectedError)
	//} else {
	//	require.NoError(t, lastErr)
	//}
	//
	//postRoot, err := test.Pre.State.GetRoot()
	//require.NoError(t, err)
	//
	//// test output message
	//broadcastedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	//if len(test.OutputMessages) > 0 || len(broadcastedMsgs) > 0 {
	//	require.Len(t, broadcastedMsgs, len(test.OutputMessages))
	//
	//	for i, msg := range test.OutputMessages {
	//		r1, _ := msg.GetRoot()
	//
	//		msg2 := &qbft.SignedMessage{}
	//		require.NoError(t, msg2.Decode(broadcastedMsgs[i].Data))
	//
	//		r2, _ := msg2.GetRoot()
	//		require.EqualValues(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
	//	}
	//}
	//
	//require.EqualValues(t, test.PostRoot, hex.EncodeToString(postRoot), "post root not valid")
}

func (test *ControllerSpecTest) TestName() string {
	return test.Name
}
