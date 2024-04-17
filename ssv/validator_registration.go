package ssv

import (
	"crypto/sha256"
	"encoding/json"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

type ValidatorRegistrationRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewValidatorRegistrationRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
) Runner {
	return &ValidatorRegistrationRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: RoleValidatorRegistration,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
	}
}

func (r *ValidatorRegistrationRunner) StartNewDuty(duty types.Duty) error {
	r.executeDuty(duty)
	// Note: Validator registration doesn't require any consensus, it can start a new duty even if previous one didn't finish
	return r.BaseRunner.baseStartNewNonBeaconDuty(r, duty.(*types.BeaconDuty))
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *ValidatorRegistrationRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *ValidatorRegistrationRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing validator registration message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	root := roots[0]
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey[:])
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root, r.GetShare().Committee)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}
	specSig := phase0.BLSSignature{}
	copy(specSig[:], fullSig)

	if err := r.beacon.SubmitValidatorRegistration(r.BaseRunner.Share[0].ValidatorPubKey[:],
		r.BaseRunner.Share[0].FeeRecipientAddress, specSig); err != nil {
		return errors.Wrap(err, "could not submit validator registration")
	}

	r.GetState().Finished = true
	return nil
}

func (r *ValidatorRegistrationRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	return errors.New("no consensus phase for validator registration")
}

func (r *ValidatorRegistrationRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	return errors.New("no post consensus phase for validator registration")
}

func (r *ValidatorRegistrationRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	vr, err := r.calculateValidatorRegistration()
	if err != nil {
		return nil, types.DomainError, errors.Wrap(err, "could not calculate validator registration")
	}
	return []ssz.HashRoot{vr}, types.DomainApplicationBuilder, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ValidatorRegistrationRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, errors.New("no post consensus roots for validator registration")
}

func (r *ValidatorRegistrationRunner) executeDuty(duty types.Duty) error {
	vr, err := r.calculateValidatorRegistration()
	if err != nil {
		return errors.Wrap(err, "could not calculate validator registration")
	}

	// sign partial randao
	msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.BeaconDuty), vr, duty.DutySlot(),
		types.DomainApplicationBuilder)
	if err != nil {
		return errors.Wrap(err, "could not sign validator registration")
	}
	msgs := &types.PartialSignatureMessages{
		Type:     types.ValidatorRegistrationPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{msg},
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], RunnerRole(r.BaseRunner.RunnerRoleType))
	msgToBroadcast, err := types.PartialSignatureMessagesToSignedSSVMessage(msgs, msgID, r.operatorSigner)
	if err != nil {
		return errors.Wrap(err, "could not sign pre-consensus partial signature message")
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial randao sig")
	}
	return nil
}

func (r *ValidatorRegistrationRunner) calculateValidatorRegistration() (*v1.ValidatorRegistration, error) {
	pk := phase0.BLSPubKey{}
	copy(pk[:], r.BaseRunner.Share[0].ValidatorPubKey[:])

	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())

	return &v1.ValidatorRegistration{
		FeeRecipient: r.BaseRunner.Share[0].FeeRecipientAddress,
		GasLimit:     types.DefaultGasLimit,
		Timestamp:    r.BaseRunner.BeaconNetwork.EpochStartTime(epoch),
		Pubkey:       pk,
	}, nil
}

func (r *ValidatorRegistrationRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *ValidatorRegistrationRunner) GetNetwork() Network {
	return r.network
}

func (r *ValidatorRegistrationRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *ValidatorRegistrationRunner) GetShare() *types.Share {
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *ValidatorRegistrationRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ValidatorRegistrationRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *ValidatorRegistrationRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *ValidatorRegistrationRunner) GetOperatorSigner() types.OperatorSigner {
	return r.operatorSigner
}

// Encode returns the encoded struct in bytes or error
func (r *ValidatorRegistrationRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *ValidatorRegistrationRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *ValidatorRegistrationRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
