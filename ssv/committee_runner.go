package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

type CommitteeRunner struct {
	BaseRunner      *BaseRunner
	beacon          BeaconNode
	network         Network
	signer          types.BeaconSigner
	operatorSigner  *types.OperatorSigner
	valCheck        qbft.ProposedValueCheckF
	submittedDuties map[types.BeaconRole]map[phase0.ValidatorIndex]struct{}
}

func NewCommitteeRunner(beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner *types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
) (Runner, error) {
	if len(share) == 0 {
		return nil, fmt.Errorf("no shares")
	}
	if err := validateShareMap(share); err != nil {
		return nil, err
	}
	return &CommitteeRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType: types.RoleCommittee,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
			QBFTController: qbftController,
		},
		beacon:          beacon,
		network:         network,
		signer:          signer,
		operatorSigner:  operatorSigner,
		valCheck:        valCheck,
		submittedDuties: make(map[types.BeaconRole]map[phase0.ValidatorIndex]struct{}),
	}, nil
}

func (cr CommitteeRunner) StartNewDuty(duty types.Duty, quorum uint64) error {
	err := cr.BaseRunner.baseStartNewDuty(cr, duty, quorum)
	if err != nil {
		return err
	}
	cr.submittedDuties[types.BNRoleAttester] = make(map[phase0.ValidatorIndex]struct{})
	cr.submittedDuties[types.BNRoleSyncCommittee] = make(map[phase0.ValidatorIndex]struct{})
	return nil
}

func (cr CommitteeRunner) GetBaseRunner() *BaseRunner {
	return cr.BaseRunner
}

func (cr CommitteeRunner) GetBeaconNode() BeaconNode {
	return cr.beacon
}

func (cr CommitteeRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return cr.valCheck
}

func (cr CommitteeRunner) GetNetwork() Network {
	return cr.network
}

func (cr CommitteeRunner) HasRunningDuty() bool {
	return cr.BaseRunner.hasRunningDuty()
}

func (cr CommitteeRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	return fmt.Errorf("no pre consensus phase for committee runner")
}

func (cr CommitteeRunner) ProcessConsensus(msg *types.SignedSSVMessage) error {
	// Decode the decided value into the slot's fork prototype: on Gloas slots it is a GloasBeaconVote
	// carrying the BN-supplied attestation index (SIP #94 §2), before Gloas a plain BeaconVote. Both
	// implement types.Encoder, so the QBFT plumbing is identical. With no state there is no duty to read
	// the fork from, but consensus processing also stops before the prototype is used, so the pre-Gloas
	// default is inert rather than a fallback anything relies on.
	var votePrototype types.Encoder = &types.BeaconVote{}
	if state := cr.BaseRunner.State; state != nil {
		votePrototype = committeeVoteForSlot(cr.beacon, state.StartingDuty.DutySlot())
	}

	decided, decidedValue, err := cr.BaseRunner.baseConsensusMsgProcessing(cr, msg, votePrototype)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	duty := cr.BaseRunner.State.StartingDuty
	version := versionForSlot(cr.beacon, duty.DutySlot())
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     duty.DutySlot(),
		Messages: []*types.PartialSignatureMessage{},
	}

	beaconVote, gloasIndex, err := baseVoteAndIndex(decidedValue)
	if err != nil {
		return err
	}
	for _, validatorDuty := range duty.(*types.CommitteeDuty).ValidatorDuties {
		switch validatorDuty.Type {
		case types.BNRoleAttester:
			attestationData := constructAttestationData(beaconVote, validatorDuty, version, gloasIndex)
			partialMsg, err := cr.BaseRunner.signBeaconObject(cr, validatorDuty, attestationData, validatorDuty.DutySlot(),
				types.DomainAttester)
			if err != nil {
				return errors.Wrap(err, "failed signing attestation data")
			}
			postConsensusMsg.Messages = append(postConsensusMsg.Messages, partialMsg)

		case types.BNRoleSyncCommittee:
			blockRoot := beaconVote.BlockRoot
			partialMsg, err := cr.BaseRunner.signBeaconObject(cr, validatorDuty, types.SSZBytes(blockRoot[:]), validatorDuty.DutySlot(),
				types.DomainSyncCommittee)
			if err != nil {
				return errors.Wrap(err, "failed signing sync committee message")
			}
			postConsensusMsg.Messages = append(postConsensusMsg.Messages, partialMsg)
		}
	}

	committeeMember := cr.BaseRunner.QBFTController.CommitteeMember
	operatorIDs := make([]types.OperatorID, len(committeeMember.Committee))
	for i, operator := range committeeMember.Committee {
		operatorIDs[i] = operator.OperatorID
	}
	committeeID := types.GetCommitteeID(operatorIDs)
	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewCommitteeMsgID(committeeMember.DomainType, committeeID, cr.BaseRunner.RunnerRoleType),
	}
	ssvMsg.Data, err = postConsensusMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	sig, err := cr.operatorSigner.SignSSVMessage(ssvMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign SSVMessage")
	}

	msgToBroadcast := &types.SignedSSVMessage{
		Signatures:  [][]byte{sig},
		OperatorIDs: []types.OperatorID{cr.BaseRunner.QBFTController.CommitteeMember.OperatorID},
		SSVMessage:  ssvMsg,
	}

	if err := cr.GetNetwork().Broadcast(ssvMsg.MsgID, msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil

}

