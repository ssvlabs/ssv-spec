package spectest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/committeemember"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/maxmsgsize"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beaconvote"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/duty"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssz"
	validatorconsensusdata "github.com/ssvlabs/ssv-spec/types/spectest/tests/validatorconsensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/validatorconsensusdata/proposer"
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
			case reflect.TypeOf(&validatorconsensusdata.EncodingTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &validatorconsensusdata.EncodingTest{}
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
			case reflect.TypeOf(&validatorconsensusdata.ValidatorConsensusDataTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &validatorconsensusdata.ValidatorConsensusDataTest{}
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
			case reflect.TypeOf(&committeemember.CommitteeMemberTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &committeemember.CommitteeMemberTest{}
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
			case reflect.TypeOf(&maxmsgsize.StructureSizeTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &maxmsgsize.StructureSizeTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))
				typedTest.Run(t)
			default:
				t.Fatalf("unknown test")
			}
		})
	}
}
