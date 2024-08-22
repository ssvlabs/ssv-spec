package spectest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	tests2 "github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/validation"
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
			case reflect.TypeOf(&validation.MessageValidationTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &validation.MessageValidationTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				tests[testName] = typedTest
				typedTest.Run(t)

			case reflect.TypeOf(&validation.MultiMessageValidationTest{}).String():
				subtests := test.(map[string]interface{})["Tests"].([]interface{})
				typedTests := make([]*validation.MessageValidationTest, 0)
				for _, subtest := range subtests {
					byts, err := json.Marshal(subtest)
					require.NoError(t, err)
					subTypedTest := &validation.MessageValidationTest{}
					require.NoError(t, json.Unmarshal(byts, &subTypedTest))

					typedTests = append(typedTests, subTypedTest)
				}

				typedTest := &validation.MultiMessageValidationTest{
					Name:  test.(map[string]interface{})["Name"].(string),
					Tests: typedTests,
				}

				typedTest.Run(t)
			default:
				panic("unsupported test type " + testType)
			}
		})
	}
}
