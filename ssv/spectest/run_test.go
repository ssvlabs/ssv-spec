package spectest

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	tests2 "github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/committee"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/partialsigcontainer"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/newduty"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/runner/duties/synccommitteeaggregator"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

func TestAll(t *testing.T) {
	t.Parallel()
	wait := sync.WaitGroup{}
	for _, testF := range AllTests {
		wait.Add(1)
		go func(f tests2.TestF) {
			test := f()
			t.Run(test.TestName(), func(t *testing.T) {
				test.Run(t)
			})
			wait.Done()
		}(testF)
	}
	wait.Wait()
}

func TestJson(t *testing.T) {
	t.Parallel()

	basedir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	path := filepath.Join(basedir, "generate", "tests.json.gz")
	untypedTests := map[string]interface{}{}
	// Open the gzip file
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open gzip file: %v", err)
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Read the decompressed data
	decompressedData, err := io.ReadAll(gzipReader)
	if err != nil {
		t.Fatalf("failed to read decompressed data: %v", err)
	}

	if err := json.Unmarshal(decompressedData, &untypedTests); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	fmt.Printf("running %d tests\n", len(untypedTests))

	var wg sync.WaitGroup
	for name, test := range untypedTests {
		wg.Add(1)
		go func(name string, test interface{}) {
			defer wg.Done()
			parseAndTest(t, name, test)
		}(name, test)
	}

	wg.Wait()
}

// parseAndTest will parse and test the spec test.
func parseAndTest(t *testing.T, name string, test interface{}) {
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
		case reflect.TypeOf(&partialsigcontainer.PartialSigContainerTest{}).String():
			byts, err := json.Marshal(test)
			require.NoError(t, err)
			typedTest := &partialsigcontainer.PartialSigContainerTest{}
			require.NoError(t, json.Unmarshal(byts, &typedTest))

			typedTest.Run(t)
		case reflect.TypeOf(&committee.CommitteeSpecTest{}).String():
			typedTest := committeeSpecTestFromMap(t, test.(map[string]interface{}))
			typedTest.Run(t)
		case reflect.TypeOf(&committee.MultiCommitteeSpecTest{}).String():
			subtests := test.(map[string]interface{})["Tests"].([]interface{})
			typedTests := make([]*committee.CommitteeSpecTest, 0)
			for _, subtest := range subtests {
				typedTests = append(typedTests, committeeSpecTestFromMap(t, subtest.(map[string]interface{})))
			}

			typedTest := &committee.MultiCommitteeSpecTest{
				Name:  test.(map[string]interface{})["Name"].(string),
				Tests: typedTests,
			}

			typedTest.Run(t)
		default:
			panic("unsupported test type " + testType)
		}
	})
}

