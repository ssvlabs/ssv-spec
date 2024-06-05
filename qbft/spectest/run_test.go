package spectest

import (
	"encoding/json"
	"github.com/ssvlabs/ssv-spec/qbft/spectest/tests/timeout"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	tests2 "github.com/ssvlabs/ssv-spec/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	t.Parallel()
	for _, testF := range AllTests {
		test := testF()
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}

func TestJson(t *testing.T) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate", "tests.json")
	untypedTests := map[string]interface{}{}
	byteValue, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(byteValue, &untypedTests); err != nil {
		panic(err.Error())
	}

	tests := make(map[string]tests2.SpecTest)
	for name, test := range untypedTests {
		testName := test.(map[string]interface{})["Name"].(string)
		t.Run(testName, func(t *testing.T) {
			testType := strings.Split(name, "_")[0]
			switch testType {
			case reflect.TypeOf(&tests2.MsgProcessingSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.MsgProcessingSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				// a little trick we do to instantiate all the internal instance params
				preByts, _ := typedTest.Pre.Encode()
				pre := qbft.NewInstance(
					testingutils.TestingConfig(testingutils.KeySetForCommitteeMember(typedTest.Pre.State.CommitteeMember)),
					typedTest.Pre.State.CommitteeMember,
					typedTest.Pre.State.ID,
					typedTest.Pre.State.Height,
				)
				err = pre.Decode(preByts)
				require.NoError(t, err)
				typedTest.Pre = pre

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&tests2.MsgSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.MsgSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&tests2.ControllerSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.ControllerSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&tests2.CreateMsgSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.CreateMsgSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&tests2.RoundRobinSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &tests2.RoundRobinSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)
			case reflect.TypeOf(&timeout.SpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &timeout.SpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				// a little trick we do to instantiate all the internal instance params
				preByts, _ := typedTest.Pre.Encode()
				pre := qbft.NewInstance(
					testingutils.TestingConfig(testingutils.KeySetForCommitteeMember(typedTest.Pre.State.CommitteeMember)),
					typedTest.Pre.State.CommitteeMember,
					typedTest.Pre.State.ID,
					typedTest.Pre.State.Height,
				)
				err = pre.Decode(preByts)
				require.NoError(t, err)
				typedTest.Pre = pre

				tests[testName] = typedTest
				typedTest.Run(t)
			default:
				panic("unsupported test type " + testType)
			}
		})
	}
}
