package ssv

import (
	"bytes"
	"sort"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// AggregatorCommitteeRunner runs aggregate and sync committee contribution duties
type AggregatorCommitteeRunner struct {
	BaseRunner     *BaseRunner
	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner *types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
	// For aggregator role: tracks by validator index only (one submission per validator)
	// For sync committee contribution role: tracks by validator index and root (multiple submissions per validator)
	submittedDuties map[types.BeaconRole]map[phase0.ValidatorIndex]map[[32]byte]struct{}

	// Pre-consensus markers:
	// - seen signers.
	preConsensusSeenSigners map[types.OperatorID]struct{}
	// - (validator index, beacon root) tuples that got checked for aggregation/scc selection with the beacon node.
	preConsensusSelectionsTested map[phase0.ValidatorIndex]map[[32]byte]struct{}
}

// NewAggregatorCommitteeRunner creates a new aggregator committee runner
func NewAggregatorCommitteeRunner(
	beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
) (Runner, error) {
	if len(share) == 0 {
		return nil, errors.New("no shares provided")
	}

	return &AggregatorCommitteeRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleAggregatorCommittee,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
			QBFTController: qbftController,
		},

		beacon:                       beacon,
		network:                      network,
		signer:                       signer,
		operatorSigner:               operatorSigner,
		valCheck:                     valCheck,
		submittedDuties:              make(map[types.BeaconRole]map[phase0.ValidatorIndex]map[[32]byte]struct{}),
		preConsensusSeenSigners:      make(map[types.OperatorID]struct{}),
		preConsensusSelectionsTested: make(map[phase0.ValidatorIndex]map[[32]byte]struct{}),
	}, nil
}

func (r *AggregatorCommitteeRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	err := r.BaseRunner.baseStartNewDuty(r, duty, quorum)
	if err != nil {
		return errors.Wrap(err, "failed to start new duty")
	}

	// Initialize submission tracking for both duty types
	r.submittedDuties[types.BNRoleAggregator] = make(map[phase0.ValidatorIndex]map[[32]byte]struct{})
	r.submittedDuties[types.BNRoleSyncCommitteeContribution] = make(map[phase0.ValidatorIndex]map[[32]byte]struct{})

	return nil
}

func (r *AggregatorCommitteeRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *AggregatorCommitteeRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	// Mark signer as seen
	r.MarkPreConsensusSignerAsSeen(signedMsg)

	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing pre-consensus message")
	}

	if !quorum {
		// If didn't get any new quorum, didn't yet start QBFT, and has received the last message, then terminate.
		if r.HasSeenAllPreConsensusSigners() {
			if r.BaseRunner.QBFTController.InstanceForHeight(qbft.Height(r.BaseRunner.State.StartingDuty.DutySlot())) == nil {
				r.BaseRunner.State.Finished = true
			}
		}
		return nil
	}

	aggregatorMap, contributionMap, err := r.expectedPreConsensusRoots()
	if err != nil {
		return errors.Wrap(err, "could not get expected pre-consensus roots")
	}

	duty := r.BaseRunner.State.StartingDuty.(*types.AggregatorCommitteeDuty)
	epoch := r.beacon.GetBeaconNetwork().EstimatedEpochAtSlot(duty.Slot)
	consensusData := &types.AggregatorCommitteeConsensusData{
		Version: r.beacon.DataVersion(epoch),
	}
	hasAnyAggregator := false

	rootSet := make(map[[32]byte]struct{})
	for _, root := range roots {
		rootSet[root] = struct{}{}
	}

	var sortedRoots [][32]byte
	for root := range rootSet {
		sortedRoots = append(sortedRoots, root)
	}
	sort.Slice(sortedRoots, func(i, j int) bool {
		return bytes.Compare(sortedRoots[i][:], sortedRoots[j][:]) < 0
	})

	var anyErr error
	for _, root := range sortedRoots {
		metadataList, found := findValidatorsForPreConsensusRoot(root, aggregatorMap, contributionMap)
		if !found {
			// Edge case: since operators may have divergent sets of validators,
			// it's possible that an operator doesn't have the validator associated to a root.
			// In this case, we simply continue.
			continue
		}

		sort.Slice(metadataList, func(i, j int) bool {
			return metadataList[i].ValidatorIndex < metadataList[j].ValidatorIndex
		})

		for _, metadata := range metadataList {
			validatorIndex := metadata.ValidatorIndex
			share := r.BaseRunner.Share[validatorIndex]
			if share == nil {
				continue
			}

			if !r.BaseRunner.State.PreConsensusContainer.HasQuorum(validatorIndex, root) {
				continue
			}

			// Reconstruct signature
			fullSig, err := r.BaseRunner.State.ReconstructBeaconSig(
				r.BaseRunner.State.PreConsensusContainer,
				root,
				share.ValidatorPubKey[:],
				validatorIndex,
			)
			if err != nil {
				// Fallback: verify each signature individually for all roots
				for root := range rootSet {
					r.BaseRunner.FallBackAndVerifyEachSignature(
						r.BaseRunner.State.PreConsensusContainer,
						root,
						share.Committee,
						validatorIndex,
					)
				}
				// Record the error and continue to next validators
				anyErr = errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
				continue
			}

			var blsSig phase0.BLSSignature
			copy(blsSig[:], fullSig)

			switch metadata.Role {
			case types.BNRoleAggregator:
				vDuty := r.findValidatorDuty(validatorIndex, types.BNRoleAggregator)
				if vDuty != nil {
					isAggregator, err := r.processAggregatorSelectionProof(blsSig, vDuty, consensusData)
					if err == nil {
						r.MarkSelectionChecked(validatorIndex, root)
						if isAggregator {
							hasAnyAggregator = true
						}
					} else {
						anyErr = errors.Wrap(err, "failed to process aggregator selection proof")
					}
				}

			case types.BNRoleSyncCommitteeContribution:
				vDuty := r.findValidatorDuty(validatorIndex, types.BNRoleSyncCommitteeContribution)
				if vDuty != nil {
					isAggregator, err := r.processSyncCommitteeSelectionProof(blsSig, metadata.ValidatorSyncCommitteeIndex, vDuty, consensusData)
					if err == nil {
						r.MarkSelectionChecked(validatorIndex, root)
						if isAggregator {
							hasAnyAggregator = true
						}
					} else {
						anyErr = errors.Wrap(err, "failed to process sync committee selection proof")
					}
				}

			default:
				// This should never happen as we build rootToMetadata ourselves with valid roles
				return errors.Errorf("unexpected role type in pre-consensus metadata: %v", metadata.Role)
			}
		}
	}

	// No aggregators selected: decide whether to finish or wait for more messages
	if !hasAnyAggregator {
		// If all duties have been tested for selection or all messages (from all operators) have been seen, terminate.
		if r.AllSelectionsChecked(aggregatorMap, contributionMap) || r.HasSeenAllPreConsensusSigners() {
			r.BaseRunner.State.Finished = true
			return nil
		}

		// If no validator was selected, but there are more possible messages (and thus selections), keep waiting for more messages.
		// This helps on our liveness vs. completeness trade-off, inherent to the pre-consensus phase of the aggregator committee runner.
		return anyErr
	}

	// Else, if some aggregators or contributors were selected (even with an error for others), proceed to consensus
	if err := consensusData.Validate(); err != nil {
		return errors.Wrap(err, "invalid aggregator consensus data")
	}

	if err := r.BaseRunner.decide(r, r.BaseRunner.State.StartingDuty.DutySlot(), consensusData); err != nil {
		return errors.Wrap(err, "failed to start consensus")
	}

	// Raise error if any
	if anyErr != nil {
		return anyErr
	}

	return nil
}