func newRunnerDutySpecTestFromMap(t *testing.T, m map[string]interface{}) *newduty.StartNewRunnerDutySpecTest {
	runnerMap := m["Runner"].(map[string]interface{})
	baseRunnerMap := runnerMap["BaseRunner"].(map[string]interface{})

	var testDuty types.Duty
	if _, ok := m["CommitteeDuty"]; ok {
		byts, err := json.Marshal(m["CommitteeDuty"])
		if err != nil {
			panic("cant marshal committee duty")
		}
		committeeDuty := &types.CommitteeDuty{}
		err = json.Unmarshal(byts, committeeDuty)
		if err != nil {
			panic("cant unmarshal committee duty")
		}
		testDuty = committeeDuty
	} else if _, ok := m["ValidatorDuty"]; ok {
		byts, err := json.Marshal(m["ValidatorDuty"])
		if err != nil {
			panic("cant marshal beacon duty")
		}
		duty := &types.ValidatorDuty{}
		err = json.Unmarshal(byts, duty)
		if err != nil {
			panic("cant unmarshal beacon duty")
		}
		testDuty = duty
	} else {
		panic("no beacon or committee duty")
	}

	outputMsgs := make([]*types.PartialSignatureMessages, 0)
	for _, msg := range m["OutputMessages"].([]interface{}) {
		byts, _ := json.Marshal(msg)
		typedMsg := &types.PartialSignatureMessages{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		outputMsgs = append(outputMsgs, typedMsg)
	}

	shareInstance := &types.Share{}
	for _, share := range baseRunnerMap["Share"].(map[string]interface{}) {
		shareBytes, err := json.Marshal(share)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(shareBytes, shareInstance)
		if err != nil {
			panic(err)
		}
	}

	ks := testingutils.KeySetForShare(shareInstance)

	runner := fixRunnerForRun(t, runnerMap, ks)

	return &newduty.StartNewRunnerDutySpecTest{
		Name:                    m["Name"].(string),
		Duty:                    testDuty,
		Runner:                  runner,
		Threshold:               ks.Threshold,
		PostDutyRunnerStateRoot: m["PostDutyRunnerStateRoot"].(string),
		ExpectedError:           m["ExpectedError"].(string),
		OutputMessages:          outputMsgs,
	}
}

func msgProcessingSpecTestFromMap(t *testing.T, m map[string]interface{}) *tests2.MsgProcessingSpecTest {
	runnerMap := m["Runner"].(map[string]interface{})
	baseRunnerMap := runnerMap["BaseRunner"].(map[string]interface{})

	var testDuty types.Duty
	if _, ok := m["CommitteeDuty"]; ok {
		byts, err := json.Marshal(m["CommitteeDuty"])
		if err != nil {
			panic("cant marshal committee duty")
		}
		committeeDuty := &types.CommitteeDuty{}
		err = json.Unmarshal(byts, committeeDuty)
		if err != nil {
			panic("cant unmarshal committee duty")
		}
		testDuty = committeeDuty
	} else if _, ok := m["ValidatorDuty"]; ok {
		byts, err := json.Marshal(m["ValidatorDuty"])
		if err != nil {
			panic("cant marshal beacon duty")
		}
		duty := &types.ValidatorDuty{}
		err = json.Unmarshal(byts, duty)
		if err != nil {
			panic("cant unmarshal beacon duty")
		}
		testDuty = duty
	} else {
		panic("no beacon or committee duty")
	}

	msgs := make([]*types.SignedSSVMessage, 0)
	for _, msg := range m["Messages"].([]interface{}) {
		byts, _ := json.Marshal(msg)
		typedMsg := &types.SignedSSVMessage{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		msgs = append(msgs, typedMsg)
	}

	outputMsgs := make([]*types.PartialSignatureMessages, 0)
	require.NotNilf(t, m["OutputMessages"], "OutputMessages can't be nil")
	for _, msg := range m["OutputMessages"].([]interface{}) {
		byts, _ := json.Marshal(msg)
		typedMsg := &types.PartialSignatureMessages{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		outputMsgs = append(outputMsgs, typedMsg)
	}

	beaconBroadcastedRoots := make([]string, 0)
	if m["BeaconBroadcastedRoots"] != nil {
		for _, r := range m["BeaconBroadcastedRoots"].([]interface{}) {
			beaconBroadcastedRoots = append(beaconBroadcastedRoots, r.(string))
		}
	}

	shareInstance := &types.Share{}
	for _, share := range baseRunnerMap["Share"].(map[string]interface{}) {
		shareBytes, err := json.Marshal(share)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(shareBytes, shareInstance)
		if err != nil {
			panic(err)
		}
	}

	ks := testingutils.KeySetForShare(shareInstance)

	// runner
	runner := fixRunnerForRun(t, runnerMap, ks)

	return &tests2.MsgProcessingSpecTest{
		Name:                    m["Name"].(string),
		Duty:                    testDuty,
		Runner:                  runner,
		Messages:                msgs,
		DecidedSlashable:        m["DecidedSlashable"].(bool),
		PostDutyRunnerStateRoot: m["PostDutyRunnerStateRoot"].(string),
		DontStartDuty:           m["DontStartDuty"].(bool),
		ExpectedError:           m["ExpectedError"].(string),
		OutputMessages:          outputMsgs,
		BeaconBroadcastedRoots:  beaconBroadcastedRoots,
	}
}

func committeeSpecTestFromMap(t *testing.T, m map[string]interface{}) *committee.CommitteeSpecTest {
	committeeMap := m["Committee"].(map[string]interface{})

	inputs := make([]interface{}, 0)
	for _, input := range m["Input"].([]interface{}) {
		byts, err := json.Marshal(input)
		if err != nil {
			panic(err)
		}

		var getDecoder = func() *json.Decoder {
			decoder := json.NewDecoder(strings.NewReader(string(byts)))
			decoder.DisallowUnknownFields()
			return decoder
		}

		committeeDuty := &types.CommitteeDuty{}
		err = getDecoder().Decode(&committeeDuty)
		if err == nil {
			inputs = append(inputs, committeeDuty)
			continue
		}

		duty := &types.ValidatorDuty{}
		err = getDecoder().Decode(&duty)
		if err == nil {
			inputs = append(inputs, duty)
			continue
		}

		msg := &types.SignedSSVMessage{}
		err = getDecoder().Decode(&msg)
		if err == nil {
			inputs = append(inputs, msg)
			continue
		}

		panic(fmt.Sprintf("Unsupported input: %T\n", input))
	}

	outputMsgs := make([]*types.PartialSignatureMessages, 0)
	require.NotNilf(t, m["OutputMessages"], "OutputMessages can't be nil")
	for _, msg := range m["OutputMessages"].([]interface{}) {
		byts, _ := json.Marshal(msg)
		typedMsg := &types.PartialSignatureMessages{}
		require.NoError(t, json.Unmarshal(byts, typedMsg))
		outputMsgs = append(outputMsgs, typedMsg)
	}

	beaconBroadcastedRoots := make([]string, 0)
	if m["BeaconBroadcastedRoots"] != nil {
		for _, r := range m["BeaconBroadcastedRoots"].([]interface{}) {
			beaconBroadcastedRoots = append(beaconBroadcastedRoots, r.(string))
		}
	}

	c := fixCommitteeForRun(t, committeeMap)

	return &committee.CommitteeSpecTest{
		Name:                   m["Name"].(string),
		Committee:              c,
		Input:                  inputs,
		PostDutyCommitteeRoot:  m["PostDutyCommitteeRoot"].(string),
		OutputMessages:         outputMsgs,
		BeaconBroadcastedRoots: beaconBroadcastedRoots,
		ExpectedError:          m["ExpectedError"].(string),
	}
}

func fixCommitteeForRun(t *testing.T, committeeMap map[string]interface{}) *ssv.Committee {

	byts, _ := json.Marshal(committeeMap)
	c := &ssv.Committee{}
	require.NoError(t, json.Unmarshal(byts, c))

	c.CreateRunnerFn = func(shareMap map[phase0.ValidatorIndex]*types.Share) *ssv.CommitteeRunner {
		return testingutils.CommitteeRunnerWithShareMap(shareMap).(*ssv.CommitteeRunner)
	}

	c.SignatureVerifier = testingutils.NewTestingVerifier()

	for slot := range c.Runners {

		var shareInstance *types.Share
		for _, share := range c.Runners[slot].BaseRunner.Share {
			shareInstance = share
			break
		}

		fixedRunner := fixRunnerForRun(t, committeeMap["Runners"].(map[string]interface{})[fmt.Sprintf("%v", slot)].(map[string]interface{}), testingutils.KeySetForShare(shareInstance))
		c.Runners[slot] = fixedRunner.(*ssv.CommitteeRunner)
	}

	return c
}

func fixRunnerForRun(t *testing.T, runnerMap map[string]interface{}, ks *testingutils.TestKeySet) ssv.Runner {
	baseRunnerMap := runnerMap["BaseRunner"].(map[string]interface{})

	base := &ssv.BaseRunner{}
	byts, _ := json.Marshal(baseRunnerMap)
	require.NoError(t, json.Unmarshal(byts, &base))

	ret := baseRunnerForRole(base.RunnerRoleType, base, ks)

	if ret.GetBaseRunner().QBFTController != nil {
		ret.GetBaseRunner().QBFTController = fixControllerForRun(t, ret, ret.GetBaseRunner().QBFTController, ks)
		if ret.GetBaseRunner().State != nil {
			if ret.GetBaseRunner().State.RunningInstance != nil {
				committeeMember := testingutils.TestingCommitteeMember(ks)
				ret.GetBaseRunner().State.RunningInstance = fixInstanceForRun(t, ret.GetBaseRunner().State.RunningInstance, ret.GetBaseRunner().QBFTController, committeeMember)
			}
		}
	}

	return ret
}

func fixControllerForRun(t *testing.T, runner ssv.Runner, contr *qbft.Controller, ks *testingutils.TestKeySet) *qbft.Controller {
	config := testingutils.TestingConfig(ks)
	config.ValueCheckF = runner.GetValCheckF()
	newContr := qbft.NewController(
		contr.Identifier,
		contr.CommitteeMember,
		config,
	)
	newContr.Height = contr.Height
	newContr.StoredInstances = contr.StoredInstances

	for i, inst := range newContr.StoredInstances {
		if inst == nil {
			continue
		}
		committeeMember := testingutils.TestingCommitteeMember(ks)
		newContr.StoredInstances[i] = fixInstanceForRun(t, inst, newContr, committeeMember)
	}
	return newContr
}

func fixInstanceForRun(t *testing.T, inst *qbft.Instance, contr *qbft.Controller, share *types.CommitteeMember) *qbft.Instance {
	newInst := qbft.NewInstance(
		contr.GetConfig(),
		share,
		contr.Identifier,
		contr.Height)

	newInst.State.DecidedValue = inst.State.DecidedValue
	newInst.State.Decided = inst.State.Decided
	newInst.State.CommitteeMember = inst.State.CommitteeMember
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
	newInst.StartValue = inst.StartValue
	return newInst
}

func baseRunnerForRole(role types.RunnerRole, base *ssv.BaseRunner, ks *testingutils.TestKeySet) ssv.Runner {
	switch role {
	case types.RoleCommittee:
		ret := testingutils.CommitteeRunner(ks)
		ret.(*ssv.CommitteeRunner).BaseRunner = base
		return ret
	case types.RoleAggregator:
		ret := testingutils.AggregatorRunner(ks)
		ret.(*ssv.AggregatorRunner).BaseRunner = base
		return ret
	case types.RoleProposer:
		ret := testingutils.ProposerRunner(ks)
		ret.(*ssv.ProposerRunner).BaseRunner = base
		return ret
	case types.RoleSyncCommitteeContribution:
		ret := testingutils.SyncCommitteeContributionRunner(ks)
		ret.(*ssv.SyncCommitteeAggregatorRunner).BaseRunner = base
		return ret
	case types.RoleValidatorRegistration:
		ret := testingutils.ValidatorRegistrationRunner(ks)
		ret.(*ssv.ValidatorRegistrationRunner).BaseRunner = base
		return ret
	case types.RoleVoluntaryExit:
		ret := testingutils.VoluntaryExitRunner(ks)
		ret.(*ssv.VoluntaryExitRunner).BaseRunner = base
		return ret
	case testingutils.UnknownDutyType:
		ret := testingutils.UnknownDutyTypeRunner(ks)
		ret.(*ssv.CommitteeRunner).BaseRunner = base
		return ret
	default:
		panic("unknown beacon role")
	}
}
