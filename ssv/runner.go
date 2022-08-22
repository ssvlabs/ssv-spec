package ssv

import (
	"bytes"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type Getters interface {
	GetNetwork() Network
	GetBeaconNetwork() types.BeaconNetwork
	GetBeaconNode() BeaconNode
	GetBeaconRole() types.BeaconRole
	GetShare() *types.Share
	GetState() *State
	GetValCheckF() qbft.ProposedValueCheckF
	GetQBFTController() *qbft.Controller
	GetSigner() types.KeyManager
}

type Runner interface {
	types.Encoder
	types.Root
	Getters

	StartNewDuty(duty *types.Duty) error
	HasRunningDuty() bool
	ProcessPreConsensus(signedMsg *SignedPartialSignatureMessage) error
	ProcessConsensus(msg *qbft.SignedMessage) error
	ProcessPostConsensus(signedMsg *SignedPartialSignatureMessage) error

	expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, spec.DomainType, error)
}

func basePreConsensusMsgProcessing(runner Runner, signedMsg *SignedPartialSignatureMessage) (bool, [][]byte, error) {
	if err := validatePreConsensusMsg(runner, signedMsg); err != nil {
		return false, nil, errors.Wrap(err, "invalid pre-consensus message")
	}

	roots := make([][]byte, 0)
	anyQuorum := false
	for _, msg := range signedMsg.Message.Messages {
		prevQuorum := runner.GetState().PreConsensusContainer.HasQuorum(msg.SigningRoot)

		if err := runner.GetState().PreConsensusContainer.AddSignature(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not add partial randao signature")
		}

		if prevQuorum {
			continue
		}

		quorum := runner.GetState().PreConsensusContainer.HasQuorum(msg.SigningRoot)
		if quorum {
			roots = append(roots, msg.SigningRoot)
			anyQuorum = true
		}
	}

	return anyQuorum, roots, nil
}

func validatePreConsensusMsg(runner Runner, signedMsg *SignedPartialSignatureMessage) error {
	if err := validatePartialSigMsg(runner, signedMsg, runner.GetState().StartingDuty.Slot); err != nil {
		return err
	}

	roots, domain, err := runner.expectedPreConsensusRootsAndDomain()
	if err != nil {
		return err
	}

	return verifyExpectedRoot(runner, signedMsg, roots, domain)
}

func baseConsensusMsgProcessing(runner Runner, msg *qbft.SignedMessage) (decided bool, decidedValue *types.ConsensusData, err error) {
	decided, decidedValueByts, err := runner.GetQBFTController().ProcessMsg(msg)
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to process consensus msg")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return false, nil, nil
	}

	decidedValue = &types.ConsensusData{}
	if err := decidedValue.Decode(decidedValueByts); err != nil {
		return true, nil, errors.Wrap(err, "failed to parse decided value to ConsensusData")
	}

	if err := validateDecidedConsensusData(runner, decidedValue); err != nil {
		return true, nil, errors.Wrap(err, "decided ConsensusData invalid")
	}

	return true, decidedValue, nil
}

func basePostConsensusMsgProcessing(runner Runner, signedMsg *SignedPartialSignatureMessage) (bool, [][]byte, error) {
	if err := canProcessPostConsensusMsg(runner, signedMsg); err != nil {
		return false, nil, errors.Wrap(err, "can't process post consensus message")
	}

	roots := make([][]byte, 0)
	anyQuorum := false
	for _, msg := range signedMsg.Message.Messages {
		prevQuorum := runner.GetState().PostConsensusContainer.HasQuorum(msg.SigningRoot)

		if err := runner.GetState().PostConsensusContainer.AddSignature(msg); err != nil {
			return false, nil, errors.Wrap(err, "could not add partial post consensus signature")
		}

		if prevQuorum {
			continue
		}

		quorum := runner.GetState().PostConsensusContainer.HasQuorum(msg.SigningRoot)

		if quorum {
			roots = append(roots, msg.SigningRoot)
			anyQuorum = true
		}
	}

	return anyQuorum, roots, nil
}

func canStartNewDuty(runner Runner, duty *types.Duty) error {
	if runner.GetState() == nil {
		return nil
	}

	// check if instance running first as we can't start new duty if it does
	if runner.GetState().RunningInstance != nil {
		// check consensus decided
		if decided, _ := runner.GetState().RunningInstance.IsDecided(); !decided {
			return errors.New("consensus on duty is running")
		}
	}
	return nil
}

