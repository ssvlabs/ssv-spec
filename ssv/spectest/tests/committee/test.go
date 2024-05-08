package committee

import (
	"encoding/hex"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type CommitteeSpecTest struct {
	Name                   string
	Committee              *ssv.Committee
	input                  []interface{} // Can be a types.Duty or a *types.SignedSSVMessage
	PostDutyCommitteeRoot  string
	PostDutyCommittee      types.Root `json:"-"` // Field is ignored by encoding/json
	OutputMessages         []*types.PartialSignatureMessages
	BeaconBroadcastedRoots []string
	ExpectedError          string
}

func (test *CommitteeSpecTest) TestName() string {
	return test.Name
}

// RunAsPartOfMultiTest runs the test as part of a MultiCommitteeSpecTest
func (test *CommitteeSpecTest) RunAsPartOfMultiTest(t *testing.T) {
	lastErr := test.runPreTesting()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	var runnerInstance *ssv.CommitteeRunner
	for _, runner := range test.Committee.Runners {
		runnerInstance = runner
		break
	}
	network := runnerInstance.GetNetwork().(*testingutils.TestingNetwork)
	beaconNetwork := runnerInstance.GetBeaconNode().(*testingutils.TestingBeaconNode)
	committee := test.Committee.Operator.Committee

	// test output message
	testingutils.ComparePartialSignatureOutputMessages(t, test.OutputMessages, network.BroadcastedMsgs, committee)

	// test beacon broadcasted msgs
	testingutils.CompareBroadcastedBeaconMsgs(t, test.BeaconBroadcastedRoots, beaconNetwork.BroadcastedRoots)

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

	for _, input := range test.input {
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
	var committee *ssv.Committee
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
	if lastErr != nil && len(test.ExpectedError) == 0 {
		return nil, lastErr
	}

	return test.Committee, nil
}
