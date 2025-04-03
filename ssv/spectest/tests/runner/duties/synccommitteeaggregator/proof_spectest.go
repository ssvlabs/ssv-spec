package synccommitteeaggregator

import (
	"encoding/hex"
	"os"
	"reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type SyncCommitteeAggregatorProofSpecTest struct {
	Name                    string
	Messages                []*types.SignedSSVMessage
	PostDutyRunnerStateRoot string
	PostDutyRunnerState     string
	ProofRootsMap           map[string]bool // if true then root returned from beacon node will be an aggregator
	ExpectedError           string
}

func (test *SyncCommitteeAggregatorProofSpecTest) TestName() string {
	return test.Name
}

func (test *SyncCommitteeAggregatorProofSpecTest) Run(t *testing.T) {
	test.overrideStateComparison(t)
	r, lastErr := test.runPreTesting()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	// post root
	postRoot, err := r.GetBaseRunner().State.GetRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(postRoot[:]))
}

func (test *SyncCommitteeAggregatorProofSpecTest) runPreTesting() (ssv.Runner, error) {
	ks := testingutils.Testing4SharesSet()
	committeeMember := testingutils.TestingCommitteeMember(ks)
	v := testingutils.BaseValidator(testingutils.KeySetForCommitteeMember(committeeMember))
	r := v.DutyRunners[types.RoleSyncCommitteeContribution]
	r.GetBeaconNode().(*testingutils.TestingBeaconNode).SetSyncCommitteeAggregatorRootHexes(test.ProofRootsMap)
	v.Beacon = r.GetBeaconNode()

	lastErr := v.StartDuty(&testingutils.TestingSyncCommitteeContributionDuty)
	for _, msg := range test.Messages {
		err := v.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}
	return r, lastErr
}

func (test *SyncCommitteeAggregatorProofSpecTest) GetPostState() (interface{}, error) {
	runner, err := test.runPreTesting()
	return runner.GetBaseRunner().State, err
}

func (test *SyncCommitteeAggregatorProofSpecTest) overrideStateComparison(t *testing.T) {
	// override state comparison
	basedir, err := os.Getwd()
	require.NoError(t, err)
	postState, err := comparable.UnmarshalStateComparison(basedir, test.Name,
		reflect.TypeOf(test).String(), &ssv.State{})
	require.NoError(t, err)

	r, err2 := postState.GetRoot()
	require.NoError(t, err2)

	test.PostDutyRunnerStateRoot = hex.EncodeToString(r[:])
}
