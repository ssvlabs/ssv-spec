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

	"github.com/ssvlabs/ssv-spec/types/spectest/tests/aggregatorcommitteeconsensusdata"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/beacon"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/encryption"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/partialsigmessage"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerconsensusdata"
	consensusdataproposer "github.com/ssvlabs/ssv-spec/types/spectest/tests/proposerconsensusdata/proposer"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/share"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/signedssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssvmsg"
	"github.com/ssvlabs/ssv-spec/types/spectest/tests/ssz"
	comparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"
)

func TestAll(t *testing.T) {
	// Run reads the state_comparison fixtures, so stale ones fail TestAll too.
	comparable.LogStaleFixturesHintOnFailure(t)
	for _, test := range AllTests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.Run(t)
		})
	}
}

func TestJson(t *testing.T) {
	comparable.LogStaleFixturesHintOnFailure(t)
	basedir, _ := os.Getwd()
	specTestsDir, err := comparable.SpecTestsDirFrom(basedir)
	if err != nil {
		t.Fatalf("Failed to resolve spec-tests dir: %v", err)
	}
	path := filepath.Join(specTestsDir, "tests.json")
	// Raw, so each test decodes straight into its type; a generic decode rounds uint64s through float64
	untypedTests := map[string]json.RawMessage{}
	byteValue, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if os.Getenv("CI") != "" {
				t.Fatalf("missing %s in CI; generate it with `make generate-jsons`", path)
			}
			t.Skipf("missing %s; generate it with `make generate-jsons`", path)
		}
		t.Fatalf("failed to read %s: %v", path, err)
	}

	if err := json.Unmarshal(byteValue, &untypedTests); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	fmt.Printf("running %d tests\n", len(untypedTests))
	for name, test := range untypedTests {
		var header struct{ Name string }
		require.NoError(t, json.Unmarshal(test, &header))
		t.Run(header.Name, func(t *testing.T) {
			testType := strings.Split(name, "_")[0]
			typedTest := newJSONTest(testType)
			if typedTest == nil {
				t.Fatalf("unsupported test type %s", testType)
			}
			require.NoError(t, json.Unmarshal(test, typedTest))
			typedTest.Run(t)
		})
	}
}

// newJSONTest returns an empty test of the type a tests.json key starts with, or nil if unsupported.
func newJSONTest(testType string) SpecTest {
	for _, test := range []SpecTest{
		&ssz.SSZSpecTest{},
		&consensusdataproposer.ProposerSpecTest{},
		&proposerconsensusdata.EncodingTest{},
		&partialsigmessage.EncodingTest{},
		&share.EncodingTest{},
		&ssvmsg.EncodingTest{},
		&encryption.EncryptionSpecTest{},
		&beacon.DepositDataSpecTest{},
		&signedssvmsg.EncodingTest{},
		&signedssvmsg.SignedSSVMessageTest{},
		&proposerconsensusdata.ProposerConsensusDataTest{},
		&partialsigmessage.MsgSpecTest{},
		&committeemember.CommitteeMemberTest{},
		&ssvmsg.SSVMessageTest{},
		&duty.DutySpecTest{},
		&beaconvote.EncodingTest{},
		&maxmsgsize.StructureSizeTest{},
		&aggregatorcommitteeconsensusdata.AggregatorCommitteeConsensusDataTest{},
		&aggregatorcommitteeconsensusdata.EncodingTest{},
		&beaconvote.ValidationTest{},
		&committeemember.ValidationTest{},
		&share.ValidationTest{},
		&duty.ValidationTest{},
	} {
		if reflect.TypeOf(test).String() == testType {
			return test
		}
	}
	return nil
}