// ProcessEnvelopeDissemination returns an error: only the envelope-proposer runner processes
// disseminated envelopes (SIP #94 §6).
func (CommitteeRunner) ProcessEnvelopeDissemination(*types.SignedSSVMessage) error {
	return types.NewError(types.EnvelopeDisseminationUnsupportedErrorCode, "runner does not process envelope dissemination")
}

func (cr CommitteeRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	// Gets all the roots that received a quorum of signatures
	quorum, rootsList, err := cr.BaseRunner.basePostConsensusMsgProcessing(&cr, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	// Get unique roots to avoid repetition
	rootSet := make(map[[32]byte]struct{})
	for _, root := range rootsList {
		rootSet[root] = struct{}{}
	}

	// Get validator-root maps for attestations and sync committees, and the root-beacon object map
	attestationMap, committeeMap, beaconObjects, err := cr.expectedPostConsensusRootsAndBeaconObjects()
	if err != nil {
		return errors.Wrap(err, "could not get expected post consensus roots and beacon objects")
	}

	var anyErr error
	attestationsToSubmit := make(map[phase0.ValidatorIndex]*spec.VersionedAttestation)
	syncCommitteeMessagesToSubmit := make(map[phase0.ValidatorIndex]*altair.SyncCommitteeMessage)

	// For each root that got at least one quorum, find the duties associated to it and try to submit
	for root := range rootSet {

		// Get validators related to the given root
		role, validators, found := findValidators(root, attestationMap, committeeMap)

		if !found {
			// Edge case: since operators may have divergent sets of validators,
			// it's possible that an operator doesn't have the validator associated to a root.
			// In this case, we simply continue.
			continue
		}

		for _, validator := range validators {

			// Skip if no quorum - We know that a root has quorum but not necessarily for the validator
			if !cr.BaseRunner.State.PostConsensusContainer.HasQuorum(validator, root) {
				continue
			}
			// Skip if already submitted
			if cr.HasSubmitted(role, validator) {
				continue
			}

			// Reconstruct signature
			share := cr.BaseRunner.Share[validator]
			pubKey := share.ValidatorPubKey
			sig, err := cr.BaseRunner.State.ReconstructBeaconSig(cr.BaseRunner.State.PostConsensusContainer, root,
				pubKey[:], validator)
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			if err != nil {
				for root := range rootSet {
					cr.BaseRunner.FallBackAndVerifyEachSignature(cr.BaseRunner.State.PostConsensusContainer, root,
						share.Committee, validator)
				}
				// Record the error and continue to next validators
				anyErr = errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
				continue
			}
			specSig := phase0.BLSSignature{}
			copy(specSig[:], sig)

			// Get the beacon object related to root
			validatorObjs, exists := beaconObjects[validator]
			if !exists {
				anyErr = errors.Wrap(err, "could not find beacon object for validator")
				continue
			}
			sszObject, exists := validatorObjs[root]
			if !exists {
				anyErr = errors.Wrap(err, "could not find beacon object for validator")
				continue
			}

			// Store objects for multiple submission
			if role == types.BNRoleSyncCommittee {
				syncMsg := sszObject.(*altair.SyncCommitteeMessage)
				// Insert signature
				syncMsg.Signature = specSig

				syncCommitteeMessagesToSubmit[validator] = syncMsg

			} else if role == types.BNRoleAttester {
				att := sszObject.(*spec.VersionedAttestation)
				// Insert signature
				att, err = VersionedAttestationWithSignature(att, specSig)
				if err != nil {
					anyErr = errors.Wrap(err, "could not insert signature in versioned attestation")
					continue
				}

				attestationsToSubmit[validator] = att
			}
		}
	}

	// Submit multiple attestations
	attestations := make([]*spec.VersionedAttestation, 0, len(attestationsToSubmit))
	for _, att := range attestationsToSubmit {
		attestations = append(attestations, att)
	}

	if len(attestations) > 0 {
		if err := cr.beacon.SubmitAttestations(attestations); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed attestation")
		}
		// Record successful submissions
		for validator := range attestationsToSubmit {
			cr.RecordSubmission(types.BNRoleAttester, validator)
		}
	}

	// Submit multiple sync committee
	syncCommitteeMessages := make([]*altair.SyncCommitteeMessage, 0, len(syncCommitteeMessagesToSubmit))
	for _, syncMsg := range syncCommitteeMessagesToSubmit {
		syncCommitteeMessages = append(syncCommitteeMessages, syncMsg)
	}
	if len(syncCommitteeMessages) > 0 {
		if err := cr.beacon.SubmitSyncMessages(syncCommitteeMessages); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed sync committee")
		}
		// Record successful submissions
		for validator := range syncCommitteeMessagesToSubmit {
			cr.RecordSubmission(types.BNRoleSyncCommittee, validator)
		}
	}

	if anyErr != nil {
		return anyErr
	}

	// Check if duty has terminated (runner has submitted for all duties)
	if cr.HasSubmittedAllValidatorDuties(attestationMap, committeeMap) {
		cr.BaseRunner.State.Finished = true
	}
	return nil
}

