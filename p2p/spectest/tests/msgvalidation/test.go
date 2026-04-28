package msgvalidation

import (
	"context"
	"crypto/rsa"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/p2p/spectest/testdoc"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type RunnerType string

const (
	RunnerTypeCommittee           RunnerType = "committee"
	RunnerTypeAggregatorCommittee RunnerType = "aggregator_committee"
	RunnerTypeProposer            RunnerType = "proposer"
)

type DutyKind string

const (
	DutyKindAttester            DutyKind = "attester"
	DutyKindProposer            DutyKind = "proposer"
	DutyKindAggregatorCommittee DutyKind = "aggregator_committee"
)

type RunnerSetupOpType string

const (
	RunnerSetupOpStartDuty           RunnerSetupOpType = "start_duty"
	RunnerSetupOpStartInstance       RunnerSetupOpType = "start_instance"
	RunnerSetupOpSetControllerHeight RunnerSetupOpType = "set_controller_height"
	RunnerSetupOpSetState            RunnerSetupOpType = "set_state"
)

type DutyPreset struct {
	Kind    DutyKind         `json:"Kind"`
	Version spec.DataVersion `json:"Version,omitempty"`
}

type RunnerStateOptions struct {
	DecidedValue           []byte `json:"DecidedValue,omitempty"`
	CreateRunningInstance  bool   `json:"CreateRunningInstance,omitempty"`
	RunningInstanceDecided bool   `json:"RunningInstanceDecided,omitempty"`
	StoreRunningInstance   bool   `json:"StoreRunningInstance,omitempty"`
	SetControllerHeight    bool   `json:"SetControllerHeight,omitempty"`
}

type RunnerSetupOp struct {
	Type         RunnerSetupOpType   `json:"Type"`
	Duty         *DutyPreset         `json:"Duty,omitempty"`
	Height       uint64              `json:"Height,omitempty"`
	Value        []byte              `json:"Value,omitempty"`
	StateOptions *RunnerStateOptions `json:"StateOptions,omitempty"`
}

type RunnerPreset struct {
	Type  RunnerType      `json:"Type"`
	Setup []RunnerSetupOp `json:"Setup,omitempty"`
}

type MsgValidationSpecTest struct {
	Name             string                  `json:"Name"`
	Type             string                  `json:"Type,omitempty"`
	Documentation    string                  `json:"Documentation"`
	Runner           RunnerPreset            `json:"Runner"`
	Message          []byte                  `json:"Message"`
	ExpectedDecision pubsub.ValidationResult `json:"ExpectedDecision"`
	// PrivateKeys are emitted for generated JSON consumers; local test execution
	// reconstructs the standard test key set in buildRunner.
	PrivateKeys *testingutils.PrivateKeyInfo `json:"PrivateKeys,omitempty"`
}

func (test *MsgValidationSpecTest) TestName() string {
	return "msg validation " + test.Name
}

func (test *MsgValidationSpecTest) Run(t *testing.T) {
	runner := buildRunner(t, test.Runner)
	result := validation.MsgValidation(runner)(context.Background(), "", PubsubMessageFromBytes(test.Message))
	require.Equal(t, test.ExpectedDecision, result)
}

func (test *MsgValidationSpecTest) GetPostState() (interface{}, error) {
	return nil, nil
}

func NewMsgValidationSpecTest(
	name string,
	documentation string,
	runner RunnerPreset,
	message []byte,
	expectedDecision pubsub.ValidationResult,
	ks *testingutils.TestKeySet,
) tests.SpecTest {
	return &MsgValidationSpecTest{
		Name:             name,
		Type:             testdoc.MsgValidationSpecTestType,
		Documentation:    documentation,
		Runner:           runner,
		Message:          message,
		ExpectedDecision: expectedDecision,
		PrivateKeys:      testingutils.BuildPrivateKeyInfo(ks),
	}
}

func NewRunnerPreset(runnerType RunnerType, setup ...RunnerSetupOp) RunnerPreset {
	return RunnerPreset{
		Type:  runnerType,
		Setup: setup,
	}
}

func AttesterDuty(version spec.DataVersion) DutyPreset {
	return DutyPreset{Kind: DutyKindAttester, Version: version}
}

func ProposerDuty(version spec.DataVersion) DutyPreset {
	return DutyPreset{Kind: DutyKindProposer, Version: version}
}

func AggregatorCommitteeDuty(version spec.DataVersion) DutyPreset {
	return DutyPreset{Kind: DutyKindAggregatorCommittee, Version: version}
}

func StartDutyOp(duty DutyPreset) RunnerSetupOp {
	return RunnerSetupOp{
		Type: RunnerSetupOpStartDuty,
		Duty: &duty,
	}
}

func StartInstanceOp(height qbft.Height, value []byte) RunnerSetupOp {
	return RunnerSetupOp{
		Type:   RunnerSetupOpStartInstance,
		Height: uint64(height),
		Value:  value,
	}
}

func SetControllerHeightOp(height qbft.Height) RunnerSetupOp {
	return RunnerSetupOp{
		Type:   RunnerSetupOpSetControllerHeight,
		Height: uint64(height),
	}
}

func StateWithoutDecidedValueOp(duty DutyPreset) RunnerSetupOp {
	return RunnerSetupOp{
		Type: RunnerSetupOpSetState,
		Duty: &duty,
	}
}

func StateWithDecidedValueOp(duty DutyPreset, decidedValue []byte) RunnerSetupOp {
	return RunnerSetupOp{
		Type: RunnerSetupOpSetState,
		Duty: &duty,
		StateOptions: &RunnerStateOptions{
			DecidedValue: decidedValue,
		},
	}
}

func StateWithRunningInstanceOp(duty DutyPreset, decidedValue []byte, decided bool) RunnerSetupOp {
	return RunnerSetupOp{
		Type: RunnerSetupOpSetState,
		Duty: &duty,
		StateOptions: &RunnerStateOptions{
			DecidedValue:           decidedValue,
			CreateRunningInstance:  true,
			RunningInstanceDecided: decided,
		},
	}
}

func MarkDecidedStateOp(duty DutyPreset, decidedValue []byte) RunnerSetupOp {
	return RunnerSetupOp{
		Type: RunnerSetupOpSetState,
		Duty: &duty,
		StateOptions: &RunnerStateOptions{
			DecidedValue:           decidedValue,
			CreateRunningInstance:  true,
			RunningInstanceDecided: true,
			StoreRunningInstance:   true,
			SetControllerHeight:    true,
		},
	}
}

func buildRunner(t *testing.T, preset RunnerPreset) ssv.Runner {
	t.Helper()

	ks := testingutils.Testing4SharesSet()

	var runner ssv.Runner
	switch preset.Type {
	case RunnerTypeCommittee:
		runner = testingutils.CommitteeRunner(ks)
	case RunnerTypeAggregatorCommittee:
		runner = testingutils.AggregatorCommitteeRunner(ks)
	case RunnerTypeProposer:
		runner = testingutils.ProposerRunner(ks)
	default:
		t.Fatalf("unknown runner type: %s", preset.Type)
	}

	for _, op := range preset.Setup {
		applyRunnerSetupOp(t, runner, op)
	}

	return runner
}

func applyRunnerSetupOp(t *testing.T, runner ssv.Runner, op RunnerSetupOp) {
	t.Helper()

	switch op.Type {
	case RunnerSetupOpStartDuty:
		duty := mustBuildDuty(t, op.Duty)
		err := runner.StartNewDuty(duty, runner.GetBaseRunner().QBFTController.CommitteeMember.GetQuorum())
		require.NoError(t, err)
	case RunnerSetupOpStartInstance:
		err := runner.GetBaseRunner().QBFTController.StartNewInstance(qbft.Height(op.Height), op.Value)
		require.NoError(t, err)
	case RunnerSetupOpSetControllerHeight:
		runner.GetBaseRunner().QBFTController.Height = qbft.Height(op.Height)
	case RunnerSetupOpSetState:
		duty := mustBuildDuty(t, op.Duty)
		opts := op.StateOptions
		if opts == nil {
			opts = &RunnerStateOptions{}
		}

		base := runner.GetBaseRunner()
		base.State = ssv.NewRunnerState(base.QBFTController.CommitteeMember.GetQuorum(), duty)
		if len(opts.DecidedValue) > 0 {
			base.State.DecidedValue = opts.DecidedValue
		}
		if opts.CreateRunningInstance {
			base.State.RunningInstance = qbft.NewInstance(
				base.QBFTController.GetConfig(),
				base.QBFTController.CommitteeMember,
				base.QBFTController.Identifier,
				qbft.Height(duty.DutySlot()),
				base.QBFTController.OperatorSigner,
			)
			if opts.RunningInstanceDecided {
				base.State.RunningInstance.State.Decided = true
				base.State.RunningInstance.State.DecidedValue = opts.DecidedValue
			}
			if opts.StoreRunningInstance {
				base.QBFTController.StoredInstances = append(
					base.QBFTController.StoredInstances,
					base.State.RunningInstance,
				)
			}
		}
		if opts.SetControllerHeight {
			base.QBFTController.Height = qbft.Height(duty.DutySlot())
		}
	default:
		t.Fatalf("unknown runner setup op: %s", op.Type)
	}
}

func mustBuildDuty(t *testing.T, preset *DutyPreset) types.Duty {
	t.Helper()

	if preset == nil {
		t.Fatal("missing duty preset")
	}

	switch preset.Kind {
	case DutyKindAttester:
		return testingutils.TestingAttesterDuty(preset.Version)
	case DutyKindProposer:
		return testingutils.TestingProposerDutyV(preset.Version)
	case DutyKindAggregatorCommittee:
		return testingutils.TestingAggregatorAndSyncCommitteeContributorDuties(preset.Version)
	default:
		t.Fatalf("unknown duty kind: %s", preset.Kind)
		return nil
	}
}

func PubsubMessageFromBytes(data []byte) *pubsub.Message {
	return &pubsub.Message{Message: &pubsubpb.Message{Data: data}}
}

func EncodeSignedSSVMessage(signedMsg *types.SignedSSVMessage) []byte {
	data, err := signedMsg.Encode()
	if err != nil {
		panic(err)
	}
	return data
}

func MutateSignedSSVMessage(signedMsg *types.SignedSSVMessage, mutate func(*types.SignedSSVMessage)) *types.SignedSSVMessage {
	ret := signedMsg.DeepCopy()
	mutate(ret)
	return ret
}

func MakeUnsupportedSSVMessage(ks *testingutils.TestKeySet) *types.SignedSSVMessage {
	duty := testingutils.TestingAggregatorAndSyncCommitteeContributorDuties(spec.DataVersionPhase0)
	keySetMap := map[phase0.ValidatorIndex]*testingutils.TestKeySet{
		testingutils.TestingValidatorIndex: ks,
	}
	partialSigMsg := testingutils.PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, 1)
	ssvMsg := testingutils.SSVMsgAggregatorCommittee(ks, nil, partialSigMsg)
	ssvMsg.MsgType = types.DKGMsgType
	sig, err := types.SignSSVMessage(ks.OperatorKeys[1], ssvMsg)
	if err != nil {
		panic(err)
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{sig},
		SSVMessage:  ssvMsg,
	}
}

