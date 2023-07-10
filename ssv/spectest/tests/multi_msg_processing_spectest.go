package tests

import (
	"github.com/bloxapp/ssv-spec/ssv"
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
			test.RunAsPartOfMultiTest(t)
		})
	}
}

// overrideStateComparison overrides the post state comparison for all tests in the multi test
func (tests *MultiMsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	for _, test := range tests.Tests {
		path := filepath.Join(tests.TestName(), test.TestName())
		overrideStateComparison(t, test, path, reflect.TypeOf(tests).String())
	}
}

func (tests *MultiMsgProcessingSpecTest) GetPostState() (interface{}, error) {
	ret := make([]ssv.Runner, len(tests.Tests))
	for i, test := range tests.Tests {
		_, err := test.runPreTesting()
		if err != nil && test.ExpectedError != err.Error() {
			return nil, err
		}
		ret[i] = test.Runner
	}
	return ret, nil
}
