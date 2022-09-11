package spectest

import (
	"encoding/json"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	tests2 "github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/messages"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/runner/duties/newduty"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
	byteValue, err := ioutil.ReadFile(path + "/" + fileName)
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(byteValue, &untypedTests); err != nil {
		panic(err.Error())
	}

	for name, test := range untypedTests {
		testType := strings.Split(name, "_")[0]
		switch testType {
		case reflect.TypeOf(&tests2.MsgProcessingSpecTest{}).String():
			runMsgProcessingSpecTest(t, msgProcessingSpecTestFromMap(t, test.(map[string]interface{})))
		case reflect.TypeOf(&tests2.MultiMsgProcessingSpecTest{}).String():
			subtests := test.(map[string]interface{})["Tests"].([]interface{})
			typedTests := make([]*tests2.MsgProcessingSpecTest, 0)
			for _, subtest := range subtests {
				typedTests = append(typedTests, msgProcessingSpecTestFromMap(t, subtest.(map[string]interface{})))
			}

			typedTest := &tests2.MultiMsgProcessingSpecTest{
				Name:  test.(map[string]interface{})["Name"].(string),
				Tests: typedTests,
			}

			t.Run(typedTest.TestName(), func(t *testing.T) {
				for _, test := range typedTest.Tests {
					runMsgProcessingSpecTest(t, test)
				}
			})

		case reflect.TypeOf(&messages.MsgSpecTest{}).String():
			byts, err := json.Marshal(test)
			require.NoError(t, err)
			typedTest := &messages.MsgSpecTest{}
			require.NoError(t, json.Unmarshal(byts, &typedTest))

			t.Run(typedTest.TestName(), func(t *testing.T) {
				typedTest.Run(t)
			})
		case reflect.TypeOf(&valcheck.MultiSpecTest{}).String():
			byts, err := json.Marshal(test)
			require.NoError(t, err)
			typedTest := &valcheck.MultiSpecTest{}
			require.NoError(t, json.Unmarshal(byts, &typedTest))

			t.Run(typedTest.TestName(), func(t *testing.T) {
				typedTest.Run(t)
			})
		case reflect.TypeOf(&newduty.MultiStartNewRunnerDutySpecTest{}).String():
			subtests := test.(map[string]interface{})["Tests"].([]interface{})
			typedTests := make([]*newduty.StartNewRunnerDutySpecTest, 0)
			for _, subtest := range subtests {
				typedTests = append(typedTests, startNewRunnerDutySpecTestFromMap(t, subtest.(map[string]interface{})))
			}

			typedTest := &newduty.MultiStartNewRunnerDutySpecTest{
				Name:  test.(map[string]interface{})["Name"].(string),
				Tests: typedTests,
			}

			t.Run(typedTest.TestName(), func(t *testing.T) {
				typedTest.Run(t)
			})
		default:
			panic("unsupported test type " + testType)
		}
	}
}

