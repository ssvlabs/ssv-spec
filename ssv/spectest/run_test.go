package spectest

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	tests2 "github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
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
	tests := map[string]*tests2.MsgProcessingSpecTest{}
	byteValue, err := ioutil.ReadFile(path + "/" + fileName)
	require.NoError(t, err)

	if err := json.Unmarshal(byteValue, &tests); err != nil {
		require.NoError(t, err)
	}

	for _, test := range tests {

		// a little trick we do to instantiate all the internal controller params
		byts, err := test.Runner.GetQBFTController().Encode()
		require.NoError(t, err)

		ks := keySetForShare(test.Runner.GetQBFTController().Share)

		newContr := qbft.NewController(
			[]byte{1, 2, 3, 4},
			test.Runner.GetQBFTController().Share,
			testingutils.TestingConfig(ks).Domain,
			testingutils.TestingConfig(ks).Signer,
			testingutils.TestingConfig(ks).ValueCheckF,
			testingutils.TestingConfig(ks).Storage,
			testingutils.TestingConfig(ks).Network,
			func(state *qbft.State, round qbft.Round) types.OperatorID {
				return 1
			},
		)
		require.NoError(t, newContr.Decode(byts))
		setControllerInRunner(test.Runner, newContr)

		for idx, i := range test.Runner.GetQBFTController().StoredInstances {
			if i == nil {
				continue
			}
			fixedInst := fixQBFTInstanceForRun(t, i, ks)
			test.Runner.GetQBFTController().StoredInstances[idx] = fixedInst

			if test.Runner.GetState() != nil &&
				test.Runner.GetState().RunningInstance != nil &&
				test.Runner.GetState().RunningInstance.GetHeight() == fixedInst.GetHeight() {
				test.Runner.GetState().RunningInstance = fixedInst
			}
		}
		t.Run(test.Name, func(t *testing.T) {
			test.Run(t)
		})
	}
}

func fixQBFTInstanceForRun(t *testing.T, i *qbft.Instance, ks *testingutils.TestKeySet) *qbft.Instance {
	// a little trick we do to instantiate all the internal instance params
	if i == nil {
		return nil
	}
	byts, _ := i.Encode()
	newInst := qbft.NewInstance(testingutils.TestingConfig(ks), i.State.Share, i.State.ID, qbft.FirstHeight)
	require.NoError(t, newInst.Decode(byts))
	return newInst
}

func keySetForShare(share *types.Share) *testingutils.TestKeySet {
	if share.Quorum == 5 {
		return testingutils.Testing7SharesSet()
	}
	if share.Quorum == 7 {
		return testingutils.Testing10SharesSet()
	}
	if share.Quorum == 9 {
		return testingutils.Testing13SharesSet()
	}
	return testingutils.Testing4SharesSet()
}

func setControllerInRunner(runner ssv.Runner, controller *qbft.Controller) {
	switch runner.GetBeaconRole() {
	case types.BNRoleAttester:
		runner.(*ssv.AttesterRunner).QBFTController = controller
	case types.BNRoleAggregator:
		runner.(*ssv.AggregatorRunner).QBFTController = controller
	case types.BNRoleProposer:
		runner.(*ssv.ProposerRunner).QBFTController = controller
	case types.BNRoleSyncCommittee:
		runner.(*ssv.SyncCommitteeRunner).QBFTController = controller
	case types.BNRoleSyncCommitteeContribution:
		runner.(*ssv.SyncCommitteeAggregatorRunner).QBFTController = controller
	}
}
