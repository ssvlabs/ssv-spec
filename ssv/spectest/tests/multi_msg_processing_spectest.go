package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type MultiMsgProcessingSpecTest struct {
	Name  string
	Tests []*MsgProcessingSpecTest
}

func (tests *MultiMsgProcessingSpecTest) TestName() string {
	return tests.Name
}

func (tests *MultiMsgProcessingSpecTest) Run(t *testing.T) {
	tests.overrideStateComparison(t)

	for _, test := range tests.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.RunAsMultiTest(t)
		})
	}
}

func (tests *MultiMsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate", "state_comparison", reflect.TypeOf(tests).String(), fmt.Sprintf("%s.json", tests.TestName()))
	byteValue, err := os.ReadFile(path)
	require.NoError(t, err)

	toDecode := make([]ssv.Runner, len(tests.Tests))
	for i, test := range tests.Tests {
		var r ssv.Runner
		switch test.Runner.(type) {
		case *ssv.AttesterRunner:
			r = &ssv.AttesterRunner{}
		case *ssv.AggregatorRunner:
			r = &ssv.AggregatorRunner{}
		case *ssv.ProposerRunner:
			r = &ssv.ProposerRunner{}
		case *ssv.SyncCommitteeRunner:
			r = &ssv.SyncCommitteeRunner{}
		case *ssv.SyncCommitteeAggregatorRunner:
			r = &ssv.SyncCommitteeAggregatorRunner{}
		case *ssv.ValidatorRegistrationRunner:
			r = &ssv.ValidatorRegistrationRunner{}
		default:
			t.Fatalf("unknown runner type")
		}
		toDecode[i] = r
	}
	require.NoError(t, json.Unmarshal(byteValue, &toDecode))

	// override
	for i, test := range tests.Tests {
		test.PostDutyRunnerState = toDecode[i]

		r, err := toDecode[i].GetRoot()
		require.NoError(t, err)

		// backwards compatability test, hard coded post root must be equal to the one loaded from file
		if len(test.PostDutyRunnerStateRoot) > 0 {
			require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(r[:]))
		}

		test.PostDutyRunnerStateRoot = hex.EncodeToString(r[:])
	}
}

func (tests *MultiMsgProcessingSpecTest) GetPostState() (interface{}, error) {
	ret := make([]ssv.Runner, len(tests.Tests))
	for i, test := range tests.Tests {
		_, err := test.runPreTesting()
		if err != nil && len(test.ExpectedError) == 0 {
			return nil, err
		}
		ret[i] = test.Runner
	}
	return ret, nil
}
