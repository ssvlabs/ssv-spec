package tests

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

// ChangeProposerFuncInstanceHeight tests with this height will return proposer operator ID 2
const ChangeProposerFuncInstanceHeight = 10

type MsgProcessingSpecTest struct {
	Name               string
	Pre                *qbft.Instance
	PostRoot           string
	PostState          types.Root `json:"-"` // Field is ignored by encoding/json
	InputMessages      []*types.SignedSSVMessage
	OutputMessages     []*types.SignedSSVMessage
	ExpectedError      string
	ExpectedTimerState *testingutils.TimerState
}

func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	// temporary to override state comparisons from file not inputted one
	test.overrideStateComparison(t)

	lastErr := test.runPreTesting()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	if test.ExpectedTimerState != nil {
		// checks round timer state
		timer, ok := test.Pre.GetConfig().GetTimer().(*testingutils.TestQBFTTimer)
		if ok && timer != nil {
			require.Equal(t, test.ExpectedTimerState.Timeouts, timer.State.Timeouts, "timer should have expected timeouts count")
			require.Equal(t, test.ExpectedTimerState.Round, timer.State.Round, "timer should have expected round")
		}
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)

	// test output message
	broadcastedSignedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	testingutils.CompareSignedSSVMessageOutputMessages(t, test.OutputMessages, broadcastedSignedMsgs, test.Pre.State.CommitteeMember.Committee)

	// test root
	if test.PostRoot != hex.EncodeToString(postRoot[:]) {
		diff := typescomparable.PrintDiff(test.Pre.State, test.PostState)
		require.Fail(t, fmt.Sprintf("expected root: %s\nactual root: %s\n\n", test.PostRoot, hex.EncodeToString(postRoot[:])), "post state not equal", diff)
	}
}

func (test *MsgProcessingSpecTest) runPreTesting() error {
	// a simple hack to change the proposer func
	if test.Pre.State.Height == ChangeProposerFuncInstanceHeight {
		test.Pre.GetConfig().(*qbft.Config).ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
			return 2
		}
	}

	var lastErr error
	for _, msg := range test.InputMessages {
		_, _, _, err := test.Pre.ProcessMsg(testingutils.ToProcessingMessage(msg))
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (test *MsgProcessingSpecTest) TestName() string {
	return "qbft message processing " + test.Name
}

func (test *MsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	basedir, err := os.Getwd()
	require.NoError(t, err)
	test.PostState, err = typescomparable.UnmarshalStateComparison(basedir, test.TestName(),
		reflect.TypeOf(test).String(),
		&qbft.State{})
	require.NoError(t, err)

	r, err := test.PostState.GetRoot()
	require.NoError(t, err)

	// backwards compatability test, hard coded post root must be equal to the one loaded from file
	if len(test.PostRoot) > 0 {
		require.EqualValues(t, test.PostRoot, hex.EncodeToString(r[:]))
	}

	test.PostRoot = hex.EncodeToString(r[:])
}

func (test *MsgProcessingSpecTest) GetPostState() (interface{}, error) {
	err := test.runPreTesting()
	if err != nil && len(test.ExpectedError) == 0 { // only non expected errors should return error
		return nil, err
	}
	return test.Pre.State, nil
}
