package tests

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/types"
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
	testsName := strings.ReplaceAll(tests.TestName(), " ", "_")
	for _, test := range tests.Tests {
		path := filepath.Join(testsName, test.TestName())
		overrideStateComparison(t, test, path, reflect.TypeOf(tests).String())
	}
}

func (tests *MultiMsgProcessingSpecTest) GetPostState() (interface{}, error) {
	ret := make(map[string]types.Root, len(tests.Tests))
	for _, test := range tests.Tests {
		_, _, err := test.runPreTesting()
		if err != nil && test.ExpectedError != err.Error() {
			return nil, err
		}
		ret[test.Name] = test.Runner
	}
	return ret, nil
}