// ProcessConsensus processes consensus messages
func (r *AggregatorCommitteeRunner) ProcessConsensus(signedMsg *types.SignedSSVMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg, &types.AggregatorCommitteeConsensusData{})
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	if !decided {
		return nil
	}

	consensusData := decidedValue.(*types.AggregatorCommitteeConsensusData)

	var messages []*types.PartialSignatureMessage

	aggProofs, err := consensusData.GetAggregateAndProofs()
	if err != nil {
		return errors.Wrap(err, "failed to get aggregate and proofs")
	}

	for i, aggProof := range aggProofs {
		validatorIndex := consensusData.Aggregators[i].ValidatorIndex

		_, exists := r.BaseRunner.Share[validatorIndex]
		if !exists {
			continue
		}

		vDuty := r.findValidatorDuty(validatorIndex, types.BNRoleAggregator)
		if vDuty == nil {
			continue
		}

		// Sign the aggregate and proof
		hashRoot, err := types.GetAggregateAndProofHashRoot(aggProof)
		if err != nil {
			return errors.Wrap(err, "failed to get aggregate and proof hash root")
		}
		msg, err := r.BaseRunner.signBeaconObject(
			r, vDuty, hashRoot,
			r.BaseRunner.State.StartingDuty.DutySlot(),
			types.DomainAggregateAndProof,
		)
		if err != nil {
			return errors.Wrap(err, "failed to sign aggregate and proof")
		}

		messages = append(messages, msg)
	}

	contributions, err := consensusData.GetSyncCommitteeContributions()
	if err != nil {
		return errors.Wrap(err, "failed to get sync committee contributions")
	}

	for i, contribution := range contributions {
		validatorIndex := consensusData.Contributors[i].ValidatorIndex

		_, exists := r.BaseRunner.Share[validatorIndex]
		if !exists {
			continue
		}

		vDuty := r.findValidatorDuty(validatorIndex, types.BNRoleSyncCommitteeContribution)
		if vDuty == nil {
			continue
		}

		contribAndProof := &altair.ContributionAndProof{
			AggregatorIndex: validatorIndex,
			Contribution:    &contribution.Contribution,
			SelectionProof:  consensusData.Contributors[i].SelectionProof,
		}

		// Sign the contribution and proof
		msg, err := r.BaseRunner.signBeaconObject(
			r, vDuty, contribAndProof,
			r.BaseRunner.State.StartingDuty.DutySlot(),
			types.DomainContributionAndProof,
		)
		if err != nil {
			return errors.Wrap(err, "failed to sign contribution and proof")
		}

		messages = append(messages, msg)
	}

	// If no messages to broadcast, we're done
	if len(messages) == 0 {
		return nil
	}

	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     r.BaseRunner.State.StartingDuty.DutySlot(),
		Messages: messages,
	}

	return r.broadcastPartialSignatureMessage(postConsensusMsg)
}