// Returns true if the runner has done submissions for all validators for the given slot
func (cr *CommitteeRunner) HasSubmittedAllValidatorDuties(attestationMap map[phase0.ValidatorIndex][32]byte, syncCommitteeMap map[phase0.ValidatorIndex][32]byte) bool {
	// Expected total
	expectedTotalSubmissions := len(attestationMap) + len(syncCommitteeMap)

	totalSubmissions := 0

	// Add submitted attestation duties
	for valIdx := range attestationMap {
		if cr.HasSubmitted(types.BNRoleAttester, valIdx) {
			totalSubmissions++
		}
	}
	// Add submitted sync committee duties
	for valIdx := range syncCommitteeMap {
		if cr.HasSubmitted(types.BNRoleSyncCommittee, valIdx) {
			totalSubmissions++
		}
	}
	return totalSubmissions >= expectedTotalSubmissions
}

// Records a submission for the (role, validator index, slot) tuple
func (cr *CommitteeRunner) RecordSubmission(role types.BeaconRole, valIdx phase0.ValidatorIndex) {
	if _, ok := cr.submittedDuties[role]; !ok {
		cr.submittedDuties[role] = make(map[phase0.ValidatorIndex]struct{})
	}
	cr.submittedDuties[role][valIdx] = struct{}{}
}

// Returns true if there is a record of submission for the (role, validator index, slot) tuple
func (cr *CommitteeRunner) HasSubmitted(role types.BeaconRole, valIdx phase0.ValidatorIndex) bool {
	if _, ok := cr.submittedDuties[role]; !ok {
		return false
	}
	_, ok := cr.submittedDuties[role][valIdx]
	return ok
}

func findValidators(
	expectedRoot [32]byte,
	attestationMap map[phase0.ValidatorIndex][32]byte,
	committeeMap map[phase0.ValidatorIndex][32]byte) (types.BeaconRole, []phase0.ValidatorIndex, bool) {
	var validators []phase0.ValidatorIndex

	// look for the expectedRoot in attestationMap
	for validator, root := range attestationMap {
		if root == expectedRoot {
			validators = append(validators, validator)
		}
	}
	if len(validators) > 0 {
		return types.BNRoleAttester, validators, true
	}
	// look for the expectedRoot in committeeMap
	for validator, root := range committeeMap {
		if root == expectedRoot {
			validators = append(validators, validator)
		}
	}
	if len(validators) > 0 {
		return types.BNRoleSyncCommittee, validators, true
	}
	return types.BNRoleUnknown, nil, false
}

