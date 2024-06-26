package tests

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/go-cmp/cmp"
	typescomparable "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type MsgProcessingSpecTest struct {
	Name     string
	Runner   ssv.Runner
	Duty     types.Duty
	Messages []*types.SignedSSVMessage
	// DecidedSlashable makes the decided value slashable. Simulates consensus instances running in parallel.
	DecidedSlashable        bool
	PostDutyRunnerStateRoot string
	PostDutyRunnerState     types.Root `json:"-"` // Field is ignored by encoding/json
	// OutputMessages compares pre/ post signed partial sigs to output. We exclude consensus msgs as it's tested in consensus
	OutputMessages         []*types.PartialSignatureMessages
	BeaconBroadcastedRoots []string
	DontStartDuty          bool // if set to true will not start a duty for the runner
	ExpectedError          string
}

func (test *MsgProcessingSpecTest) TestName() string {
	return test.Name
}

// RunAsPartOfMultiTest runs the test as part of a MultiMsgProcessingSpecTest
func (test *MsgProcessingSpecTest) RunAsPartOfMultiTest(t *testing.T) {
	v, c, lastErr := test.runPreTesting()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	network := &testingutils.TestingNetwork{}
	beaconNetwork := &testingutils.TestingBeaconNode{}
	var committee []*types.Operator
	switch test.Runner.(type) {
	case *ssv.CommitteeRunner:
		var runnerInstance *ssv.CommitteeRunner
		for _, runner := range c.Runners {
			runnerInstance = runner
			break
		}
		network = runnerInstance.GetNetwork().(*testingutils.TestingNetwork)
		beaconNetwork = runnerInstance.GetBeaconNode().(*testingutils.TestingBeaconNode)
		committee = c.CommitteeMember.Committee
	default:
		network = v.Network.(*testingutils.TestingNetwork)
		committee = v.CommitteeMember.Committee
		beaconNetwork = test.Runner.GetBeaconNode().(*testingutils.TestingBeaconNode)
	}

	// test output message
	testingutils.ComparePartialSignatureOutputMessages(t, test.OutputMessages, network.BroadcastedMsgs, committee)

	// test beacon broadcasted msgs
	testingutils.CompareBroadcastedBeaconMsgs(t, test.BeaconBroadcastedRoots, beaconNetwork.BroadcastedRoots)

	// post root
	postRoot, err := test.Runner.GetRoot()
	require.NoError(t, err)

	if test.PostDutyRunnerStateRoot != hex.EncodeToString(postRoot[:]) {
		t.Errorf("post runner state not equal: %v", cmp.Diff(test.Runner, test.PostDutyRunnerState, cmp.Exporter(func(p reflect.Type) bool { return true })))
	}
}

// Run as an individual test
func (test *MsgProcessingSpecTest) Run(t *testing.T) {
	test.overrideStateComparison(t)
	test.RunAsPartOfMultiTest(t)
}

func (test *MsgProcessingSpecTest) runPreTesting() (*ssv.Validator, *ssv.Committee, error) {
	var share *types.Share
	ketSetMap := make(map[phase0.ValidatorIndex]*testingutils.TestKeySet)
	if len(test.Runner.GetBaseRunner().Share) == 0 {
		panic("No share in base runner for tests")
	}
	for _, validatorShare := range test.Runner.GetBaseRunner().Share {
		share = validatorShare
		break
	}
	for valIdx, validatorShare := range test.Runner.GetBaseRunner().Share {
		ketSetMap[valIdx] = testingutils.KeySetForShare(validatorShare)
	}

	var v *ssv.Validator
	var c *ssv.Committee
	var lastErr error
	switch test.Runner.(type) {
	case *ssv.CommitteeRunner:
		c = testingutils.BaseCommitteeWithRunner(ketSetMap, test.Runner.(*ssv.CommitteeRunner))

		if !test.DontStartDuty {
			lastErr = c.StartDuty(test.Duty.(*types.CommitteeDuty))
		} else {
			c.Runners[test.Duty.DutySlot()] = test.Runner.(*ssv.CommitteeRunner)
		}

		for _, msg := range test.Messages {
			err := c.ProcessMessage(msg)
			if err != nil {
				lastErr = err
			}
			if test.DecidedSlashable && IsQBFTProposalMessage(msg) {
				for _, validatorShare := range test.Runner.GetBaseRunner().Share {
					test.Runner.GetSigner().(*testingutils.TestingKeyManager).AddSlashableDataRoot(validatorShare.
						SharePubKey, testingutils.TestingAttestationDataRoot[:])
				}
			}
		}

	default:
		v = testingutils.BaseValidator(testingutils.KeySetForShare(share))
		v.DutyRunners[test.Runner.GetBaseRunner().RunnerRoleType] = test.Runner
		v.Network = test.Runner.GetNetwork()

		if !test.DontStartDuty {
			lastErr = v.StartDuty(test.Duty)
		}
		for _, msg := range test.Messages {
			err := v.ProcessMessage(msg)
			if err != nil {
				lastErr = err
			}
		}
	}

	return v, c, lastErr
}