func (r *AggregatorCommitteeRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, rootsList, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	rootSet := make(map[[32]byte]struct{})
	for _, root := range rootsList {
		rootSet[root] = struct{}{}
	}

	var sortedRoots [][32]byte
	for root := range rootSet {
		sortedRoots = append(sortedRoots, root)
	}
	sort.Slice(sortedRoots, func(i, j int) bool {
		return bytes.Compare(sortedRoots[i][:], sortedRoots[j][:]) < 0
	})

	aggregatorMap, contributionMap, beaconObjects, err := r.expectedPostConsensusRootsAndBeaconObjects()
	if err != nil {
		return errors.Wrap(err, "could not get expected post consensus roots and beacon objects")
	}

	var anyErr error

	type AggregationSubmission struct {
		spec.VersionedSignedAggregateAndProof
		root [32]byte
	}
	type ContributionSubmission struct {
		*altair.SignedContributionAndProof
		root [32]byte
	}
	aggregationsToSubmit := make(map[phase0.ValidatorIndex]*AggregationSubmission)
	syncCommitteeContributionsToSubmit := make(map[phase0.ValidatorIndex][]*ContributionSubmission)

	// For each root that got at least one quorum, find the duties and try to submit
	for _, root := range sortedRoots {
		// Get validators related to the given root
		role, validators, found := findValidatorsForPostConsensusRoot(root, aggregatorMap, contributionMap)

		if !found {
			// Edge case: operator doesn't have the validator associated to a root
			continue
		}

		for _, validator := range validators {
			if !r.BaseRunner.State.PostConsensusContainer.HasQuorum(validator, root) {
				continue
			}

			if r.HasSubmitted(role, validator, root) {
				continue
			}

			// Reconstruct signature
			share := r.BaseRunner.Share[validator]
			if share == nil {
				continue
			}

			sig, err := r.BaseRunner.State.ReconstructBeaconSig(
				r.BaseRunner.State.PostConsensusContainer,
				root,
				share.ValidatorPubKey[:],
				validator,
			)
			if err != nil {
				// Fallback: verify each signature individually
				for root := range rootSet {
					r.BaseRunner.FallBackAndVerifyEachSignature(
						r.BaseRunner.State.PostConsensusContainer,
						root,
						share.Committee,
						validator,
					)
				}
				anyErr = types.WrapError(types.PostConsensusQuorumWithInvalidSignatures, errors.Wrap(err, "got post-consensus quorum but it has invalid signatures"))
				continue
			}

			var blsSig phase0.BLSSignature
			copy(blsSig[:], sig)

			validatorObjs, exists := beaconObjects[validator]
			if !exists {
				anyErr = errors.New("could not find beacon objects for validator")
				continue
			}
			beaconObj, exists := validatorObjs[root]
			if !exists {
				anyErr = errors.New("could not find beacon object for root")
				continue
			}

			// Submit based on role
			switch role {
			case types.BNRoleAggregator:
				aggregateAndProof := beaconObj.(*spec.VersionedAggregateAndProof)
				signedAgg, err := r.constructSignedAggregateAndProof(aggregateAndProof, blsSig)
				if err != nil {
					anyErr = errors.Wrap(err, "failed to construct signed aggregate and proof")
					continue
				}

				// Store for later submission
				aggregationsToSubmit[validator] = &AggregationSubmission{
					VersionedSignedAggregateAndProof: *signedAgg,
					root:                             root,
				}

			case types.BNRoleSyncCommitteeContribution:
				contribAndProof := beaconObj.(*altair.ContributionAndProof)
				signedContrib := &altair.SignedContributionAndProof{
					Message:   contribAndProof,
					Signature: blsSig,
				}

				// Store for later submission
				if _, ok := syncCommitteeContributionsToSubmit[validator]; !ok {
					syncCommitteeContributionsToSubmit[validator] = make([]*ContributionSubmission, 0)
				}
				syncCommitteeContributionsToSubmit[validator] = append(syncCommitteeContributionsToSubmit[validator], &ContributionSubmission{
					SignedContributionAndProof: signedContrib,
					root:                       root,
				})
			default:
				return errors.Errorf("unexpected role type in post-consensus: %v", role)
			}
		}
	}

	// Submit aggregates and contributions
	for validator, aggSubmission := range aggregationsToSubmit {
		if err := r.beacon.SubmitSignedAggregateAndProof(&aggSubmission.VersionedSignedAggregateAndProof); err != nil {
			anyErr = errors.Wrap(err, "failed to submit signed aggregate and proof")
			continue
		}
		r.RecordSubmission(types.BNRoleAggregator, validator, aggSubmission.root)
	}
	for validator, contribSubmissions := range syncCommitteeContributionsToSubmit {
		for _, contribSubmission := range contribSubmissions {
			if err := r.beacon.SubmitSignedContributionAndProof(contribSubmission.SignedContributionAndProof); err != nil {
				anyErr = errors.Wrap(err, "failed to submit signed contribution and proof")
				continue
			}
			r.RecordSubmission(types.BNRoleSyncCommitteeContribution, validator, contribSubmission.root)
		}
	}

	if anyErr != nil {
		return anyErr
	}

	if r.HasSubmittedAllDuties() {
		r.BaseRunner.State.Finished = true
	}

	return nil
}

