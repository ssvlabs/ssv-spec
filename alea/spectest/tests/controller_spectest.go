package tests

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/alea"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/stretchr/testify/require"
)

type DecidedState struct {
	DecidedVal               []byte
	DecidedCnt               uint
	BroadcastedDecided       *alea.SignedMessage
	CalledSyncDecidedByRange bool
	DecidedByRangeValues     [2]alea.Height
}

type RunInstanceData struct {
	InputValue           []byte
	InputMessages        []*alea.SignedMessage
	ControllerPostRoot   string
	ExpectedTimerState   *testingutils.TimerStateAlea
	ExpectedDecidedState DecidedState
}

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []*RunInstanceData
	OutputMessages  []*alea.SignedMessage
	ExpectedError   string
}

func (test *ControllerSpecTest) TestName() string {
	return "alea controller " + test.Name
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	config := testingutils.TestingConfigAlea(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingAleaController(
		identifier[:],
		testingutils.TestingShareAlea(testingutils.Testing4SharesSet()),
		config,
	)

	var lastErr error
	for _, runData := range test.RunInstanceData {
		if err := test.runInstanceWithData(t, contr, config, identifier, runData); err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *ControllerSpecTest) testTimer(
	t *testing.T,
	config *alea.Config,
	runData *RunInstanceData,
) {
	if runData.ExpectedTimerState != nil {
		if timer, ok := config.GetTimer().(*testingutils.TestQBFTTimerAlea); ok {
			require.Equal(t, runData.ExpectedTimerState.Timeouts, timer.State.Timeouts)
			require.Equal(t, runData.ExpectedTimerState.Round, timer.State.Round)
		}
	}
}

func (test *ControllerSpecTest) testProcessMsg(
	t *testing.T,
	contr *alea.Controller,
	config *alea.Config,
	runData *RunInstanceData,
) error {
	// decidedCnt := 0
	var lastErr error
	for _, msg := range runData.InputMessages {
		_, err := contr.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
		// if decided != nil {
		// 	decidedCnt++

		// 	data, _ := decided.Message.GetCommitData()
		// 	require.EqualValues(t, runData.ExpectedDecidedState.DecidedVal, data.Data)
		// }
	}
	// require.EqualValues(t, runData.ExpectedDecidedState.DecidedCnt, decidedCnt)

	// verify sync decided by range calls
	// if runData.ExpectedDecidedState.CalledSyncDecidedByRange {
	// 	require.EqualValues(t, runData.ExpectedDecidedState.DecidedByRangeValues, config.GetNetwork().(*testingutils.TestingNetworkAlea).DecidedByRange)
	// } else {
	// 	require.EqualValues(t, [2]alea.Height{0, 0}, config.GetNetwork().(*testingutils.TestingNetworkAlea).DecidedByRange)
	// }

	return lastErr
}

func (test *ControllerSpecTest) testBroadcastedDecided(
	t *testing.T,
	config *alea.Config,
	identifier types.MessageID,
	runData *RunInstanceData,
) {
	if runData.ExpectedDecidedState.BroadcastedDecided != nil {
		// test broadcasted
		broadcastedMsgs := config.GetNetwork().(*testingutils.TestingNetworkAlea).BroadcastedMsgs
		require.Greater(t, len(broadcastedMsgs), 0)
		found := false
		for _, msg := range broadcastedMsgs {
			if !bytes.Equal(identifier[:], msg.MsgID[:]) {
				continue
			}

			msg1 := &alea.SignedMessage{}
			require.NoError(t, msg1.Decode(msg.Data))
			r1, err := msg1.GetRoot()
			require.NoError(t, err)

			r2, err := runData.ExpectedDecidedState.BroadcastedDecided.GetRoot()
			require.NoError(t, err)

			if bytes.Equal(r1, r2) &&
				reflect.DeepEqual(runData.ExpectedDecidedState.BroadcastedDecided.Signers, msg1.Signers) &&
				reflect.DeepEqual(runData.ExpectedDecidedState.BroadcastedDecided.Signature, msg1.Signature) {
				require.False(t, found)
				found = true
			}
		}
		require.True(t, found)
	}
}

func (test *ControllerSpecTest) runInstanceWithData(
	t *testing.T,
	contr *alea.Controller,
	config *alea.Config,
	identifier types.MessageID,
	runData *RunInstanceData,
) error {
	err := contr.StartNewInstance(runData.InputValue)
	var lastErr error
	if err != nil {
		lastErr = err
	}

	test.testTimer(t, config, runData)

	if err := test.testProcessMsg(t, contr, config, runData); err != nil {
		lastErr = err
	}

	test.testBroadcastedDecided(t, config, identifier, runData)

	// test root
	r, err := contr.GetRoot()
	require.NoError(t, err)
	require.EqualValues(t, runData.ControllerPostRoot, hex.EncodeToString(r))

	return lastErr
}
