package spectest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests/msgvalidation"
)

func TestAll(t *testing.T) {
	t.Parallel()

	for _, testF := range AllTests {
		test := testF()
		t.Run(test.TestName(), func(t *testing.T) {
			t.Parallel()
			test.Run(t)
		})
	}
}

func TestJson(t *testing.T) {
	t.Parallel()

	basedir, err := os.Getwd()
	require.NoError(t, err)

	path := filepath.Join(basedir, "generate", "tests.json")
	byteValue, err := os.ReadFile(path)
	require.NoError(t, err)

	untypedTests := map[string]interface{}{}
	require.NoError(t, json.Unmarshal(byteValue, &untypedTests))

	for name, test := range untypedTests {
		testName := test.(map[string]interface{})["Name"].(string)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			testType := strings.Split(name, "_")[0]
			switch testType {
			case reflect.TypeOf(&msgvalidation.MsgValidationSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)

				typedTest := &msgvalidation.MsgValidationSpecTest{}
				require.NoError(t, json.Unmarshal(byts, typedTest))
				typedTest.Run(t)
			default:
				t.Fatalf("unsupported test type %s", testType)
			}
		})
	}
}