// executeDuty executes the aggregator committee duty by creating selection proofs
func (r *AggregatorCommitteeRunner) executeDuty(duty types.Duty) error {
	aggCommitteeDuty, ok := duty.(*types.AggregatorCommitteeDuty)
	if !ok {
		return errors.New("invalid duty type for aggregator committee runner")
	}

	// Validate duty
	valIdxs := make(map[phase0.ValidatorIndex]struct{})
	for idx := range r.BaseRunner.Share {
		valIdxs[idx] = struct{}{}
	}
	if err := aggCommitteeDuty.Validate(valIdxs); err != nil {
		return err
	}

	msg := &types.PartialSignatureMessages{
		Type:     types.AggregatorCommitteePartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{},
	}

	// Generate selection proofs for all validators and duties
	for _, vDuty := range aggCommitteeDuty.ValidatorDuties {
		if _, ok := r.BaseRunner.Share[vDuty.ValidatorIndex]; !ok {
			continue
		}

		switch vDuty.Type {
		case types.BNRoleAggregator:
			// Sign slot for aggregator selection proof
			partialSig, err := r.BaseRunner.signBeaconObject(
				r, vDuty, types.SSZUint64(duty.DutySlot()),
				duty.DutySlot(), types.DomainSelectionProof)
			if err != nil {
				return errors.Wrap(err, "failed to sign aggregator selection proof")
			}

			msg.Messages = append(msg.Messages, partialSig)

		case types.BNRoleSyncCommitteeContribution:
			// Sign sync committee selection proofs for each subcommittee
			for _, index := range vDuty.ValidatorSyncCommitteeIndices {
				subnet := r.GetBeaconNode().SyncCommitteeSubnetID(phase0.CommitteeIndex(index))

				data := &altair.SyncAggregatorSelectionData{
					Slot:              duty.DutySlot(),
					SubcommitteeIndex: subnet,
				}

				partialSig, err := r.BaseRunner.signBeaconObject(
					r, vDuty, data, duty.DutySlot(),
					types.DomainSyncCommitteeSelectionProof)
				if err != nil {
					return errors.Wrap(err, "failed to sign sync committee selection proof")
				}

				msg.Messages = append(msg.Messages, partialSig)
			}

		default:
			return errors.Errorf("invalid validator duty type for aggregator committee: %v", vDuty.Type)
		}
	}

	// Early exit if no selection proofs needed
	if len(msg.Messages) == 0 {
		r.BaseRunner.State.Finished = true
		return nil
	}

	// Broadcast the selection proofs
	return r.broadcastPartialSignatureMessage(msg)
}

// expectedPreConsensusRootsAndDomain returns expected pre-consensus roots
func (r *AggregatorCommitteeRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	panic("not in use")
}

// expectedPostConsensusRootsAndDomain returns expected post-consensus roots
func (r *AggregatorCommitteeRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	panic("not in use")
}

func (r *AggregatorCommitteeRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *AggregatorCommitteeRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

// broadcastPartialSignatureMessage handles the common pattern of creating and broadcasting SSV messages
func (r *AggregatorCommitteeRunner) broadcastPartialSignatureMessage(
	msg *types.PartialSignatureMessages,
) error {
	committeeMember := r.BaseRunner.QBFTController.CommitteeMember
	if committeeMember == nil {
		return errors.New("committee member is nil")
	}

	operatorIDs := make([]types.OperatorID, 0, len(committeeMember.Committee))
	for _, op := range committeeMember.Committee {
		operatorIDs = append(operatorIDs, op.OperatorID)
	}
	committeeID := types.GetCommitteeID(operatorIDs)

	msgID := types.NewMsgID(committeeMember.DomainType, committeeID[:], r.BaseRunner.RunnerRoleType)

	encodedMsg, err := msg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode message")
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}

	// Sign the message
	sig, err := r.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{committeeMember.OperatorID},
		SSVMessage:  ssvMsg,
	}

	if err := r.GetNetwork().Broadcast(msgID, msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast message")
	}

	return nil
}

func (r *AggregatorCommitteeRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *AggregatorCommitteeRunner) GetSigner() types.BeaconSigner {
	return r.signer
}

func (r *AggregatorCommitteeRunner) GetOperatorSigner() *types.OperatorSigner {
	return r.operatorSigner
}

func (r *AggregatorCommitteeRunner) GetNetwork() Network {
	return r.network
}

// findValidatorDuty finds the validator duty for a specific role
func (r *AggregatorCommitteeRunner) findValidatorDuty(validatorIndex phase0.ValidatorIndex, role types.BeaconRole) *types.ValidatorDuty {
	duty := r.BaseRunner.State.StartingDuty.(*types.AggregatorCommitteeDuty)

	for _, d := range duty.ValidatorDuties {
		if d.ValidatorIndex == validatorIndex && d.Type == role {
			return d
		}
	}

	return nil
}

