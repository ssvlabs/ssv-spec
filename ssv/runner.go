package ssv

import (
	"bytes"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type Getters interface {
	GetBaseRunner() *BaseRunner
	GetBeaconNode() BeaconNode
	GetValCheckF() qbft.ProposedValueCheckF
	GetSigner() types.KeyManager
	GetNetwork() Network
}

type Runner interface {
	types.Encoder
	types.Root
	Getters

	StartNewDuty(duty *types.Duty) error
	HasRunningDuty() bool
	ProcessPreConsensus(signedMsg *SignedPartialSignature) error
	ProcessConsensus(msg *types.Message) error
	ProcessPostConsensus(signedMsg *SignedPartialSignature) error

	expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error)
	executeDuty(duty *types.Duty) error
}

type BaseRunner struct {
	State          *State
	Share          *types.Share
	QBFTController *qbft.Controller
	BeaconNetwork  types.BeaconNetwork
	BeaconRoleType types.BeaconRole
}

func (b *BaseRunner) baseStartNewDuty(runner Runner, duty *types.Duty) error {
	if err := b.canStartNewDuty(); err != nil {
		return err
	}
	b.State = NewRunnerState(b.Share.Quorum, duty)
	return runner.executeDuty(duty)
}

func (b *BaseRunner) canStartNewDuty() error {
	if b.State == nil {
		return nil
	}

	// check if instance running first as we can't start new duty if it does
	if b.State.RunningInstance != nil {
		// check consensus decided
		if decided, _ := b.State.RunningInstance.IsDecided(); !decided {
			return errors.New("consensus on duty is running")
		}
	}
	return nil
}

func (b *BaseRunner) basePreConsensusMsgProcessing(runner Runner, signedMsg *SignedPartialSignature) (bool, [][]byte, error) {
	if err := b.validatePreConsensusMsg(runner, signedMsg); err != nil {
		return false, nil, errors.Wrap(err, "invalid pre-consensus message")
	}

	roots := make([][]byte, 0)
	anyQuorum := false
	for _, msg := range signedMsg.Message.Messages {
		prevQuorum := b.State.PreConsensusContainer.HasQuorum(msg.SigningRoot)

		if err := b.State.PreConsensusContainer.AddSignature(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not add partial randao signature")
		}

		if prevQuorum {
			continue
		}

		quorum := b.State.PreConsensusContainer.HasQuorum(msg.SigningRoot)
		if quorum {
			roots = append(roots, msg.SigningRoot)
			anyQuorum = true
		}
	}

	return anyQuorum, roots, nil
}

func (b *BaseRunner) baseConsensusMsgProcessing(runner Runner, msg *types.Message) (decided bool, decidedValue *types.ConsensusData, err error) {
	if err := b.validateConsensusMsg(); err != nil {
		return false, nil, errors.Wrap(err, "invalid consensus message")
	}

	prevDecided, _ := b.State.RunningInstance.IsDecided()

	decidedMsg, err := b.QBFTController.ProcessMsg(msg)
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to process consensus msg")
	}

	if decideCorrectly, err := b.didDecideCorrectly(prevDecided, decidedMsg); !decideCorrectly {
		return false, nil, err
	}

	decidedValue = &types.ConsensusData{}
	if err := decidedValue.UnmarshalSSZ(decidedMsg.InputSource); err != nil {
		return true, nil, errors.Wrap(err, "failed to parse decided value to ConsensusData")
	}

	if err := b.validateDecidedConsensusData(runner, decidedValue); err != nil {
		return true, nil, errors.Wrap(err, "decided ConsensusData invalid")
	}

	runner.GetBaseRunner().State.DecidedValue = decidedValue

	return true, decidedValue, nil
}

func (b *BaseRunner) basePostConsensusMsgProcessing(signedMsg *SignedPartialSignature) (bool, [][]byte, error) {
	if err := b.validatePostConsensusMsg(signedMsg); err != nil {
		return false, nil, errors.Wrap(err, "invalid post-consensus message")
	}

	roots := make([][]byte, 0)
	anyQuorum := false
	for _, msg := range signedMsg.Message.Messages {
		prevQuorum := b.State.PostConsensusContainer.HasQuorum(msg.SigningRoot)

		if err := b.State.PostConsensusContainer.AddSignature(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not add partial post consensus signature")
		}

		if prevQuorum {
			continue
		}

		quorum := b.State.PostConsensusContainer.HasQuorum(msg.SigningRoot)

		if quorum {
			roots = append(roots, msg.SigningRoot)
			anyQuorum = true
		}
	}

	return anyQuorum, roots, nil
}

func (b *BaseRunner) didDecideCorrectly(prevDecided bool, decidedMsg *qbft.SignedMessage) (bool, error) {
	decided := decidedMsg != nil
	decidedRunningInstance := decided && decidedMsg.Message.Height == b.State.RunningInstance.GetHeight()

	if !decided {
		return false, nil
	}
	if !decidedRunningInstance {
		return false, errors.New("decided wrong instance")
	}
	// verify we decided running instance only, if not we do not proceed
	if prevDecided {
		return false, nil
	}

	return true, nil
}

func (b *BaseRunner) validatePreConsensusMsg(runner Runner, signedMsg *SignedPartialSignature) error {
	if !b.HashRunningDuty() {
		return errors.New("no running duty")
	}

	if err := b.validatePartialSigMsg(signedMsg, b.State.StartingDuty.Slot); err != nil {
		return err
	}

	roots, domain, err := runner.expectedPreConsensusRootsAndDomain()
	if err != nil {
		return err
	}

	return b.verifyExpectedRoot(runner, signedMsg, roots, domain)
}

