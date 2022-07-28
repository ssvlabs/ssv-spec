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
		DecidedCnt    uint
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

		decidedCnt := 0
		for _, msg := range runData.InputMessages {
			decided, _, err := contr.ProcessMsg(msg)
			if err != nil {
				lastErr = err
			}
			if decided {
				decidedCnt++
			}
		}

		require.EqualValues(t, runData.DecidedCnt, decidedCnt)

		isDecided, decidedVal := contr.InstanceForHeight(contr.Height).IsDecided()
		require.EqualValues(t, runData.Decided, isDecided)
		require.EqualValues(t, runData.DecidedVal, decidedVal)
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *ControllerSpecTest) TestName() string {
	return "qbft controller " + test.Name
}