// processAggregatorSelectionProof handles aggregator selection proofs
func (r *AggregatorCommitteeRunner) processAggregatorSelectionProof(
	selectionProof phase0.BLSSignature,
	vDuty *types.ValidatorDuty,
	aggregatorData *types.AggregatorCommitteeConsensusData,
) (bool, error) {
	isAggregator := r.beacon.IsAggregator(
		vDuty.Slot, vDuty.CommitteeIndex, vDuty.CommitteeLength, selectionProof[:])

	if !isAggregator {
		// Not selected as aggregator
		return false, nil
	}

	// Check if attestation for committee index was already included
	for _, idx := range aggregatorData.AggregatorsCommitteeIndexes {
		if idx == uint64(vDuty.CommitteeIndex) {
			// If so, just add to aggregators and return
			aggregatorData.Aggregators = append(aggregatorData.Aggregators, types.AssignedAggregator{
				ValidatorIndex: vDuty.ValidatorIndex,
				SelectionProof: selectionProof,
				CommitteeIndex: uint64(vDuty.CommitteeIndex),
			})
			return true, nil
		}
	}

	// Else, fetch attestation and include everything (if successful)

	attestation, err := r.beacon.GetAggregateAttestation(vDuty.Slot, vDuty.CommitteeIndex)
	if err != nil {
		return true, errors.Wrap(err, "failed to get aggregate attestation")
	}

	aggregatorData.Aggregators = append(aggregatorData.Aggregators, types.AssignedAggregator{
		ValidatorIndex: vDuty.ValidatorIndex,
		SelectionProof: selectionProof,
		CommitteeIndex: uint64(vDuty.CommitteeIndex),
	})

	// Marshal attestation for storage
	attestationBytes, err := attestation.MarshalSSZ()
	if err != nil {
		return true, errors.Wrap(err, "failed to marshal attestation")
	}

	aggregatorData.AggregatorsCommitteeIndexes = append(aggregatorData.AggregatorsCommitteeIndexes, uint64(vDuty.CommitteeIndex))
	aggregatorData.AggregatedAttestations = append(aggregatorData.AggregatedAttestations, attestationBytes)

	return true, nil
}

// processSyncCommitteeSelectionProof handles sync committee selection proofs with known index
func (r *AggregatorCommitteeRunner) processSyncCommitteeSelectionProof(
	selectionProof phase0.BLSSignature,
	validatorSyncCommitteeIndex types.ValidatorSyncCommitteeIndex,
	vDuty *types.ValidatorDuty,
	aggregatorData *types.AggregatorCommitteeConsensusData,
) (bool, error) {
	subnetID := r.beacon.SyncCommitteeSubnetID(phase0.CommitteeIndex(validatorSyncCommitteeIndex))

	isAggregator := r.beacon.IsSyncCommitteeAggregator(selectionProof[:])

	if !isAggregator {
		return false, nil // Not selected as sync committee aggregator
	}

	// Check if we already have a contribution for this sync committee subnet ID
	for _, contrib := range aggregatorData.SyncCommitteeContributions {
		if contrib.SubcommitteeIndex == subnetID {
			// If so, just add to contributors and return
			aggregatorData.Contributors = append(aggregatorData.Contributors, types.AssignedAggregator{
				ValidatorIndex: vDuty.ValidatorIndex,
				SelectionProof: selectionProof,
				CommitteeIndex: subnetID,
			})
			return true, nil
		}
	}

	// Else, fetch contribution and include everything (if successful)

	contributions, _, err := r.GetBeaconNode().GetSyncCommitteeContribution(
		vDuty.Slot, []phase0.BLSSignature{selectionProof}, []uint64{subnetID})
	if err != nil {
		return true, err
	}

	// Type assertion to get the actual Contributions object
	contribs, ok := contributions.(*types.Contributions)
	if !ok {
		return true, errors.Errorf("unexpected contributions type: %T", contributions)
	}

	if len(*contribs) == 0 {
		return true, errors.New("no contributions found")
	}

	// Append the contribution(s)
	for _, contrib := range *contribs {
		if contrib.Contribution.SubcommitteeIndex != subnetID {
			continue
		}

		aggregatorData.Contributors = append(aggregatorData.Contributors, types.AssignedAggregator{
			ValidatorIndex: vDuty.ValidatorIndex,
			SelectionProof: selectionProof,
			CommitteeIndex: subnetID,
		})

		aggregatorData.SyncCommitteeContributions = append(aggregatorData.SyncCommitteeContributions, contrib.Contribution)
	}

	return true, nil
}

// expectedAggregatorSelectionRoot calculates the expected signing root for aggregator selection
func (r *AggregatorCommitteeRunner) expectedAggregatorSelectionRoot(
	_ *types.ValidatorDuty,
	slot phase0.Slot,
) ([32]byte, error) {
	epoch := r.beacon.GetBeaconNetwork().EstimatedEpochAtSlot(slot)
	domain, err := r.beacon.DomainData(epoch, types.DomainSelectionProof)
	if err != nil {
		return [32]byte{}, err
	}

	return types.ComputeETHSigningRoot(types.SSZUint64(slot), domain)
}

// expectedSyncCommitteeSelectionRoot calculates the expected signing root for sync committee selection
func (r *AggregatorCommitteeRunner) expectedSyncCommitteeSelectionRoot(
	_ *types.ValidatorDuty,
	slot phase0.Slot,
	validatorSyncCommitteeIndex types.ValidatorSyncCommitteeIndex,
) ([32]byte, error) {
	subnet := r.beacon.SyncCommitteeSubnetID(phase0.CommitteeIndex(validatorSyncCommitteeIndex))

	data := &altair.SyncAggregatorSelectionData{
		Slot:              slot,
		SubcommitteeIndex: subnet,
	}

	epoch := r.beacon.GetBeaconNetwork().EstimatedEpochAtSlot(slot)
	domain, err := r.beacon.DomainData(epoch, types.DomainSyncCommitteeSelectionProof)
	if err != nil {
		return [32]byte{}, err
	}

	return types.ComputeETHSigningRoot(data, domain)
}

