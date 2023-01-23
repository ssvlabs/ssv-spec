package spectest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/MatheusFranco99/ssv-spec-AleaBFT/qbft"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	tests2 "github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/messages"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/runner/duties/newduty"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/runner/duties/synccommitteeaggregator"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv/spectest/tests/valcheck"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types/testingutils"
	"github.com/stretchr/testify/require"
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
			case reflect.TypeOf(&tests2.MsgProcessingSpecTest{}).String():
				typedTest := msgProcessingSpecTestFromMap(t, test.(map[string]interface{}))
				typedTest.Run(t)
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

				typedTest.Run(t)
			case reflect.TypeOf(&messages.MsgSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &messages.MsgSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				typedTest.Run(t)
			case reflect.TypeOf(&valcheck.SpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &valcheck.SpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				typedTest.Run(t)
			case reflect.TypeOf(&valcheck.MultiSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &valcheck.MultiSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				typedTest.Run(t)
			case reflect.TypeOf(&synccommitteeaggregator.SyncCommitteeAggregatorProofSpecTest{}).String():
				byts, err := json.Marshal(test)
				require.NoError(t, err)
				typedTest := &synccommitteeaggregator.SyncCommitteeAggregatorProofSpecTest{}
				require.NoError(t, json.Unmarshal(byts, &typedTest))

				typedTest.Run(t)
			case reflect.TypeOf(&newduty.MultiStartNewRunnerDutySpecTest{}).String():
				subtests := test.(map[string]interface{})["Tests"].([]interface{})
				typedTests := make([]*newduty.StartNewRunnerDutySpecTest, 0)
				for _, subtest := range subtests {
					typedTests = append(typedTests, newRunnerDutySpecTestFromMap(t, subtest.(map[string]interface{})))
				}

				typedTest := &newduty.MultiStartNewRunnerDutySpecTest{
					Name:  test.(map[string]interface{})["Name"].(string),
					Tests: typedTests,
				}

				typedTest.Run(t)
			default:
				panic("unsupported test type " + testType)
			}
		})
	}
}