func startNewRunnerDutySpecTestFromMap(t *testing.T, m map[string]interface{}) *newduty.StartNewRunnerDutySpecTest {
	runnerMap := m["Runner"].(map[string]interface{})

	duty := &types.Duty{}
	byts, _ := json.Marshal(m["Duty"])
	require.NoError(t, json.Unmarshal(byts, duty))

	outputMsgs := make([]*ssv.SignedPartialSignatureMessage, 0)
	for _, msg := range m["OutputMessages"].([]interface{}) {
		byts, _ = json.Marshal(msg)
		typedMsg := &ssv.SignedPartialSignatureMessage{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		outputMsgs = append(outputMsgs, typedMsg)
	}

	ks := testingutils.KeySetForShare(&types.Share{Quorum: uint64(runnerMap["Share"].(map[string]interface{})["Quorum"].(float64))})

	runner := runnerFromRole(types.BeaconRole(runnerMap["BeaconRoleType"].(float64)), ks)

	return &newduty.StartNewRunnerDutySpecTest{
		Name:                    m["Name"].(string),
		Duty:                    duty,
		Runner:                  runner,
		PostDutyRunnerStateRoot: m["PostDutyRunnerStateRoot"].(string),
		ExpectedError:           m["ExpectedError"].(string),
		OutputMessages:          outputMsgs,
	}
}

func msgProcessingSpecTestFromMap(t *testing.T, m map[string]interface{}) *tests2.MsgProcessingSpecTest {
	runnerMap := m["Runner"].(map[string]interface{})

	duty := &types.Duty{}
	byts, _ := json.Marshal(m["Duty"])
	require.NoError(t, json.Unmarshal(byts, duty))

	msgs := make([]*types.SSVMessage, 0)
	for _, msg := range m["Messages"].([]interface{}) {
		byts, _ = json.Marshal(msg)
		typedMsg := &types.SSVMessage{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		msgs = append(msgs, typedMsg)
	}

	outputMsgs := make([]*ssv.SignedPartialSignatureMessage, 0)
	for _, msg := range m["OutputMessages"].([]interface{}) {
		byts, _ = json.Marshal(msg)
		typedMsg := &ssv.SignedPartialSignatureMessage{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		outputMsgs = append(outputMsgs, typedMsg)
	}

	ks := testingutils.KeySetForShare(&types.Share{Quorum: uint64(runnerMap["Share"].(map[string]interface{})["Quorum"].(float64))})

	runner := runnerFromRole(types.BeaconRole(runnerMap["BeaconRoleType"].(float64)), ks)

	return &tests2.MsgProcessingSpecTest{
		Name:                    m["Name"].(string),
		Duty:                    duty,
		Runner:                  runner,
		Messages:                msgs,
		PostDutyRunnerStateRoot: m["PostDutyRunnerStateRoot"].(string),
		DontStartDuty:           m["DontStartDuty"].(bool),
		ExpectedError:           m["ExpectedError"].(string),
		OutputMessages:          outputMsgs,
	}
}

func runnerFromRole(role types.BeaconRole, ks *testingutils.TestKeySet) ssv.Runner {
	switch role {
	case types.BNRoleAttester:
		return testingutils.AttesterRunner(ks)
	case types.BNRoleAggregator:
		return testingutils.AggregatorRunner(ks)
	case types.BNRoleProposer:
		return testingutils.ProposerRunner(ks)
	case types.BNRoleSyncCommittee:
		return testingutils.SyncCommitteeRunner(ks)
	case types.BNRoleSyncCommitteeContribution:
		return testingutils.SyncCommitteeContributionRunner(ks)
	case testingutils.UnknownDutyType:
		return testingutils.UnknownDutyTypeRunner(ks)
	default:
		panic("unknown beacon role")
	}
}

func runMsgProcessingSpecTest(t *testing.T, test *tests2.MsgProcessingSpecTest) {
	// a little trick we do to instantiate all the internal controller params
	byts, err := test.Runner.GetQBFTController().Encode()
	require.NoError(t, err)

	ks := testingutils.KeySetForShare(test.Runner.GetQBFTController().Share)

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

		if test.Runner.GetBaseRunner().State != nil &&
			test.Runner.GetBaseRunner().State.RunningInstance != nil &&
			test.Runner.GetBaseRunner().State.RunningInstance.GetHeight() == fixedInst.GetHeight() {
			test.Runner.GetBaseRunner().State.RunningInstance = fixedInst
		}
	}
	t.Run(test.Name, func(t *testing.T) {
		test.Run(t)
	})
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

func setControllerInRunner(runner ssv.Runner, controller *qbft.Controller) {
	switch runner.GetBaseRunner().BeaconRoleType {
	case types.BNRoleAttester:
		runner.(*ssv.AttesterRunner).BaseRunner.QBFTController = controller
	case types.BNRoleAggregator:
		runner.(*ssv.AggregatorRunner).BaseRunner.QBFTController = controller
	case types.BNRoleProposer:
		runner.(*ssv.ProposerRunner).BaseRunner.QBFTController = controller
	case types.BNRoleSyncCommittee:
		runner.(*ssv.SyncCommitteeRunner).BaseRunner.QBFTController = controller
	case types.BNRoleSyncCommitteeContribution:
		runner.(*ssv.SyncCommitteeAggregatorRunner).BaseRunner.QBFTController = controller
	}
}