// expectedPreConsensusRoots returns the expected roots for the pre-consensus phase.
// It returns the aggregator and sync committee validator to root maps.
func (r *AggregatorCommitteeRunner) expectedPreConsensusRoots() (
	aggregatorMap map[phase0.ValidatorIndex][32]byte,
	contributionMap map[phase0.ValidatorIndex]map[types.ValidatorSyncCommitteeIndex][32]byte,
	error error,
) {
	aggregatorMap = make(map[phase0.ValidatorIndex][32]byte)
	contributionMap = make(map[phase0.ValidatorIndex]map[types.ValidatorSyncCommitteeIndex][32]byte)

	duty := r.BaseRunner.State.StartingDuty.(*types.AggregatorCommitteeDuty)

	for _, vDuty := range duty.ValidatorDuties {
		if vDuty == nil {
			continue
		}

		switch vDuty.Type {
		case types.BNRoleAggregator:
			root, err := r.expectedAggregatorSelectionRoot(vDuty, duty.Slot)
			if err != nil {
				continue
			}
			aggregatorMap[vDuty.ValidatorIndex] = root

		case types.BNRoleSyncCommitteeContribution:
			if _, ok := contributionMap[vDuty.ValidatorIndex]; !ok {
				contributionMap[vDuty.ValidatorIndex] = make(map[uint64][32]byte)
			}

			for _, index := range vDuty.ValidatorSyncCommitteeIndices {
				root, err := r.expectedSyncCommitteeSelectionRoot(vDuty, duty.Slot, index)
				if err != nil {
					continue
				}
				contributionMap[vDuty.ValidatorIndex][index] = root
			}

		default:
			return nil, nil, errors.Errorf("invalid duty type in aggregator committee duty: %v", vDuty.Type)
		}
	}

	return aggregatorMap, contributionMap, nil
}

// expectedPostConsensusRootsAndBeaconObjects returns the expected roots and beacon objects for the post-consensus
// phase. It returns the aggregate and sync committee contribution validator to root maps, as well as beacon objects.
func (r *AggregatorCommitteeRunner) expectedPostConsensusRootsAndBeaconObjects() (
	aggregatorMap map[phase0.ValidatorIndex][32]byte,
	contributionMap map[phase0.ValidatorIndex][][32]byte,
	beaconObjects map[phase0.ValidatorIndex]map[[32]byte]interface{},
	error error,
) {
	aggregatorMap = make(map[phase0.ValidatorIndex][32]byte)
	contributionMap = make(map[phase0.ValidatorIndex][][32]byte)
	beaconObjects = make(map[phase0.ValidatorIndex]map[[32]byte]interface{})

	consensusData := &types.AggregatorCommitteeConsensusData{}
	if err := consensusData.Decode(r.BaseRunner.State.DecidedValue); err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not decode consensus data")
	}

	epoch := r.beacon.GetBeaconNetwork().EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.DutySlot())

	aggregateAndProofs, err := consensusData.GetAggregateAndProofs()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not get aggregate and proofs")
	}

	for i, aggregateAndProof := range aggregateAndProofs {
		validatorIndex := consensusData.Aggregators[i].ValidatorIndex
		hashRoot, err := types.GetAggregateAndProofHashRoot(aggregateAndProof)
		if err != nil {
			continue
		}

		// Calculate signing root for aggregate and proof
		domain, err := r.beacon.DomainData(epoch, types.DomainAggregateAndProof)
		if err != nil {
			continue
		}

		root, err := types.ComputeETHSigningRoot(hashRoot, domain)
		if err != nil {
			continue
		}

		aggregatorMap[validatorIndex] = root

		// Store beacon object
		if _, ok := beaconObjects[validatorIndex]; !ok {
			beaconObjects[validatorIndex] = make(map[[32]byte]interface{})
		}
		beaconObjects[validatorIndex][root] = aggregateAndProof
	}

	contributions, err := consensusData.GetSyncCommitteeContributions()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not get sync committee contributions")
	}
	for i, contribution := range contributions {
		validatorIndex := consensusData.Contributors[i].ValidatorIndex

		// Create contribution and proof
		contribAndProof := &altair.ContributionAndProof{
			AggregatorIndex: validatorIndex,
			Contribution:    &contribution.Contribution,
			SelectionProof:  consensusData.Contributors[i].SelectionProof,
		}

		// Calculate signing root
		domain, err := r.beacon.DomainData(epoch, types.DomainContributionAndProof)
		if err != nil {
			continue
		}

		root, err := types.ComputeETHSigningRoot(contribAndProof, domain)
		if err != nil {
			continue
		}

		contributionMap[validatorIndex] = append(contributionMap[validatorIndex], root)

		// Store beacon object
		if _, ok := beaconObjects[validatorIndex]; !ok {
			beaconObjects[validatorIndex] = make(map[[32]byte]interface{})
		}
		beaconObjects[validatorIndex][root] = contribAndProof
	}

	return aggregatorMap, contributionMap, beaconObjects, nil
}

// findValidatorsForRoot finds which validators are associated with a given root
// preConsensusMetadata holds metadata for pre-consensus processing