// Unneeded since no preconsensus phase
func (cr CommitteeRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	panic("not in use")
}

// This function signature returns only one domain type... but we can have mixed domains
// instead we rely on expectedPostConsensusRootsAndBeaconObjects that is called later
func (cr CommitteeRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	panic("not in use")
}

// expectedPostConsensusRootsAndBeaconObjects returns the expected roots and beacon objects for the post consensus
// phase. It returns the attestation and sync committee validator to root map, as well as a root to beacon object map.
func (cr *CommitteeRunner) expectedPostConsensusRootsAndBeaconObjects() (
	attestationMap map[phase0.ValidatorIndex][32]byte,
	syncCommitteeMap map[phase0.ValidatorIndex][32]byte,
	beaconObjects map[phase0.ValidatorIndex]map[[32]byte]interface{}, error error,
) {
	attestationMap = make(map[phase0.ValidatorIndex][32]byte)
	syncCommitteeMap = make(map[phase0.ValidatorIndex][32]byte)
	beaconObjects = make(map[phase0.ValidatorIndex]map[[32]byte]interface{})
	duty := cr.BaseRunner.State.StartingDuty.(*types.CommitteeDuty)
	beaconVoteData := cr.BaseRunner.State.DecidedValue

	slot := duty.DutySlot()
	epoch := cr.beacon.GetBeaconNetwork().EstimatedEpochAtSlot(slot)

	// Decode the decided value into the slot's fork prototype, then split into the base vote and the
	// optional Gloas attestation index (SIP #94 §2).
	decodedVote := committeeVoteForSlot(cr.beacon, slot)
	dataVersion := versionForSlot(cr.beacon, slot)
	if err := decodedVote.Decode(beaconVoteData); err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not decode beacon vote")
	}
	beaconVote, gloasIndex, err := baseVoteAndIndex(decodedVote)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := beaconVote.Validate(); err != nil {
		return nil, nil, nil, errors.Wrap(err, "invalid beacon vote")
	}

	for _, validatorDuty := range duty.ValidatorDuties {
		if validatorDuty == nil {
			continue
		}

		switch validatorDuty.Type {
		case types.BNRoleAttester:

			// Attestation object
			attestationData := constructAttestationData(beaconVote, validatorDuty, dataVersion, gloasIndex)
			attestationResponse, err := ConstructVersionedAttestationWithoutSignature(attestationData, dataVersion, validatorDuty)
			if err != nil {
				continue
			}

			// Root
			domain, err := cr.GetBeaconNode().DomainData(epoch, types.DomainAttester)
			if err != nil {
				continue
			}
			root, err := types.ComputeETHSigningRoot(attestationData, domain)
			if err != nil {
				continue
			}

			// Add to map
			attestationMap[validatorDuty.ValidatorIndex] = root
			if _, ok := beaconObjects[validatorDuty.ValidatorIndex]; !ok {
				beaconObjects[validatorDuty.ValidatorIndex] = make(map[[32]byte]interface{})
			}
			beaconObjects[validatorDuty.ValidatorIndex][root] = attestationResponse
		case types.BNRoleSyncCommittee:
			// Sync committee beacon object
			syncMsg := &altair.SyncCommitteeMessage{
				Slot:            slot,
				BeaconBlockRoot: beaconVote.BlockRoot,
				ValidatorIndex:  validatorDuty.ValidatorIndex,
			}

			// Root
			domain, err := cr.GetBeaconNode().DomainData(epoch, types.DomainSyncCommittee)
			if err != nil {
				continue
			}
			// Eth root
			blockRoot := types.SSZBytes(beaconVote.BlockRoot[:])
			root, err := types.ComputeETHSigningRoot(blockRoot, domain)
			if err != nil {
				continue
			}

			// Set root and beacon object
			syncCommitteeMap[validatorDuty.ValidatorIndex] = root
			if _, ok := beaconObjects[validatorDuty.ValidatorIndex]; !ok {
				beaconObjects[validatorDuty.ValidatorIndex] = make(map[[32]byte]interface{})
			}
			beaconObjects[validatorDuty.ValidatorIndex][root] = syncMsg
		}
	}
	return attestationMap, syncCommitteeMap, beaconObjects, nil
}

