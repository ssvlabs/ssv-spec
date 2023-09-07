package newduty

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/bloxapp/ssv-spec/types/testingutils/comparable"
)

type StartNewRunnerDutySpecTest struct {
	Name                    string
	Runner                  ssv.Runner
	Duty                    *types.Duty
	PostDutyRunnerStateRoot string
	PostDutyRunnerState     types.Root `json:"-"` // Field is ignored by encoding/json
	OutputMessages          []*types.SignedPartialSignatureMessage
	ExpectedError           string
}

func (test *StartNewRunnerDutySpecTest) TestName() string {
	return test.Name
}

// overrideStateComparison overrides the state comparison to compare the runner state
func (test *StartNewRunnerDutySpecTest) overrideStateComparison(t *testing.T) {
	overrideStateComparison(t, test, test.Name, reflect.TypeOf(test).String())
}

// RunAsPartOfMultiTest runs the test as part of a MultiMsgProcessingSpecTest.
// It simply runs without calling oveerideStateComparison
func (test *StartNewRunnerDutySpecTest) RunAsPartOfMultiTest(t *testing.T) {
	err := test.runPreTesting()
	if len(test.ExpectedError) > 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	// test output message
	broadcastedMsgs := test.Runner.GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	if len(broadcastedMsgs) > 0 {
		index := 0
		for _, msg := range broadcastedMsgs {
			if msg.MsgType != types.SSVPartialSignatureMsgType {
				continue
			}

			msg1 := &types.SignedPartialSignatureMessage{}
			require.NoError(t, msg1.Decode(msg.Data))
			msg2 := test.OutputMessages[index]
			require.Len(t, msg1.Message.Messages, len(msg2.Message.Messages))

			// messages are not guaranteed to be in order so we map them and then test all roots to be equal
			roots := make(map[string]string)
			for i, partialSigMsg2 := range msg2.Message.Messages {
				r2, err := partialSigMsg2.GetRoot()
				require.NoError(t, err)
				if _, found := roots[hex.EncodeToString(r2[:])]; !found {
					roots[hex.EncodeToString(r2[:])] = ""
				} else {
					roots[hex.EncodeToString(r2[:])] = hex.EncodeToString(r2[:])
				}

				partialSigMsg1 := msg1.Message.Messages[i]
				r1, err := partialSigMsg1.GetRoot()
				require.NoError(t, err)

				if _, found := roots[hex.EncodeToString(r1[:])]; !found {
					roots[hex.EncodeToString(r1[:])] = ""
				} else {
					roots[hex.EncodeToString(r1[:])] = hex.EncodeToString(r1[:])
				}
			}
			for k, v := range roots {
				require.EqualValues(t, k, v, "missing output msg")
			}

			index++
		}

		require.Len(t, test.OutputMessages, index)
	}

	// post root
	postRoot, err := test.Runner.GetRoot()
	require.NoError(t, err)

	if test.PostDutyRunnerStateRoot != hex.EncodeToString(postRoot[:]) {
		diff := comparable.PrintDiff(test.Runner, test.PostDutyRunnerState)
		require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(postRoot[:]), fmt.Sprintf("post runner state not equal\n%s\n", diff))
	}
}

func (test *StartNewRunnerDutySpecTest) Run(t *testing.T) {
	test.overrideStateComparison(t)
	test.RunAsPartOfMultiTest(t)
}

// runPreTesting runs the spec logic before testing the output
// It simply starts a new duty
func (test *StartNewRunnerDutySpecTest) runPreTesting() error {
	err := test.Runner.StartNewDuty(test.Duty)
	return err
}

func (test *StartNewRunnerDutySpecTest) GetPostState() (interface{}, error) {
	err := test.runPreTesting()
	return test.Runner, err
}

type MultiStartNewRunnerDutySpecTest struct {
	Name  string
	Tests []*StartNewRunnerDutySpecTest
}

func (tests *MultiStartNewRunnerDutySpecTest) TestName() string {
	return tests.Name
}

func (tests *MultiStartNewRunnerDutySpecTest) Run(t *testing.T) {
	tests.overrideStateComparison(t)

	for _, test := range tests.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.RunAsPartOfMultiTest(t)
		})
	}
}

func (tests *MultiStartNewRunnerDutySpecTest) GetPostState() (interface{}, error) {
	ret := make(map[string]types.Root, len(tests.Tests))
	for _, test := range tests.Tests {
		err := test.runPreTesting()
		if err != nil && test.ExpectedError != err.Error() {
			return nil, err
		}
		ret[test.Name] = test.Runner
	}
	return ret, nil
}

// overrideStateComparison overrides the post state comparison for all tests in the multi test
func (tests *MultiStartNewRunnerDutySpecTest) overrideStateComparison(t *testing.T) {
	testsName := strings.ReplaceAll(tests.TestName(), " ", "_")
	for _, test := range tests.Tests {
		path := filepath.Join(testsName, test.TestName())
		overrideStateComparison(t, test, path, reflect.TypeOf(tests).String())
	}
}

func overrideStateComparison(t *testing.T, test *StartNewRunnerDutySpecTest, name string, testType string) {
	var runner ssv.Runner
	switch test.Runner.(type) {
	case *ssv.AttesterRunner:
		runner = &ssv.AttesterRunner{}
	case *ssv.AggregatorRunner:
		runner = &ssv.AggregatorRunner{}
	case *ssv.ProposerRunner:
		runner = &ssv.ProposerRunner{}
	case *ssv.SyncCommitteeRunner:
		runner = &ssv.SyncCommitteeRunner{}
	case *ssv.SyncCommitteeAggregatorRunner:
		runner = &ssv.SyncCommitteeAggregatorRunner{}
	case *ssv.ValidatorRegistrationRunner:
		runner = &ssv.ValidatorRegistrationRunner{}
	case *ssv.VoluntaryExitRunner:
		runner = &ssv.VoluntaryExitRunner{}
	default:
		t.Fatalf("unknown runner type")
	}
	basedir, err := os.Getwd()
	require.NoError(t, err)
	runner, err = comparable.UnmarshalStateComparison(basedir, name, testType, runner)
	require.NoError(t, err)

	// override
	test.PostDutyRunnerState = runner

	root, err := runner.GetRoot()
	require.NoError(t, err)

	test.PostDutyRunnerStateRoot = hex.EncodeToString(root[:])
}
