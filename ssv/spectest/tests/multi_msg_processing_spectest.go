package tests

import (
    "path/filepath"
    "reflect"
    "strings"
    "testing"

    "github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
    "github.com/ssvlabs/ssv-spec/types"
    "github.com/ssvlabs/ssv-spec/types/spectest/tests/errcodes"
    "github.com/ssvlabs/ssv-spec/types/testingutils"
    "github.com/pkg/errors"
)

type MultiMsgProcessingSpecTest struct {
	Name          string
	Type          string
	Documentation string
	Tests         []*MsgProcessingSpecTest
	PrivateKeys   *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
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
        if err != nil && errcodes.FromError(err) != errcodes.FromError(errors.New(test.ExpectedError)) {
            return nil, err
        }
        ret[test.Name] = test.Runner
    }
	return ret, nil
}

func NewMultiMsgProcessingSpecTest(name, documentation string, tests []*MsgProcessingSpecTest, ks *testingutils.TestKeySet) *MultiMsgProcessingSpecTest {
	return &MultiMsgProcessingSpecTest{
		Name:          name,
		Type:          testdoc.MultiMsgProcessingSpecTestType,
		Documentation: documentation,
		Tests:         tests,
		PrivateKeys:   testingutils.BuildPrivateKeyInfo(ks),
	}
}
