package spectest

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/beacon"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/consensusdata"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/encryption"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/share"
	"github.com/bloxapp/ssv-spec/types/spectest/tests/ssvmsg"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	for _, test := range AllTests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}

func TestJson(t *testing.T) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate")
	fileName := "tests.json"
	untypedTests := map[string]interface{}{}
	byteValue, err := os.ReadFile(path + "/" + fileName)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(byteValue, &untypedTests); err != nil {
		panic(err.Error())
	}

	fmt.Printf("running %d tests\n", len(untypedTests))
	for name, test := range untypedTests {
		testName := test.(map[string]interface{})["Name"].(string)
		t.Run(testName, func(t *testing.T) {
			testType := strings.Split(name, "_")[0]
			switch testType {
			case reflect.TypeOf(&consensusdata.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &consensusdata.EncodingTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&partialsigmessage.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &partialsigmessage.EncodingTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&share.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &share.EncodingTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&ssvmsg.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &ssvmsg.EncodingTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&encryption.EncryptionSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &encryption.EncryptionSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&beacon.DepositDataSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &beacon.DepositDataSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			default:
				t.Fatalf("unknown test")
			}
		})
	}
}