func newRunnerDutySpecTestFromMap(t *testing.T, m map[string]interface{}) *newduty.StartNewRunnerDutySpecTest {
	runnerMap := m["Runner"].(map[string]interface{})
	baseRunnerMap := runnerMap["BaseRunner"].(map[string]interface{})

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

	ks := testingutils.KeySetForShare(&types.Share{Quorum: uint64(baseRunnerMap["Share"].(map[string]interface{})["Quorum"].(float64))})

	runner := fixRunnerForRun(t, runnerMap, ks)

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
	baseRunnerMap := runnerMap["BaseRunner"].(map[string]interface{})

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
	require.NotNilf(t, m["OutputMessages"], "OutputMessages can't be nil")
	for _, msg := range m["OutputMessages"].([]interface{}) {
		byts, _ = json.Marshal(msg)
		typedMsg := &ssv.SignedPartialSignatureMessage{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		outputMsgs = append(outputMsgs, typedMsg)
	}

	beaconBroadcastedRoots := make([]string, 0)
	if m["BeaconBroadcastedRoots"] != nil {
		for _, r := range m["BeaconBroadcastedRoots"].([]interface{}) {
			beaconBroadcastedRoots = append(beaconBroadcastedRoots, r.(string))
		}
	}

	ks := testingutils.KeySetForShare(&types.Share{Quorum: uint64(baseRunnerMap["Share"].(map[string]interface{})["Quorum"].(float64))})

	// runner
	runner := fixRunnerForRun(t, runnerMap, ks)

	return &tests2.MsgProcessingSpecTest{
		Name:                    m["Name"].(string),
		Duty:                    duty,
		Runner:                  runner,
		Messages:                msgs,
		PostDutyRunnerStateRoot: m["PostDutyRunnerStateRoot"].(string),
		DontStartDuty:           m["DontStartDuty"].(bool),
		ExpectedError:           m["ExpectedError"].(string),
		OutputMessages:          outputMsgs,
		BeaconBroadcastedRoots:  beaconBroadcastedRoots,
	}
}

func fixRunnerForRun(t *testing.T, runnerMap map[string]interface{}, ks *testingutils.TestKeySet) ssv.Runner {
	baseRunnerMap := runnerMap["BaseRunner"].(map[string]interface{})

	base := &ssv.BaseRunner{}
	byts, _ := json.Marshal(baseRunnerMap)
	require.NoError(t, json.Unmarshal(byts, &base))

	ret := baseRunnerForRole(base.BeaconRoleType, base, ks)

	// specific for blinded block
	if blindedBlocks, ok := runnerMap["ProducesBlindedBlocks"]; ok {
		ret.(*ssv.ProposerRunner).ProducesBlindedBlocks = blindedBlocks.(bool)
	}

	if ret.GetBaseRunner().QBFTController != nil {
		ret.GetBaseRunner().QBFTController = fixControllerForRun(t, ret, ret.GetBaseRunner().QBFTController, ks)
		if ret.GetBaseRunner().State != nil {
			if ret.GetBaseRunner().State.RunningInstance != nil {
				ret.GetBaseRunner().State.RunningInstance = fixInstanceForRun(t, ret.GetBaseRunner().State.RunningInstance, ret.GetBaseRunner().QBFTController, ret.GetBaseRunner().Share)
			}
		}
	}

	return ret
}

func fixControllerForRun(t *testing.T, runner ssv.Runner, contr *qbft.Controller, ks *testingutils.TestKeySet) *qbft.Controller {
	config := testingutils.TestingConfig(ks)
	newContr := qbft.NewController(
		contr.Identifier,
		contr.Share,
		testingutils.TestingConfig(ks).Domain,
		config,
	)
	newContr.Height = contr.Height
	newContr.Domain = contr.Domain
	newContr.StoredInstances = contr.StoredInstances

	for i, inst := range newContr.StoredInstances {
		if inst == nil {
			continue
		}
		newContr.StoredInstances[i] = fixInstanceForRun(t, inst, newContr, runner.GetBaseRunner().Share)
	}
	return newContr
}

func fixInstanceForRun(t *testing.T, inst *qbft.Instance, contr *qbft.Controller, share *types.Share) *qbft.Instance {
	newInst := qbft.NewInstance(
		contr.GetConfig(),
		share,
		contr.Identifier,
		contr.Height)

	newInst.State.DecidedValue = inst.State.DecidedValue
	newInst.State.Decided = inst.State.Decided
	newInst.State.Share = inst.State.Share
	newInst.State.Round = inst.State.Round
	newInst.State.Height = inst.State.Height
	newInst.State.ProposalAcceptedForCurrentRound = inst.State.ProposalAcceptedForCurrentRound
	newInst.State.ID = inst.State.ID
	newInst.State.LastPreparedValue = inst.State.LastPreparedValue
	newInst.State.LastPreparedRound = inst.State.LastPreparedRound
	newInst.State.ProposeContainer = inst.State.ProposeContainer
	newInst.State.PrepareContainer = inst.State.PrepareContainer
	newInst.State.CommitContainer = inst.State.CommitContainer
	newInst.State.RoundChangeContainer = inst.State.RoundChangeContainer
	return newInst
}

func baseRunnerForRole(role types.BeaconRole, base *ssv.BaseRunner, ks *testingutils.TestKeySet) ssv.Runner {
	switch role {
	case types.BNRoleAttester:
		ret := testingutils.AttesterRunner(ks)
		ret.(*ssv.AttesterRunner).BaseRunner = base
		return ret
	case types.BNRoleAggregator:
		ret := testingutils.AggregatorRunner(ks)
		ret.(*ssv.AggregatorRunner).BaseRunner = base
		return ret
	case types.BNRoleProposer:
		ret := testingutils.ProposerRunner(ks)
		ret.(*ssv.ProposerRunner).BaseRunner = base
		return ret
	case types.BNRoleSyncCommittee:
		ret := testingutils.SyncCommitteeRunner(ks)
		ret.(*ssv.SyncCommitteeRunner).BaseRunner = base
		return ret
	case types.BNRoleSyncCommitteeContribution:
		ret := testingutils.SyncCommitteeContributionRunner(ks)
		ret.(*ssv.SyncCommitteeAggregatorRunner).BaseRunner = base
		return ret
	case types.BNRoleValidatorRegistration:
		ret := testingutils.ValidatorRegistrationRunner(ks)
		ret.(*ssv.ValidatorRegistrationRunner).BaseRunner = base
		return ret
	case testingutils.UnknownDutyType:
		ret := testingutils.UnknownDutyTypeRunner(ks)
		ret.(*ssv.AttesterRunner).BaseRunner = base
		return ret
	default:
		panic("unknown beacon role")
	}
}