func MakeExistingInstanceConsensusMessage(ks *testingutils.TestKeySet, height qbft.Height) *types.SignedSSVMessage {
	return testingutils.TestingProposalMessageWithIdentifierAndFullData(
		ks.OperatorKeys[1],
		1,
		testingutils.CommitteeMsgID(ks),
		testingutils.TestBeaconVoteByts,
		height,
	)
}

func MakeFutureDecidedConsensusMessage(ks *testingutils.TestKeySet) *types.SignedSSVMessage {
	return testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
		[]*rsa.PrivateKey{ks.OperatorKeys[1], ks.OperatorKeys[2], ks.OperatorKeys[3]},
		[]types.OperatorID{1, 2, 3},
		qbft.Height(testingutils.TestingDutySlot),
		testingutils.CommitteeMsgID(ks),
		testingutils.TestBeaconVoteByts,
	)
}

func resignSignedMessage(operatorID types.OperatorID, sk *rsa.PrivateKey, ssvMessage *types.SSVMessage) *types.SignedSSVMessage {
	sig, err := types.SignSSVMessage(sk, ssvMessage)
	if err != nil {
		panic(err)
	}

	return &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{operatorID},
		Signatures:  [][]byte{sig},
		SSVMessage:  ssvMessage,
	}
}

func MutatePartialSignatureSignedMessage(
	signedMsg *types.SignedSSVMessage,
	operatorID types.OperatorID,
	sk *rsa.PrivateKey,
	mutate func(*types.PartialSignatureMessages),
) *types.SignedSSVMessage {
	partialSigMsg := &types.PartialSignatureMessages{}
	if err := partialSigMsg.Decode(signedMsg.SSVMessage.Data); err != nil {
		panic(err)
	}

	mutate(partialSigMsg)

	data, err := partialSigMsg.Encode()
	if err != nil {
		panic(err)
	}

	ssvMessage := &types.SSVMessage{
		MsgType: signedMsg.SSVMessage.MsgType,
		MsgID:   signedMsg.SSVMessage.MsgID,
		Data:    data,
	}

	return resignSignedMessage(operatorID, sk, ssvMessage)
}
