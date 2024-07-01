package ssv

import (
	"bytes"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type SyncCommitteeAggregatorRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewSyncCommitteeAggregatorRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) Runner {
	return &SyncCommitteeAggregatorRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     types.RoleSyncCommitteeContribution,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
	}
}

func (r *SyncCommitteeAggregatorRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	return r.BaseRunner.baseStartNewDuty(r, duty, quorum)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *SyncCommitteeAggregatorRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *SyncCommitteeAggregatorRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
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
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root,
			r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			for _, root := range roots {
				r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root,
					r.GetShare().Committee, r.GetShare().ValidatorIndex)
			}
			return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
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
		subnet, err := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(r.GetState().StartingDuty.(*types.ValidatorDuty).ValidatorSyncCommitteeIndices[i]))
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

	duty := r.GetState().StartingDuty.(*types.ValidatorDuty)

	// fetch contributions
	contributions, ver, err := r.GetBeaconNode().GetSyncCommitteeContribution(duty.DutySlot(), selectionProofs, subnets)
	if err != nil {
		return errors.Wrap(err, "could not get sync committee contribution")
	}
	byts, err := contributions.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal contributions")
	}

	// create consensus object
	input := &types.ValidatorConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: byts,
	}

	if err := r.BaseRunner.decide(r, input.Duty.DutySlot(), input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (r *SyncCommitteeAggregatorRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	cd := decidedValue.(*types.ValidatorConsensusData)
	contributions, err := cd.GetSyncCommitteeContributions()
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

		signed, err := r.BaseRunner.signBeaconObject(r, r.BaseRunner.State.StartingDuty.(*types.ValidatorDuty), contribAndProof,
			cd.Duty.Slot,
			types.DomainContributionAndProof)
		if err != nil {
			return errors.Wrap(err, "failed to sign aggregate and proof")
		}

		msgs = append(msgs, signed)
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     cd.Duty.Slot,
		Messages: msgs,
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := postConsensusMsg.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{r.operatorSigner.GetOperatorID()},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast.SSVMessage.GetID(), msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *SyncCommitteeAggregatorRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	// get contributions
	validatorConsensusData := &types.ValidatorConsensusData{}
	err = validatorConsensusData.Decode(r.GetState().DecidedValue)
	if err != nil {
		return errors.Wrap(err, "could not create consensus data")
	}
	contributions, err := validatorConsensusData.GetSyncCommitteeContributions()
	if err != nil {
		return errors.Wrap(err, "could not get contributions")
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
		if err != nil {
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			for _, root := range roots {
				r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PostConsensusContainer, root,
					r.GetShare().Committee, r.GetShare().ValidatorIndex)
			}
			return errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
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

			signedContrib, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey[:], r.GetShare().ValidatorIndex)
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
		AggregatorIndex: r.GetState().StartingDuty.(*types.ValidatorDuty).ValidatorIndex,
		Contribution:    &contrib,
		SelectionProof:  proof,
	}

	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.GetState().StartingDuty.DutySlot())
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
	for _, index := range r.GetState().StartingDuty.(*types.ValidatorDuty).ValidatorSyncCommitteeIndices {
		subnet, err := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(index))
		if err != nil {
			return nil, types.DomainError, errors.Wrap(err, "could not get sync committee subnet ID")
		}
		data := &altair.SyncAggregatorSelectionData{
			Slot:              r.GetState().StartingDuty.DutySlot(),
			SubcommitteeIndex: subnet,
		}
		sszIndexes = append(sszIndexes, data)
	}
	return sszIndexes, types.DomainSyncCommitteeSelectionProof, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *SyncCommitteeAggregatorRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	// get contributions
	validatorConsensusData := &types.ValidatorConsensusData{}
	err := validatorConsensusData.Decode(r.GetState().DecidedValue)
	if err != nil {
		return nil, types.DomainError, errors.Wrap(err, "could not create consensus data")
	}
	contributions, err := validatorConsensusData.GetSyncCommitteeContributions()
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
func (r *SyncCommitteeAggregatorRunner) executeDuty(duty types.Duty) error {
	// sign selection proofs
	msgs := &types.PartialSignatureMessages{
		Type:     types.ContributionProofs,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{},
	}
	for _, index := range r.GetState().StartingDuty.(*types.ValidatorDuty).ValidatorSyncCommitteeIndices {
		subnet, err := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(index))
		if err != nil {
			return errors.Wrap(err, "could not get sync committee subnet ID")
		}
		data := &altair.SyncAggregatorSelectionData{
			Slot:              duty.DutySlot(),
			SubcommitteeIndex: subnet,
		}
		msg, err := r.BaseRunner.signBeaconObject(r, duty.(*types.ValidatorDuty), data, duty.DutySlot(),
			types.DomainSyncCommitteeSelectionProof)
		if err != nil {
			return errors.Wrap(err, "could not sign sync committee selection proof")
		}

		msgs.Messages = append(msgs.Messages, msg)
	}

	msgID := types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := msgs.Encode()
	if err != nil {
		return err
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{r.operatorSigner.GetOperatorID()},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgToBroadcast.SSVMessage.GetID(), msgToBroadcast); err != nil {
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
	// there is only one share
	for _, share := range r.BaseRunner.Share {
		return share
	}
	return nil
}

func (r *SyncCommitteeAggregatorRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *SyncCommitteeAggregatorRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *SyncCommitteeAggregatorRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *SyncCommitteeAggregatorRunner) GetOperatorSigner() types.OperatorSigner {
	return r.operatorSigner
}
