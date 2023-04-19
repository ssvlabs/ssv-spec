package tests

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	qbftcomparable "github.com/bloxapp/ssv-spec/qbft/spectest/comparable"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	typescomparable "github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

type DecidedState struct {
	DecidedVal               []byte
	DecidedCnt               uint
	BroadcastedDecided       *qbft.SignedMessage
	CalledSyncDecidedByRange bool
	DecidedByRangeValues     [2]qbft.Height
}

type RunInstanceData struct {
	InputValue           []byte
	InputMessages        []*qbft.SignedMessage
	ControllerPostRoot   string
	ExpectedTimerState   *testingutils.TimerState
	ExpectedDecidedState DecidedState
}

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []*RunInstanceData
	OutputMessages  []*qbft.SignedMessage
	ExpectedError   string
}

func (test *ControllerSpecTest) TestName() string {
	return "qbft controller " + test.Name
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	identifier := []byte{1, 2, 3, 4}
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
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
	config *qbft.Config,
	runData *RunInstanceData,
) {
	if runData.ExpectedTimerState != nil {
		if timer, ok := config.GetTimer().(*testingutils.TestQBFTTimer); ok {
			require.Equal(t, runData.ExpectedTimerState.Timeouts, timer.State.Timeouts)
			require.Equal(t, runData.ExpectedTimerState.Round, timer.State.Round)
		}
	}
}

func (test *ControllerSpecTest) testProcessMsg(
	t *testing.T,
	contr *qbft.Controller,
	config *qbft.Config,
	runData *RunInstanceData,
) error {
	decidedCnt := 0
	var lastErr error
	for _, msg := range runData.InputMessages {
		decided, err := contr.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
		if decided != nil {
			decidedCnt++

			require.EqualValues(t, runData.ExpectedDecidedState.DecidedVal, decided.FullData)
		}
	}
	require.EqualValues(t, runData.ExpectedDecidedState.DecidedCnt, decidedCnt)

	// verify sync decided by range calls
	if runData.ExpectedDecidedState.CalledSyncDecidedByRange {
		require.EqualValues(t, runData.ExpectedDecidedState.DecidedByRangeValues, config.GetNetwork().(*testingutils.TestingNetwork).DecidedByRange)
	} else {
		require.EqualValues(t, [2]qbft.Height{0, 0}, config.GetNetwork().(*testingutils.TestingNetwork).DecidedByRange)
	}

	return lastErr
}

func (test *ControllerSpecTest) testBroadcastedDecided(
	t *testing.T,
	config *qbft.Config,
	identifier []byte,
	runData *RunInstanceData,
) {
	if runData.ExpectedDecidedState.BroadcastedDecided != nil {
		// test broadcasted
		broadcastedMsgs := config.GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
		require.Greater(t, len(broadcastedMsgs), 0)
		found := false
		for _, msg := range broadcastedMsgs {

			// a hack for testing non standard messageID identifiers since we copy them into a MessageID this fixes it
			msgID := types.MessageID{}
			copy(msgID[:], identifier)

			if !bytes.Equal(msgID[:], msg.MsgID[:]) {
				continue
			}

			msg1 := &qbft.SignedMessage{}
			require.NoError(t, msg1.Decode(msg.Data))
			r1, err := msg1.GetRoot()
			require.NoError(t, err)

			r2, err := runData.ExpectedDecidedState.BroadcastedDecided.GetRoot()
			require.NoError(t, err)

			if r1 == r2 &&
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
	contr *qbft.Controller,
	config *qbft.Config,
	identifier []byte,
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
	if runData.ControllerPostRoot != hex.EncodeToString(r) {
		diff := typescomparable.PrintDiff(contr, qbftcomparable.RootRegister[runData.ControllerPostRoot])
		require.Fail(t, "post state not equal", diff)
	}
	return lastErr
}