func (b *BaseRunner) validateConsensusMsg() error {
	if !b.HashRunningDuty() {
		return errors.New("no running duty")
	}
	return nil
}

func (b *BaseRunner) validatePostConsensusMsg(msg *SignedPartialSignature) error {
	if !b.HashRunningDuty() {
		return errors.New("no running duty")
	}

	if err := b.validatePartialSigMsg(msg, b.State.DecidedValue.Duty.Slot); err != nil {
		return errors.Wrap(err, "post consensus msg invalid")
	}

	return nil
}

func (b *BaseRunner) decide(runner Runner, input *types.ConsensusData) error {
	source, err := input.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not encode ConsensusData")
	}
	root, err := input.HashTreeRoot()
	if err != nil {
		return nil
	}
	inputData := &qbft.Data{
		Root:   root,
		Source: source,
	}

	if err := runner.GetValCheckF()(inputData); err != nil {
		return errors.Wrap(err, "input data invalid")
	}

	if err := runner.GetBaseRunner().QBFTController.StartNewInstance(inputData); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := runner.GetBaseRunner().QBFTController.InstanceForHeight(runner.GetBaseRunner().QBFTController.Height)
	if newInstance == nil {
		return errors.New("could not find newly created QBFT instance")
	}

	runner.GetBaseRunner().State.RunningInstance = newInstance
	return nil
}

func (b *BaseRunner) HashRunningDuty() bool {
	if b.State == nil {
		return false
	}
	return !b.State.Finished
}

func (b *BaseRunner) signBeaconObject(
	runner Runner,
	obj ssz.HashRoot,
	slot phase0.Slot,
	domainType phase0.DomainType,
) (*PartialSignature, error) {
	epoch := runner.GetBaseRunner().BeaconNetwork.EstimatedEpochAtSlot(slot)
	domain, err := runner.GetBeaconNode().DomainData(epoch, domainType)
	if err != nil {
		return nil, errors.Wrap(err, "could not get beacon domain")
	}

	sig, r, err := runner.GetSigner().SignBeaconObject(obj, domain, runner.GetBaseRunner().Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign beacon object")
	}

	return &PartialSignature{
		Slot:             slot,
		PartialSignature: sig,
		SigningRoot:      r,
		Signer:           runner.GetBaseRunner().Share.OperatorID,
	}, nil
}

func (b *BaseRunner) validatePartialSigMsg(
	signedMsg *SignedPartialSignature,
	slot phase0.Slot,
) error {
	if err := signedMsg.Validate(); err != nil {
		return errors.Wrap(err, "SignedPartialSignature invalid")
	}

	if err := signedMsg.GetSignature().VerifyByOperators(signedMsg, b.Share.DomainType, types.PartialSignatureType, b.Share.Committee); err != nil {
		return errors.Wrap(err, "failed to verify PartialSignature")
	}

	for _, msg := range signedMsg.Message.Messages {
		if slot != msg.Slot {
			return errors.New("wrong slot")
		}

		if err := b.verifyBeaconPartialSignature(msg); err != nil {
			return errors.Wrap(err, "could not verify Beacon partial Signature")
		}
	}

	return nil
}

func (b *BaseRunner) verifyBeaconPartialSignature(msg *PartialSignature) error {
	signer := msg.Signer
	signature := msg.PartialSignature
	root := msg.SigningRoot

	for _, n := range b.Share.Committee {
		if n.GetID() == signer {
			pk := &bls.PublicKey{}
			if err := pk.Deserialize(n.GetPublicKey()); err != nil {
				return errors.Wrap(err, "could not deserialized pk")
			}
			sig := &bls.Sign{}
			if err := sig.Deserialize(signature); err != nil {
				return errors.Wrap(err, "could not deserialized Signature")
			}

			// verify
			if !sig.VerifyByte(pk, root) {
				return errors.New("wrong signature")
			}
			return nil
		}
	}
	return errors.New("unknown signer")
}

func (b *BaseRunner) validateDecidedConsensusData(runner Runner, val *types.ConsensusData) error {
	source, err := val.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not encode decided value")
	}
	root, err := val.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not get decided root")
	}
	if err := runner.GetValCheckF()(&qbft.Data{
		Root:   root,
		Source: source,
	}); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	return nil
}

func (b *BaseRunner) verifyExpectedRoot(runner Runner, signedMsg *SignedPartialSignature, expectedRootObjs []ssz.HashRoot, domain phase0.DomainType) error {
	if len(expectedRootObjs) != len(signedMsg.Message.Messages) {
		return errors.New("wrong expected roots count")
	}
	for i, msg := range signedMsg.Message.Messages {
		epoch := b.BeaconNetwork.EstimatedEpochAtSlot(b.State.StartingDuty.Slot)
		d, err := runner.GetBeaconNode().DomainData(epoch, domain)
		if err != nil {
			return errors.Wrap(err, "could not get pre consensus root domain")
		}

		r, err := types.ComputeETHSigningRoot(expectedRootObjs[i], d)
		if err != nil {
			return errors.Wrap(err, "could not compute ETH signing root")
		}
		if !bytes.Equal(r[:], msg.SigningRoot) {
			return errors.New("wrong pre consensus signing root")
		}
	}
	return nil
}

func (b *BaseRunner) signPostConsensusMsg(runner Runner, msg *PartialSignatures) (*SignedPartialSignature, error) {
	signature, err := runner.GetSigner().SignRoot(msg, types.PartialSignatureType, b.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign PartialSignature for PostConsensusContainer")
	}

	return &SignedPartialSignature{
		Message:   *msg,
		Signature: signature,
		Signer:    b.Share.OperatorID,
	}, nil
}