// IsQBFTProposalMessage checks if the message is a QBFT proposal message
func IsQBFTProposalMessage(msg *types.SignedSSVMessage) bool {
	if msg.SSVMessage.MsgType == types.SSVConsensusMsgType {
		qbftMsg := qbft.Message{}
		err := qbftMsg.Decode(msg.SSVMessage.Data)
		if err != nil {
			panic("could not decode message")
		}
		return qbftMsg.MsgType == qbft.ProposalMsgType
	}
	return false
}

func (test *MsgProcessingSpecTest) overrideStateComparison(t *testing.T) {
	overrideStateComparison(t, test, test.Name, reflect.TypeOf(test).String())
}

func overrideStateComparison(t *testing.T, test *MsgProcessingSpecTest, name string, testType string) {
	var runner ssv.Runner
	switch test.Runner.(type) {
	case *ssv.CommitteeRunner:
		runner = &ssv.CommitteeRunner{}
	case *ssv.AggregatorRunner:
		runner = &ssv.AggregatorRunner{}
	case *ssv.ProposerRunner:
		runner = &ssv.ProposerRunner{}
	case *ssv.SyncCommitteeAggregatorRunner:
		runner = &ssv.SyncCommitteeAggregatorRunner{}
	case *ssv.ValidatorRegistrationRunner:
		runner = &ssv.ValidatorRegistrationRunner{}
	case *ssv.VoluntaryExitRunner:
		runner = &ssv.VoluntaryExitRunner{}
	default:
		t.Fatalf("unknown runner type")
	}
	basedir, err := os.Getwd()
	require.NoError(t, err)
	runner, err = typescomparable.UnmarshalStateComparison(basedir, name, testType, runner)
	require.NoError(t, err)

	// override
	test.PostDutyRunnerState = runner

	root, err := runner.GetRoot()
	require.NoError(t, err)

	test.PostDutyRunnerStateRoot = hex.EncodeToString(root[:])
}

func (test *MsgProcessingSpecTest) GetPostState() (interface{}, error) {
	_, _, lastErr := test.runPreTesting()
	if lastErr != nil && len(test.ExpectedError) == 0 {
		return nil, lastErr
	}

	return test.Runner, nil
}

// Create alias without duty
type MsgProcessingSpecTestAlias struct {
	Name   string
	Runner ssv.Runner
	// No duty from type types.Duty
	Messages                []*types.SignedSSVMessage
	DecidedSlashable        bool
	PostDutyRunnerStateRoot string
	PostDutyRunnerState     types.Root `json:"-"`
	OutputMessages          []*types.PartialSignatureMessages
	BeaconBroadcastedRoots  []string
	DontStartDuty           bool
	ExpectedError           string
	ValidatorDuty           *types.ValidatorDuty `json:"ValidatorDuty,omitempty"`
	CommitteeDuty           *types.CommitteeDuty `json:"CommitteeDuty,omitempty"`
}

func (t *MsgProcessingSpecTest) MarshalJSON() ([]byte, error) {
	alias := &MsgProcessingSpecTestAlias{
		Name:                    t.Name,
		Runner:                  t.Runner,
		Messages:                t.Messages,
		DecidedSlashable:        t.DecidedSlashable,
		PostDutyRunnerStateRoot: t.PostDutyRunnerStateRoot,
		PostDutyRunnerState:     t.PostDutyRunnerState,
		OutputMessages:          t.OutputMessages,
		BeaconBroadcastedRoots:  t.BeaconBroadcastedRoots,
		DontStartDuty:           t.DontStartDuty,
		ExpectedError:           t.ExpectedError,
	}

	if t.Duty != nil {
		if duty, ok := t.Duty.(*types.ValidatorDuty); ok {
			alias.ValidatorDuty = duty
		} else if committeeDuty, ok := t.Duty.(*types.CommitteeDuty); ok {
			alias.CommitteeDuty = committeeDuty
		} else {
			return nil, errors.New("can't marshal StartNewRunnerDutySpecTest because t.Duty isn't ValidatorDuty or CommitteeDuty")
		}
	}
	byts, err := json.Marshal(alias)

	return byts, err
}

func (t *MsgProcessingSpecTest) UnmarshalJSON(data []byte) error {
	aux := &MsgProcessingSpecTestAlias{}

	// Unmarshal the JSON data into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Name = aux.Name
	t.Runner = aux.Runner
	t.DecidedSlashable = aux.DecidedSlashable
	t.Messages = aux.Messages
	t.PostDutyRunnerStateRoot = aux.PostDutyRunnerStateRoot
	t.PostDutyRunnerState = aux.PostDutyRunnerState
	t.OutputMessages = aux.OutputMessages
	t.BeaconBroadcastedRoots = aux.BeaconBroadcastedRoots
	t.DontStartDuty = aux.DontStartDuty
	t.ExpectedError = aux.ExpectedError

	// Determine which type of duty was marshaled
	if aux.ValidatorDuty != nil {
		t.Duty = aux.ValidatorDuty
	} else if aux.CommitteeDuty != nil {
		t.Duty = aux.CommitteeDuty
	}

	return nil
}