func (cr CommitteeRunner) executeDuty(duty types.Duty) error {
	slot := duty.DutySlot()
	attData, _, err := cr.GetBeaconNode().GetAttestationData(slot)
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}

	// On Gloas slots the consensus value is a GloasBeaconVote carrying the BN-supplied attestation
	// index (SIP #94 §2); before Gloas it is a plain BeaconVote.
	var input types.Encoder
	if versionForSlot(cr.beacon, slot) >= gloas.DataVersionGloas {
		input = &types.GloasBeaconVote{
			BlockRoot:            attData.BeaconBlockRoot,
			Source:               attData.Source,
			Target:               attData.Target,
			AttestationDataIndex: attData.Index,
		}
	} else {
		input = &types.BeaconVote{
			BlockRoot: attData.BeaconBlockRoot,
			Source:    attData.Source,
			Target:    attData.Target,
		}
	}

	if err := cr.BaseRunner.decide(cr, slot, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (cr CommitteeRunner) GetSigner() types.BeaconSigner {
	return cr.signer
}

func (cr CommitteeRunner) GetOperatorSigner() *types.OperatorSigner {
	return cr.operatorSigner
}

// versionForSlot returns the data version the duty at slot runs under, per the local beacon node's
// fork schedule.
//
// Comparisons of this version against gloas.DataVersionGloas are >=, not ==, because any version at
// or beyond Gloas is Gloas or a later fork, all of which keep the Gloas-shaped values until a later
// fork's spec changes them. Contrast the version stamped *inside* a received consensus-data value
// (e.g. ProposerConsensusData.Version), which is attacker-controlled and therefore pinned against
// the slot's fork where it is consumed (see ProposerValueCheckF).
func versionForSlot(beacon BeaconNode, slot phase0.Slot) spec.DataVersion {
	return beacon.DataVersion(beacon.GetBeaconNetwork().EstimatedEpochAtSlot(slot))
}

// committeeVoteForSlot returns the empty committee consensus value to decode into for the slot's fork:
// a GloasBeaconVote from Gloas on (SIP #94 §2), a BeaconVote before. Every committee path derives the
// fork through this and versionForSlot so the decisions cannot drift apart.
func committeeVoteForSlot(beacon BeaconNode, slot phase0.Slot) types.Encoder {
	if versionForSlot(beacon, slot) >= gloas.DataVersionGloas {
		return &types.GloasBeaconVote{}
	}
	return &types.BeaconVote{}
}

// baseVoteAndIndex splits a decided committee value into its base BeaconVote fields and, on Gloas,
// the BN-supplied attestation index (SIP #94 §2). Pre-Gloas the value is already a *types.BeaconVote
// and the returned index is nil.
func baseVoteAndIndex(decidedValue interface{}) (*types.BeaconVote, *phase0.CommitteeIndex, error) {
	switch vote := decidedValue.(type) {
	case *types.GloasBeaconVote:
		index := vote.AttestationDataIndex
		return &types.BeaconVote{BlockRoot: vote.BlockRoot, Source: vote.Source, Target: vote.Target}, &index, nil
	case *types.BeaconVote:
		return vote, nil, nil
	default:
		// Unreachable: the caller selects the prototype. Fail loudly rather than returning a nil vote
		// that would nil-deref inside constructAttestationData.
		return nil, nil, fmt.Errorf("unexpected committee decided value type %T", decidedValue)
	}
}

func constructAttestationData(vote *types.BeaconVote, duty *types.ValidatorDuty, version spec.DataVersion, gloasIndex *phase0.CommitteeIndex) *phase0.AttestationData {
	attData := &phase0.AttestationData{
		Slot:            duty.Slot,
		Index:           duty.CommitteeIndex,
		BeaconBlockRoot: vote.BlockRoot,
		Source:          vote.Source,
		Target:          vote.Target,
	}

	switch {
	case gloasIndex != nil:
		// SIP #94 §2: under Gloas the index is the decided payload-status value (0/1), not a committee index.
		attData.Index = *gloasIndex
	case version >= spec.DataVersionElectra:
		attData.Index = 0 // EIP-7549: Index should be set to 0
	}

	return attData
}

func VersionedAttestationWithSignature(att *spec.VersionedAttestation, specSig phase0.BLSSignature) (*spec.VersionedAttestation, error) {
	switch att.Version {
	case spec.DataVersionPhase0:
		if att.Phase0 == nil {
			return att, fmt.Errorf("no Phase0 attestation")
		}
		att.Phase0.Signature = specSig
	case spec.DataVersionAltair:
		if att.Altair == nil {
			return att, fmt.Errorf("no Altair attestation")
		}
		att.Altair.Signature = specSig
	case spec.DataVersionBellatrix:
		if att.Bellatrix == nil {
			return att, fmt.Errorf("no Bellatrix attestation")
		}
		att.Bellatrix.Signature = specSig
	case spec.DataVersionCapella:
		if att.Capella == nil {
			return att, fmt.Errorf("no Capella attestation")
		}
		att.Capella.Signature = specSig
	case spec.DataVersionDeneb:
		if att.Deneb == nil {
			return att, fmt.Errorf("no Deneb attestation")
		}
		att.Deneb.Signature = specSig
	case spec.DataVersionElectra:
		if att.Electra == nil {
			return att, fmt.Errorf("no Electra attestation")
		}
		att.Electra.Signature = specSig
	case spec.DataVersionFulu:
		if att.Fulu == nil {
			return att, fmt.Errorf("no Fulu attestation")
		}
		att.Fulu.Signature = specSig
	case gloas.DataVersionGloas:
		// spec.VersionedAttestation has no Gloas field; Gloas reuses the Electra container shape
		// (local API wrapper, not wire format — SIP #94 §2), distinguished by the Version tag.
		if att.Electra == nil {
			return att, fmt.Errorf("no Electra attestation")
		}
		att.Electra.Signature = specSig
	default:
		return nil, fmt.Errorf("unknown version: %s", att.Version)
	}

	return att, nil
}

func ConstructPhase0AttestationWithoutSignature(attestationData *phase0.AttestationData, validatorDuty *types.ValidatorDuty) *phase0.Attestation {
	aggregationBitfield := bitfield.NewBitlist(validatorDuty.CommitteeLength)
	aggregationBitfield.SetBitAt(validatorDuty.ValidatorCommitteeIndex, true)
	return &phase0.Attestation{
		Data:            attestationData,
		AggregationBits: aggregationBitfield,
	}
}

func ConstructElectraAttestationWithoutSignature(attestationData *phase0.AttestationData, validatorDuty *types.ValidatorDuty) *electra.Attestation {
	aggregationBitfield := bitfield.NewBitlist(validatorDuty.CommitteeLength)
	aggregationBitfield.SetBitAt(validatorDuty.ValidatorCommitteeIndex, true)

	committeeBits := bitfield.NewBitvector64()
	committeeBits.SetBitAt(uint64(validatorDuty.CommitteeIndex), true)

	return &electra.Attestation{
		Data:            attestationData,
		AggregationBits: aggregationBitfield,
		CommitteeBits:   committeeBits,
	}
}

func ConstructVersionedAttestationWithoutSignature(attestationData *phase0.AttestationData, dataVersion spec.DataVersion, validatorDuty *types.ValidatorDuty) (*spec.VersionedAttestation, error) {
	ret := &spec.VersionedAttestation{
		Version:        dataVersion,
		ValidatorIndex: &validatorDuty.ValidatorIndex,
	}

	switch dataVersion {
	case spec.DataVersionPhase0:
		ret.Phase0 = ConstructPhase0AttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case spec.DataVersionAltair:
		ret.Altair = ConstructPhase0AttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case spec.DataVersionBellatrix:
		ret.Bellatrix = ConstructPhase0AttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case spec.DataVersionCapella:
		ret.Capella = ConstructPhase0AttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case spec.DataVersionDeneb:
		ret.Deneb = ConstructPhase0AttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case spec.DataVersionElectra:
		ret.Electra = ConstructElectraAttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case spec.DataVersionFulu:
		ret.Fulu = ConstructElectraAttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	case gloas.DataVersionGloas:
		// Gloas reuses the Electra attestation container (SIP #94 §2); the Version tag distinguishes it.
		ret.Electra = ConstructElectraAttestationWithoutSignature(attestationData, validatorDuty)
		return ret, nil
	default:
		return nil, fmt.Errorf("unknown version")
	}
}
