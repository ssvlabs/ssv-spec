package tests

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type MultiMsgProcessingSpecTest struct {
	Name          string
	Type          string
	Documentation string
	Tests         []*MsgProcessingSpecTest
	PrivateKeys   *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

func (mTest *MultiMsgProcessingSpecTest) TestName() string {
	return mTest.Name
}

func (mTest *MultiMsgProcessingSpecTest) Run(t *testing.T) {
	mTest.overrideStateComparison(t)

	for _, test := range mTest.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.RunAsPartOfMultiTest(t)
		})
	}
}

// overrideStateComparison overrides the post state comparison for all tests in the multi test
func (mTest *MultiMsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	testsName := strings.ReplaceAll(mTest.TestName(), " ", "_")
	for _, test := range mTest.Tests {
		path := filepath.Join(testsName, test.TestName())
		overrideStateComparison(t, test, path, reflect.TypeOf(mTest).String())
	}
}

func (mTest *MultiMsgProcessingSpecTest) GetPostState() (interface{}, error) {
	ret := make(map[string]types.Root, len(mTest.Tests))
	for _, test := range mTest.Tests {
		_, _, err := test.runPreTesting()
		if err != nil && !MatchesErrorCode(test.ExpectedErrorCode, err) {
			return nil, fmt.Errorf(
				"(%s) expected error with code: %d, got error: %s",
				test.TestName(),
				test.ExpectedErrorCode,
				err,
			)
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
