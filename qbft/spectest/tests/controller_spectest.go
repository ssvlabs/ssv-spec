package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

type DecidedState struct {
	DecidedVal []byte
	DecidedCnt uint
}

type RunInstanceData struct {
	InputValue           []byte
	InputMessages        []*types.SignedSSVMessage
	ControllerPostRoot   string
	ControllerPostState  types.Root `json:"-"` // Field is ignored by encoding/json
	ExpectedTimerState   *testingutils.TimerState
	ExpectedDecidedState DecidedState
	Height               *qbft.Height `json:"omitempty"`
}

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []*RunInstanceData
	OutputMessages  []*types.SignedSSVMessage
	ExpectedError   string
	StartHeight     *qbft.Height `json:"omitempty"`
}

func (test *ControllerSpecTest) TestName() string {
	return "qbft controller " + test.Name
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	//temporary to override state comparisons from file not inputted one
	test.overrideStateComparison(t)

	contr := test.generateController()

	if test.StartHeight != nil {
		contr.Height = *test.StartHeight
	}

	var lastErr error
	for i, runData := range test.RunInstanceData {
		height := qbft.Height(i)
		if runData.Height != nil {
			height = *runData.Height
		}
		if err := test.runInstanceWithData(t, height, contr, runData); err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *ControllerSpecTest) generateController() *qbft.Controller {
	identifier := []byte{1, 2, 3, 4}
	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)
	committeeMember := testingutils.TestingCommitteeMember(ks)
	return testingutils.NewTestingQBFTController(
		identifier[:],
		committeeMember,
		config,
		testingutils.TestingOperatorSigner(ks),
	)
}

func (test *ControllerSpecTest) testTimer(
	t *testing.T,
	config qbft.IConfig,
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

	return lastErr
}

func (test *ControllerSpecTest) runInstanceWithData(
	t *testing.T,
	height qbft.Height,
	contr *qbft.Controller,
	runData *RunInstanceData,
) error {
	err := contr.StartNewInstance(height, runData.InputValue)
	var lastErr error
	if err != nil {
		lastErr = err
	}

	test.testTimer(t, contr.GetConfig(), runData)

	if err := test.testProcessMsg(t, contr, runData); err != nil {
		lastErr = err
	}

	// test root
	r, err := contr.GetRoot()
	require.NoError(t, err)
	if runData.ControllerPostRoot != hex.EncodeToString(r[:]) {
		diff := typescomparable.PrintDiff(contr, runData.ControllerPostState)
		require.Fail(
			t,
			fmt.Sprintf(
				"post state not equal\nexpected: %s\nreceived: %s",
				runData.ControllerPostRoot,
				hex.EncodeToString(r[:]),
			),
			diff,
		)
	}

	return lastErr
}

func (test *ControllerSpecTest) overrideStateComparison(t *testing.T) {
	basedir, err := os.Getwd()
	require.NoError(t, err)
	basedir = filepath.Join(basedir, "generate")
	dir := typescomparable.GetSCDir(basedir, reflect.TypeOf(test).String())
	path := filepath.Join(dir, fmt.Sprintf("%s.json", test.TestName()))
	byteValue, err := os.ReadFile(filepath.Clean(path))
	require.NoError(t, err)
	sc := make([]*qbft.Controller, len(test.RunInstanceData))
	require.NoError(t, json.Unmarshal(byteValue, &sc))

	for i, runData := range test.RunInstanceData {
		runData.ControllerPostState = sc[i]

		r, err := sc[i].GetRoot()
		require.NoError(t, err)

		runData.ControllerPostRoot = hex.EncodeToString(r[:])
	}
}

func (test *ControllerSpecTest) GetPostState() (interface{}, error) {
	contr := test.generateController()

	if test.StartHeight != nil {
		contr.Height = *test.StartHeight
	}

	ret := make([]*qbft.Controller, len(test.RunInstanceData))
	for i, runData := range test.RunInstanceData {
		height := qbft.Height(i)
		if runData.Height != nil {
			height = *runData.Height
		}
		err := contr.StartNewInstance(height, runData.InputValue)
		if err != nil && len(test.ExpectedError) == 0 {
			return nil, err
		}

		for _, msg := range runData.InputMessages {
			_, err := contr.ProcessMsg(msg)
			if err != nil && len(test.ExpectedError) == 0 {
				return nil, err
			}
		}

		// copy controller
		byts, err := contr.Encode()
		if err != nil {
			return nil, err
		}
		copied := &qbft.Controller{}
		if err := copied.Decode(byts); err != nil {
			return nil, err
		}
		ret[i] = copied
	}
	return ret, nil
}
