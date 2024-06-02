package spectest

import (
	"encoding/json"
	"fmt"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/operator"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beaconvote"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/duty"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/consensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/consensusdata/proposer"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssz"
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
	path := filepath.Join(basedir, "generate", "tests.json")
	untypedTests := map[string]interface{}{}
	byteValue, err := os.ReadFile(path)
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
			case reflect.TypeOf(&ssz.SSZSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &ssz.SSZSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&consensusdataproposer.ProposerSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &consensusdataproposer.ProposerSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
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
			case reflect.TypeOf(&signedssvmsg.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &signedssvmsg.EncodingTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&signedssvmsg.SignedSSVMessageTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &signedssvmsg.SignedSSVMessageTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&consensusdata.ConsensusDataTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &consensusdata.ConsensusDataTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&partialsigmessage.MsgSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &partialsigmessage.MsgSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&share.ShareTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &share.ShareTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&operator.OperatorTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &operator.OperatorTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&ssvmsg.SSVMessageTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &ssvmsg.SSVMessageTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&duty.DutySpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &duty.DutySpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			case reflect.TypeOf(&beaconvote.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &beaconvote.EncodingTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			default:
				t.Fatalf("unknown test")
			}
		})
	}
}