// Note: findValidatorsForPreConsensusRoot and findValidatorsForPostConsensusRoot are similar but serve different purposes:
// - findValidatorsForPreConsensusRoot: Returns detailed metadata including sync committee indices (pre-consensus)
// - findValidatorsForPostConsensusRoot: Returns just the role and validator list (post-consensus)
type preConsensusMetadata struct {
	ValidatorIndex              phase0.ValidatorIndex
	Role                        types.BeaconRole
	ValidatorSyncCommitteeIndex types.ValidatorSyncCommitteeIndex // only for sync committee role
}

// findValidatorsForPreConsensusRoot finds all validators that have the given root in pre-consensus
func findValidatorsForPreConsensusRoot(
	expectedRoot [32]byte,
	aggregatorMap map[phase0.ValidatorIndex][32]byte,
	contributionMap map[phase0.ValidatorIndex]map[types.ValidatorSyncCommitteeIndex][32]byte,
) ([]preConsensusMetadata, bool) {
	var metadata []preConsensusMetadata

	// Check aggregator map
	for validator, root := range aggregatorMap {
		if root == expectedRoot {
			metadata = append(metadata, preConsensusMetadata{
				ValidatorIndex: validator,
				Role:           types.BNRoleAggregator,
			})
		}
	}

	// Check sync committee contribution map
	for validator, indexMap := range contributionMap {
		for index, root := range indexMap {
			if root == expectedRoot {
				metadata = append(metadata, preConsensusMetadata{
					ValidatorIndex:              validator,
					Role:                        types.BNRoleSyncCommitteeContribution,
					ValidatorSyncCommitteeIndex: index,
				})
			}
		}
	}

	return metadata, len(metadata) > 0
}

func findValidatorsForPostConsensusRoot(
	expectedRoot [32]byte,
	aggregatorMap map[phase0.ValidatorIndex][32]byte,
	contributionMap map[phase0.ValidatorIndex][][32]byte,
) (types.BeaconRole, []phase0.ValidatorIndex, bool) {
	var validators []phase0.ValidatorIndex

	// Check aggregator map
	for validator, root := range aggregatorMap {
		if root == expectedRoot {
			validators = append(validators, validator)
		}
	}
	if len(validators) > 0 {
		return types.BNRoleAggregator, validators, true
	}

	// Check contribution map
	for validator, roots := range contributionMap {
		for _, root := range roots {
			if root == expectedRoot {
				validators = append(validators, validator)
				break
			}
		}
	}
	if len(validators) > 0 {
		return types.BNRoleSyncCommitteeContribution, validators, true
	}

	return types.BNRoleUnknown, nil, false
}

// constructSignedAggregateAndProof constructs a signed aggregate and proof from versioned data
func (r *AggregatorCommitteeRunner) constructSignedAggregateAndProof(
	aggregateAndProof *spec.VersionedAggregateAndProof,
	signature phase0.BLSSignature,
) (*spec.VersionedSignedAggregateAndProof, error) {
	ret := &spec.VersionedSignedAggregateAndProof{
		Version: aggregateAndProof.Version,
	}

	switch ret.Version {
	case spec.DataVersionPhase0:
		ret.Phase0 = &phase0.SignedAggregateAndProof{
			Message:   aggregateAndProof.Phase0,
			Signature: signature,
		}
	case spec.DataVersionAltair:
		ret.Altair = &phase0.SignedAggregateAndProof{
			Message:   aggregateAndProof.Altair,
			Signature: signature,
		}
	case spec.DataVersionBellatrix:
		ret.Bellatrix = &phase0.SignedAggregateAndProof{
			Message:   aggregateAndProof.Bellatrix,
			Signature: signature,
		}
	case spec.DataVersionCapella:
		ret.Capella = &phase0.SignedAggregateAndProof{
			Message:   aggregateAndProof.Capella,
			Signature: signature,
		}
	case spec.DataVersionDeneb:
		ret.Deneb = &phase0.SignedAggregateAndProof{
			Message:   aggregateAndProof.Deneb,
			Signature: signature,
		}
	case spec.DataVersionElectra:
		if aggregateAndProof.Electra == nil {
			return nil, errors.New("nil Electra aggregate and proof")
		}
		ret.Electra = &electra.SignedAggregateAndProof{
			Message:   aggregateAndProof.Electra,
			Signature: signature,
		}
	case spec.DataVersionFulu:
		if aggregateAndProof.Fulu == nil {
			return nil, errors.New("nil Fulu aggregate and proof")
		}
		ret.Fulu = &electra.SignedAggregateAndProof{
			Message:   aggregateAndProof.Fulu,
			Signature: signature,
		}

	default:
		return nil, errors.Errorf("unknown version %s", ret.Version.String())
	}

	return ret, nil
}

func (r *AggregatorCommitteeRunner) RecordSubmission(role types.BeaconRole, validatorIndex phase0.ValidatorIndex, root [32]byte) {
	if _, ok := r.submittedDuties[role]; !ok {
		r.submittedDuties[role] = make(map[phase0.ValidatorIndex]map[[32]byte]struct{})
	}
	if _, ok := r.submittedDuties[role][validatorIndex]; !ok {
		r.submittedDuties[role][validatorIndex] = make(map[[32]byte]struct{})
	}
	r.submittedDuties[role][validatorIndex][root] = struct{}{}
}

