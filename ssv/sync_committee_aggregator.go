package ssv

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

type SyncCommitteeAggregatorRunner struct {
	BaseRunner *BaseRunner

	beacon   BeaconNode
	network  Network
	signer   types.KeyManager
	valCheck qbft.ProposedValueCheckF
}

func NewSyncCommitteeAggregatorRunner(
	beaconNetwork types.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.KeyManager,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) Runner {
	return &SyncCommitteeAggregatorRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType:     types.BNRoleSyncCommitteeContribution,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:   beacon,
		network:  network,
		signer:   signer,
		valCheck: valCheck,
	}
}

func (r *SyncCommitteeAggregatorRunner) StartNewDuty(duty *types.Duty) error {
	return r.BaseRunner.baseStartNewDuty(r, duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *SyncCommitteeAggregatorRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *SyncCommitteeAggregatorRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing sync committee selection proof message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// collect selection proofs and subnets
	var (
		selectionProofs []phase0.BLSSignature
		subnets         []uint64
	)
	for i, root := range roots {
		// reconstruct selection proof sig
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct sync committee index root")
		}
		blsSigSelectionProof := phase0.BLSSignature{}
		copy(blsSigSelectionProof[:], sig)

		aggregator, err := r.GetBeaconNode().IsSyncCommitteeAggregator(sig)
		if err != nil {
			return errors.Wrap(err, "could not check if sync committee aggregator")
		}
		if !aggregator {
			continue
		}

		// fetch sync committee contribution
		subnet, err := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(r.GetState().StartingDuty.ValidatorSyncCommitteeIndices[i]))
		if err != nil {
			return errors.Wrap(err, "could not get sync committee subnet ID")
		}

		selectionProofs = append(selectionProofs, blsSigSelectionProof)
		subnets = append(subnets, subnet)
	}
	if len(selectionProofs) == 0 {
		// there aren't any aggregators
		r.GetState().Finished = true
		return nil
	}

	duty := r.GetState().StartingDuty

	// fetch contributions
	contributions, ver, err := r.GetBeaconNode().GetSyncCommitteeContribution(duty.Slot, selectionProofs, subnets)
	if err != nil {
		return errors.Wrap(err, "could not get sync committee contribution")
	}
	byts, err := contributions.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal contributions")
	}

	// create consensus object
	input := &types.ConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (r *SyncCommitteeAggregatorRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	contributions, err := decidedValue.GetSyncCommitteeContributions()
	if err != nil {
		return errors.Wrap(err, "could not get contributions")
	}

	// specific duty sig
	msgs := make([]*types.PartialSignatureMessage, 0)
	for _, c := range contributions {
		contribAndProof, _, err := r.generateContributionAndProof(c.Contribution, c.SelectionProofSig)
		if err != nil {
			return errors.Wrap(err, "could not generate contribution and proof")
		}

		signed, err := r.BaseRunner.signBeaconObject(r, contribAndProof, decidedValue.Duty.Slot, types.DomainContributionAndProof)
		if err != nil {
			return errors.Wrap(err, "failed to sign aggregate and proof")
		}

		msgs = append(msgs, signed)
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     decidedValue.Duty.Slot,
		Messages: msgs,
	}

	postSignedMsg, err := r.BaseRunner.signPostConsensusMsg(r, postConsensusMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign post consensus msg")
	}

	data, err := postSignedMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.BeaconRoleType),
		Data:    data,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *SyncCommitteeAggregatorRunner) ProcessPostConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	// get contributions
	contributions, err := r.GetState().DecidedValue.GetSyncCommitteeContributions()
	if err != nil {
		return errors.Wrap(err, "could not get contributions")
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey)
		if err != nil {
			return errors.Wrap(err, "could not reconstruct post consensus signature")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		for _, contribution := range contributions {
			// match the right contrib and proof root to signed root
			contribAndProof, contribAndProofRoot, err := r.generateContributionAndProof(contribution.Contribution, contribution.SelectionProofSig)
			if err != nil {
				return errors.Wrap(err, "could not generate contribution and proof")
			}
			if !bytes.Equal(root[:], contribAndProofRoot[:]) {
				continue // not the correct root
			}

			signedContrib, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey)
			if err != nil {
				return errors.Wrap(err, "could not reconstruct contribution and proof sig")
			}
			blsSignedContribAndProof := phase0.BLSSignature{}
			copy(blsSignedContribAndProof[:], signedContrib)
			signedContribAndProof := &altair.SignedContributionAndProof{
				Message:   contribAndProof,
				Signature: blsSignedContribAndProof,
			}

			if err := r.GetBeaconNode().SubmitSignedContributionAndProof(signedContribAndProof); err != nil {
				return errors.Wrap(err, "could not submit to Beacon chain reconstructed contribution and proof")
			}
			break
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *SyncCommitteeAggregatorRunner) generateContributionAndProof(contrib altair.SyncCommitteeContribution, proof phase0.BLSSignature) (*altair.ContributionAndProof, phase0.Root, error) {
	contribAndProof := &altair.ContributionAndProof{
		AggregatorIndex: r.GetState().DecidedValue.Duty.ValidatorIndex,
		Contribution:    &contrib,
		SelectionProof:  proof,
	}

	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.GetState().DecidedValue.Duty.Slot)
	dContribAndProof, err := r.GetBeaconNode().DomainData(epoch, types.DomainContributionAndProof)
	if err != nil {
		return nil, phase0.Root{}, errors.Wrap(err, "could not get domain data")
	}
	contribAndProofRoot, err := types.ComputeETHSigningRoot(contribAndProof, dContribAndProof)
	if err != nil {
		return nil, phase0.Root{}, errors.Wrap(err, "could not compute signing root")
	}
	return contribAndProof, contribAndProofRoot, nil
}

func (r *SyncCommitteeAggregatorRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	sszIndexes := make([]ssz.HashRoot, 0)
	for _, index := range r.GetState().StartingDuty.ValidatorSyncCommitteeIndices {
		subnet, err := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(index))
		if err != nil {
			return nil, types.DomainError, errors.Wrap(err, "could not get sync committee subnet ID")
		}
		data := &altair.SyncAggregatorSelectionData{
			Slot:              r.GetState().StartingDuty.Slot,
			SubcommitteeIndex: subnet,
		}
		sszIndexes = append(sszIndexes, data)
	}
	return sszIndexes, types.DomainSyncCommitteeSelectionProof, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *SyncCommitteeAggregatorRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	// get contributions
	contributions, err := r.GetState().DecidedValue.GetSyncCommitteeContributions()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get contributions")
	}

	ret := make([]ssz.HashRoot, 0)
	for _, contrib := range contributions {
		contribAndProof, _, err := r.generateContributionAndProof(contrib.Contribution, contrib.SelectionProofSig)
		if err != nil {
			return nil, types.DomainError, errors.Wrap(err, "could not generate contribution and proof")
		}
		ret = append(ret, contribAndProof)
	}
	return ret, types.DomainContributionAndProof, nil
}

// executeDuty steps:
// 1) sign a partial contribution proof (for each subcommittee index) and wait for 2f+1 partial sigs from peers
// 2) Reconstruct contribution proofs, check IsSyncCommitteeAggregator and start consensus on duty + contribution data
// 3) Once consensus decides, sign partial contribution data (for each subcommittee) and broadcast
// 4) collect 2f+1 partial sigs, reconstruct and broadcast valid SignedContributionAndProof (for each subcommittee) sig to the BN
func (r *SyncCommitteeAggregatorRunner) executeDuty(duty *types.Duty) error {
	// sign selection proofs
	msgs := types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     duty.Slot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for _, index := range r.GetState().StartingDuty.ValidatorSyncCommitteeIndices {
		subnet, err := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(index))
		if err != nil {
			return errors.Wrap(err, "could not get sync committee subnet ID")
		}
		data := &altair.SyncAggregatorSelectionData{
			Slot:              duty.Slot,
			SubcommitteeIndex: subnet,
		}
		msg, err := r.BaseRunner.signBeaconObject(r, data, duty.Slot, types.DomainSyncCommitteeSelectionProof)
		if err != nil {
			return errors.Wrap(err, "could not sign sync committee selection proof")
		}

		msgs.Messages = append(msgs.Messages, msg)
	}

	// package into signed partial sig
	signature, err := r.GetSigner().SignRoot(msgs, types.PartialSignatureType, r.GetShare().SharePubKey)
	if err != nil {
		return errors.Wrap(err, "could not sign PartialSignatureMessage for contribution proofs")
	}
	signedPartialMsg := &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: signature,
		Signer:    r.GetShare().OperatorID,
	}

	// broadcast
	data, err := signedPartialMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode contribution proofs pre-consensus signature msg")
	}
	msgToBroadcast := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.BeaconRoleType),
		Data:    data,
	}
	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial contribution proof sig")
	}
	return nil
}

func (r *SyncCommitteeAggregatorRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *SyncCommitteeAggregatorRunner) GetNetwork() Network {
	return r.network
}

func (r *SyncCommitteeAggregatorRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *SyncCommitteeAggregatorRunner) GetShare() *types.Share {
	return r.BaseRunner.Share
}

func (r *SyncCommitteeAggregatorRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *SyncCommitteeAggregatorRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *SyncCommitteeAggregatorRunner) GetSigner() types.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *SyncCommitteeAggregatorRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *SyncCommitteeAggregatorRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *SyncCommitteeAggregatorRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
