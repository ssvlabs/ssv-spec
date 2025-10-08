package committee

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/go-cmp/cmp"

	"github.com/ssvlabs/ssv-spec/ssv/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type CommitteeSpecTest struct {
	Name                   string
	Type                   string `json:"Type,omitempty"`
	Documentation          string
	Committee              *ssv.Committee
	Input                  []interface{} // Can be a types.Duty or a *types.SignedSSVMessage
	PostDutyCommitteeRoot  string
	PostDutyCommittee      types.Root `json:"-"` // Field is ignored by encoding/json
	OutputMessages         []*types.PartialSignatureMessages
	BeaconBroadcastedRoots []string
	ExpectedErrorCode      int
	PrivateKeys            *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

func (test *CommitteeSpecTest) TestName() string {
	return test.Name
}

// RunAsPartOfMultiTest runs the test as part of a MultiCommitteeSpecTest
func (test *CommitteeSpecTest) RunAsPartOfMultiTest(t *testing.T) {
	lastErr := test.runPreTesting()
	tests.AssertErrorCode(t, test.ExpectedErrorCode, lastErr)

	broadcastedMsgs := make([]*types.SignedSSVMessage, 0)
	broadcastedRoots := make([]phase0.Root, 0)
	for _, runner := range test.Committee.Runners {
		network := runner.GetNetwork().(*testingutils.TestingNetwork)
		beaconNetwork := runner.GetBeaconNode().(*testingutils.TestingBeaconNode)
		broadcastedMsgs = append(broadcastedMsgs, network.BroadcastedMsgs...)
		broadcastedRoots = append(broadcastedRoots, beaconNetwork.BroadcastedRoots...)
	}

	// test output message (in asynchronous order)
	testingutils.ComparePartialSignatureOutputMessagesInAsynchronousOrder(t, test.OutputMessages, broadcastedMsgs, test.Committee.CommitteeMember.Committee)

	// test beacon broadcasted msgs
	testingutils.CompareBroadcastedBeaconMsgs(t, test.BeaconBroadcastedRoots, broadcastedRoots)

	// post root
	postRoot, err := test.Committee.GetRoot()
	require.NoError(t, err)

	if test.PostDutyCommitteeRoot != hex.EncodeToString(postRoot[:]) {
		t.Errorf("post runner state not equal: %v", cmp.Diff(test.Committee, test.PostDutyCommittee, cmp.Exporter(func(p reflect.Type) bool { return true })))
	}
}

// Run as an individual test
func (test *CommitteeSpecTest) Run(t *testing.T) {
	test.overrideStateComparison(t)
	test.RunAsPartOfMultiTest(t)
}

func (test *CommitteeSpecTest) runPreTesting() error {

	var lastErr error

	for _, input := range test.Input {
		var err error
		switch input := input.(type) {
		case types.Duty:
			err = test.Committee.StartDuty(input.(*types.CommitteeDuty))
		case *types.SignedSSVMessage:
			err = test.Committee.ProcessMessage(input)
		default:
			panic("input is neither duty or SignedSSVMessage")
		}
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (test *CommitteeSpecTest) overrideStateComparison(t *testing.T) {
	overrideStateComparison(t, test, test.Name, reflect.TypeOf(test).String())
}

func overrideStateComparison(t *testing.T, test *CommitteeSpecTest, name string, testType string) {
	committee := &ssv.Committee{}
	basedir, err := os.Getwd()
	require.NoError(t, err)
	committee, err = typescomparable.UnmarshalStateComparison(basedir, name, testType, committee)
	require.NoError(t, err)

	// override
	test.PostDutyCommittee = committee

	root, err := committee.GetRoot()
	require.NoError(t, err)

	test.PostDutyCommitteeRoot = hex.EncodeToString(root[:])
}

func (test *CommitteeSpecTest) GetPostState() (interface{}, error) {
	lastErr := test.runPreTesting()
	if lastErr != nil && test.ExpectedErrorCode == 0 {
		return nil, lastErr
	}

	return test.Committee, nil
}

type MultiCommitteeSpecTest struct {
	Name          string
	Type          string
	Documentation string
	Tests         []*CommitteeSpecTest
	PrivateKeys   *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

func (test *MultiCommitteeSpecTest) TestName() string {
	return test.Name
}

func (test *MultiCommitteeSpecTest) Run(t *testing.T) {
	test.overrideStateComparison(t)

	for _, test := range test.Tests {
		t.Run(test.TestName(), func(t *testing.T) {
			test.RunAsPartOfMultiTest(t)
		})
	}
}

// overrideStateComparison overrides the post state comparison for all tests in the multi test
func (test *MultiCommitteeSpecTest) overrideStateComparison(t *testing.T) {
	testsName := strings.ReplaceAll(test.TestName(), " ", "_")
	for _, test := range test.Tests {
		path := filepath.Join(testsName, test.TestName())
		overrideStateComparison(t, test, path, reflect.TypeOf(test).String())
	}
}

func (test *MultiCommitteeSpecTest) GetPostState() (interface{}, error) {
	ret := make(map[string]types.Root, len(test.Tests))
	for _, test := range test.Tests {
		err := test.runPreTesting()
		if err != nil && !tests.MatchesErrorCode(test.ExpectedErrorCode, err) {
			return nil, fmt.Errorf(
				"(%s) expected error with code: %d, got error: %s",
				test.TestName(),
				test.ExpectedErrorCode,
				err,
			)
		}
		ret[test.Name] = test.Committee
	}
	return ret, nil
}

func NewMultiCommitteeSpecTest(name, documentation string, tests []*CommitteeSpecTest, ks *testingutils.TestKeySet) *MultiCommitteeSpecTest {
	return &MultiCommitteeSpecTest{
		Name:          name,
		Type:          testdoc.MultiCommitteeSpecTestType,
		Documentation: documentation,
		Tests:         tests,
		PrivateKeys:   testingutils.BuildPrivateKeyInfo(ks),
	}
}

func NewCommitteeSpecTest(name, documentation string, committee *ssv.Committee, input []interface{}, postDutyCommitteeRoot string, postDutyCommittee types.Root, outputMessages []*types.PartialSignatureMessages, beaconBroadcastedRoots []string, expectedErrorCode int, ks *testingutils.TestKeySet) *CommitteeSpecTest {
	return &CommitteeSpecTest{
		Name:                   name,
		Type:                   testdoc.CommitteeSpecTestType,
		Documentation:          documentation,
		Committee:              committee,
		Input:                  input,
		PostDutyCommitteeRoot:  postDutyCommitteeRoot,
		PostDutyCommittee:      postDutyCommittee,
		OutputMessages:         outputMessages,
		BeaconBroadcastedRoots: beaconBroadcastedRoots,
		ExpectedErrorCode:      expectedErrorCode,
		PrivateKeys:            testingutils.BuildPrivateKeyInfo(ks),
	}
}