func decide(runner Runner, input *types.ConsensusData) error {
	byts, err := input.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode ConsensusData")
	}

	if err := runner.GetValCheckF()(byts); err != nil {
		return errors.Wrap(err, "input data invalid")
	}

	if err := runner.GetQBFTController().StartNewInstance(byts); err != nil {
		return errors.Wrap(err, "could not start new QBFT instance")
	}
	newInstance := runner.GetQBFTController().InstanceForHeight(runner.GetQBFTController().Height)
	if newInstance == nil {
		return errors.New("could not find newly created QBFT instance")
	}

	runner.GetState().RunningInstance = newInstance
	return nil
}

func signBeaconObject(
	runner Runner,
	obj ssz.HashRoot,
	slot spec.Slot,
	domainType spec.DomainType,
) (*PartialSignatureMessage, error) {
	epoch := runner.GetBeaconNetwork().EstimatedEpochAtSlot(slot)
	domain, err := runner.GetBeaconNode().DomainData(epoch, domainType)
	if err != nil {
		return nil, errors.Wrap(err, "could not get beacon domain")
	}

	sig, r, err := runner.GetSigner().SignBeaconObject(obj, domain, runner.GetShare().SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign beacon object")
	}

	return &PartialSignatureMessage{
		Slot:             slot,
		PartialSignature: sig,
		SigningRoot:      r,
		Signer:           runner.GetShare().OperatorID,
	}, nil
}

func validatePartialSigMsg(
	runner Runner,
	signedMsg *SignedPartialSignatureMessage,
	slot spec.Slot,
) error {
	if err := signedMsg.Validate(); err != nil {
		return errors.Wrap(err, "SignedPartialSignatureMessage invalid")
	}

	if err := signedMsg.GetSignature().VerifyByOperators(signedMsg, runner.GetShare().DomainType, types.PartialSignatureType, runner.GetShare().Committee); err != nil {
		return errors.Wrap(err, "failed to verify PartialSignature")
	}

	for _, msg := range signedMsg.Message.Messages {
		if slot != msg.Slot {
			return errors.New("wrong slot")
		}

		if err := verifyBeaconPartialSignature(runner, msg); err != nil {
			return errors.Wrap(err, "could not verify Beacon partial Signature")
		}
	}

	return nil
}

func verifyBeaconPartialSignature(runner Runner, msg *PartialSignatureMessage) error {
	signer := msg.Signer
	signature := msg.PartialSignature
	root := msg.SigningRoot

	for _, n := range runner.GetShare().Committee {
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
	return errors.New("Beacon partial Signature Signer not found")
}

func validateDecidedConsensusData(runner Runner, val *types.ConsensusData) error {
	byts, err := val.Encode()
	if err != nil {
		return errors.Wrap(err, "could not encode decided value")
	}
	if err := runner.GetValCheckF()(byts); err != nil {
		return errors.Wrap(err, "decided value is invalid")
	}

	if runner.GetBeaconRole() != val.Duty.Type {
		return errors.New("decided value's duty has wrong beacon role type")
	}

	if !bytes.Equal(runner.GetShare().ValidatorPubKey, val.Duty.PubKey[:]) {
		return errors.New("decided value's validator pk is wrong")
	}

	return nil
}

func canProcessPostConsensusMsg(runner Runner, msg *SignedPartialSignatureMessage) error {
	if runner.GetState().DecidedValue == nil {
		return errors.New("consensus didn't decide")
	}

	if err := validatePartialSigMsg(runner, msg, runner.GetState().DecidedValue.Duty.Slot); err != nil {
		return errors.Wrap(err, "post consensus msg invalid")
	}

	return nil
}

func verifyExpectedRoot(runner Runner, signedMsg *SignedPartialSignatureMessage, expectedRootObjs []ssz.HashRoot, domain spec.DomainType) error {
	if len(expectedRootObjs) != len(signedMsg.Message.Messages) {
		return errors.New("wrong expected roots count")
	}
	for i, msg := range signedMsg.Message.Messages {
		epoch := runner.GetBeaconNetwork().EstimatedEpochAtSlot(runner.GetState().StartingDuty.Slot)
		d, err := runner.GetBeaconNode().DomainData(epoch, domain)
		if err != nil {
			return errors.Wrap(err, "could not get pre consensus root domain")
		}

		r, err := types.ComputeETHSigningRoot(expectedRootObjs[i], d)
		if !bytes.Equal(r[:], msg.SigningRoot) {
			return errors.New("wrong pre consensus signing root")
		}
	}
	return nil
}

func signPostConsensusMsg(runner Runner, msg *PartialSignatureMessages) (*SignedPartialSignatureMessage, error) {
	signature, err := runner.GetSigner().SignRoot(msg, types.PartialSignatureType, runner.GetShare().SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign PartialSignatureMessage for PostConsensusContainer")
	}

	return &SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: signature,
		Signer:    runner.GetShare().OperatorID,
	}, nil
}