// HasSubmitted checks if a validator has already submitted for a given role and root
func (r *AggregatorCommitteeRunner) HasSubmitted(role types.BeaconRole, validatorIndex phase0.ValidatorIndex, root [32]byte) bool {
	if _, ok := r.submittedDuties[role]; !ok {
		return false
	}
	if _, ok := r.submittedDuties[role][validatorIndex]; !ok {
		return false
	}
	_, submitted := r.submittedDuties[role][validatorIndex][root]
	return submitted
}

// HasSubmittedAllDuties checks if all expected duties have been submitted.
// For aggregator role we expect exactly one submission per validator.
// For sync committee contribution role we expect one submission per expected root
// (i.e., per subcommittee index assigned to that validator for this slot).
func (r *AggregatorCommitteeRunner) HasSubmittedAllDuties() bool {
	duty := r.BaseRunner.State.StartingDuty.(*types.AggregatorCommitteeDuty)

	// Build the expected post-consensus roots per validator/role from the decided data.
	aggregatorMap, contributionMap, _, err := r.expectedPostConsensusRootsAndBeaconObjects()
	if err != nil {
		// If we can't resolve the expected set, do not finish yet.
		return false
	}

	for _, vDuty := range duty.ValidatorDuties {
		if vDuty == nil {
			continue
		}

		// Only consider validators this operator actually runs.
		if _, hasShare := r.BaseRunner.Share[vDuty.ValidatorIndex]; !hasShare {
			continue
		}

		switch vDuty.Type {
		case types.BNRoleAggregator:
			// Expect exactly one aggregate root for this validator.
			expectedRoot, ok := aggregatorMap[vDuty.ValidatorIndex]
			if !ok {
				// If consensus did not include this validator's aggregate, we haven't finished.
				return false
			}
			if !r.HasSubmitted(types.BNRoleAggregator, vDuty.ValidatorIndex, expectedRoot) {
				return false
			}

		case types.BNRoleSyncCommitteeContribution:
			// Expect a submission for every contribution root assigned to this validator.
			expectedRoots, ok := contributionMap[vDuty.ValidatorIndex]
			if !ok || len(expectedRoots) == 0 {
				// The duty indicates sync committee work but no expected roots were found.
				return false
			}
			for _, root := range expectedRoots {
				if !r.HasSubmitted(types.BNRoleSyncCommitteeContribution, vDuty.ValidatorIndex, root) {
					return false
				}
			}

		default:
			// Unknown role type: don't allow finishing.
			return false
		}
	}

	return true
}

// HasSeenSignerInPreConsensus checks if we've already seen a pre-consensus message from such signer
func (r *AggregatorCommitteeRunner) HasSeenSignerInPreConsensus(signedMsg *types.PartialSignatureMessages) bool {
	for _, m := range signedMsg.Messages {
		if _, seen := r.preConsensusSeenSigners[m.Signer]; seen {
			return true
		}
	}
	return false
}

// MarkPreConsensusSignerAsSeen marks the signer of the given pre-consensus message as seen
func (r *AggregatorCommitteeRunner) MarkPreConsensusSignerAsSeen(signedMsg *types.PartialSignatureMessages) {
	for _, m := range signedMsg.Messages {
		r.preConsensusSeenSigners[m.Signer] = struct{}{}
		break // all message signers are equal, no need to check further
	}
}

// HasSeenAllPreConsensusSigners checks if we've seen messages from all committee operators
func (r *AggregatorCommitteeRunner) HasSeenAllPreConsensusSigners() bool {
	committeeSize := len(r.BaseRunner.QBFTController.CommitteeMember.Committee)
	return len(r.preConsensusSeenSigners) >= committeeSize
}

// MarkSelectionChecked marks a given (validator index, root) selection as checked in pre-consensus
func (r *AggregatorCommitteeRunner) MarkSelectionChecked(validatorIndex phase0.ValidatorIndex, root [32]byte) {
	if _, ok := r.preConsensusSelectionsTested[validatorIndex]; !ok {
		r.preConsensusSelectionsTested[validatorIndex] = make(map[[32]byte]struct{})
	}
	r.preConsensusSelectionsTested[validatorIndex][root] = struct{}{}
}

// AllSelectionsChecked checks if we have tested all possible (validator index, root) selections
func (r *AggregatorCommitteeRunner) AllSelectionsChecked(
	aggregatorMap map[phase0.ValidatorIndex][32]byte,
	contributionMap map[phase0.ValidatorIndex]map[types.ValidatorSyncCommitteeIndex][32]byte,
) bool {
	for validatorIndex, root := range aggregatorMap {
		if _, hasShare := r.BaseRunner.Share[validatorIndex]; !hasShare {
			continue
		}
		if _, ok := r.preConsensusSelectionsTested[validatorIndex]; !ok {
			return false
		}
		if _, ok := r.preConsensusSelectionsTested[validatorIndex][root]; !ok {
			return false
		}
	}
	for validatorIndex, syncCommitteeRoots := range contributionMap {
		if _, hasShare := r.BaseRunner.Share[validatorIndex]; !hasShare {
			continue
		}
		if _, ok := r.preConsensusSelectionsTested[validatorIndex]; !ok {
			return false
		}
		for _, root := range syncCommitteeRoots {
			if _, ok := r.preConsensusSelectionsTested[validatorIndex][root]; !ok {
				return false
			}
		}
	}
	return true
}
