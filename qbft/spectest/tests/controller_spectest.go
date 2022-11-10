package tests

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"reflect"
	"sort"
	"testing"
)

type RunInstanceData struct {
	InputValue         []byte
	InputMessages      []*qbft.SignedMessage
	DecidedVal         []byte
	DecidedCnt         uint
	SavedDecided       *qbft.SignedMessage
	BroadcastedDecided *qbft.SignedMessage
	ControllerPostRoot string
	ExpectedTimerState *testingutils.TimerState
}

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []*RunInstanceData
	OutputMessages  []*qbft.SignedMessage
	ExpectedError   string
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	config := testingutils.TestingConfig(testingutils.Testing4SharesSet())
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		config,
	)

	var lastErr error
	for _, runData := range test.RunInstanceData {
		err := contr.StartNewInstance(runData.InputValue)
		if err != nil {
			lastErr = err
		} else if runData.ExpectedTimerState != nil {
			if timer, ok := config.GetTimer().(*testingutils.TestQBFTTimer); ok {
				require.Equal(t, runData.ExpectedTimerState.Timeouts, timer.State.Timeouts)
				require.Equal(t, runData.ExpectedTimerState.Round, timer.State.Round)
			}
		}

		decidedCnt := 0
		for _, msg := range runData.InputMessages {
			decided, err := contr.ProcessMsg(msg)
			if err != nil {
				lastErr = err
			}
			if decided != nil {
				decidedCnt++

				data, _ := decided.Message.GetCommitData()
				require.EqualValues(t, runData.DecidedVal, data.Data)
			}
		}

		require.EqualValues(t, runData.DecidedCnt, decidedCnt)

		if runData.SavedDecided != nil {
			// test saved to storage
			decided, err := config.GetStorage().GetHighestDecided(identifier[:])
			require.NoError(t, err)
			require.NotNil(t, decided)
			r1, err := decided.GetRoot()
			require.NoError(t, err)

			r2, err := runData.SavedDecided.GetRoot()
			require.NoError(t, err)

			require.EqualValues(t, r2, r1)
			require.EqualValues(t, runData.SavedDecided.Signers, decided.Signers)
			require.EqualValues(t, runData.SavedDecided.Signature, decided.Signature)
		}
		if runData.BroadcastedDecided != nil {
			// test broadcasted
			broadcastedMsgs := config.GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
			require.Greater(t, len(broadcastedMsgs), 0)
			found := false
			for _, msg := range broadcastedMsgs {
				if !bytes.Equal(identifier[:], msg.MsgID[:]) {
					continue
				}

				msg1 := &qbft.SignedMessage{}
				require.NoError(t, msg1.Decode(msg.Data))
				r1, err := msg1.GetRoot()
				require.NoError(t, err)

				r2, err := runData.BroadcastedDecided.GetRoot()
				require.NoError(t, err)

				if bytes.Equal(r1, r2) &&
					reflect.DeepEqual(runData.BroadcastedDecided.Signers, msg1.Signers) &&
					reflect.DeepEqual(runData.BroadcastedDecided.Signature, msg1.Signature) {
					require.False(t, found)
					found = true
				}
			}
			require.True(t, found)
		}

		//r, err := contr.GetRoot()
		r, err := ControllerHistoricalRoot(contr)
		require.NoError(t, err)
		require.EqualValues(t, runData.ControllerPostRoot, hex.EncodeToString(r))
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

// ControllerHistoricalStruct returns ctrl historical root struct. TODO need to align all root in tests and remove this patch
func ControllerHistoricalStruct(ctrl *qbft.Controller) (interface{}, error) {
	rootStruct := struct {
		Identifier             []byte
		Height                 qbft.Height
		InstanceRoots          [][]byte
		HigherReceivedMessages map[types.OperatorID]qbft.Height
		Domain                 types.DomainType
		Share                  *types.Share
	}{
		Identifier:             ctrl.Identifier,
		Height:                 ctrl.Height,
		HigherReceivedMessages: ctrl.FutureMsgsContainer,
		Domain:                 ctrl.Domain,
		Share:                  ctrl.Share,
	}

	if ctrl.GetConfig() == nil {
		return rootStruct, nil
	}

	testingStorage := ctrl.GetConfig().GetStorage().(*testingutils.TestingStorage)
	states, err := testingStorage.GetAllInstancesState(ctrl.Identifier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all instances state")
	}

	// sort in order maintain the same root
	sort.Slice(states, func(i, j int) bool {
		return states[i].Height > states[j].Height
	})

	var roots [][]byte
	for i := 0; i < len(states); i++ {
		r, err := qbft.NewInstanceFromState(ctrl.GetConfig(), states[i]).GetRoot()
		if err != nil {
			return nil, errors.Wrap(err, "failed getting instance root")
		}
		roots = append(roots, r)
	}

	rootStruct.InstanceRoots = roots
	return rootStruct, nil
}

// ControllerHistoricalRoot supports historical root. TODO need to align all root in tests and remove this patch
func ControllerHistoricalRoot(ctrl *qbft.Controller) ([]byte, error) {
	rootStruct, err := ControllerHistoricalStruct(ctrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ctrl root struct")
	}
	r, err := json.Marshal(rootStruct)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(r)
	return ret[:], nil
}
