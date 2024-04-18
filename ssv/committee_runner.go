package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"
)

type CommitteeRunner struct {
	BaseRunner     *BaseRunner
	beacon         BeaconNode
	network        Network
	signer         types.BeaconSigner
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewCommitteeRunner(beaconNetwork types.BeaconNetwork,
	share map[phase0.ValidatorIndex]*types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.BeaconSigner,
	operatorSigner types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot) Runner {
	return &CommitteeRunner{
		BaseRunner: &BaseRunner{
			RunnerRoleType:     RoleCommittee,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},
		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
	}
}

func (cr CommitteeRunner) StartNewDuty(duty types.Duty) error {
	return cr.BaseRunner.baseStartNewDuty(cr, duty)
}

func (cr CommitteeRunner) Encode() ([]byte, error) {
	return json.Marshal(cr)
}

// StopDuty stops the duty for the given validator
func (cr *CommitteeRunner) StopDuty(validator types.ValidatorPK) {
	for _, duty := range cr.BaseRunner.State.StartingDuty.(*types.CommitteeDuty).BeaconDuties {
		if types.ValidatorPK(duty.PubKey) == validator {
			duty.IsStopped = true
		}
	}
}

func (cr CommitteeRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &cr)
}

func (cr CommitteeRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := cr.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

func (cr CommitteeRunner) GetBaseRunner() *BaseRunner {
	//TODO implement me
	panic("implement me")
}

func (cr CommitteeRunner) GetBeaconNode() BeaconNode {
	//TODO implement me
	panic("implement me")
}

func (cr CommitteeRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return cr.valCheck
}

func (cr CommitteeRunner) GetNetwork() Network {
	//TODO implement me
	panic("implement me")
}

func (cr CommitteeRunner) HasRunningDuty() bool {
	return cr.BaseRunner.hasRunningDuty()
}

func (cr CommitteeRunner) ProcessPreConsensus(signedMsg *types.PartialSignatureMessages) error {
	//TODO implement me
	panic("implement me")
}

func (cr CommitteeRunner) ProcessConsensus(msg *types.SignedSSVMessage) error {
	decided, decidedValue, err := cr.BaseRunner.baseConsensusMsgProcessing(cr, msg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	beaconVote, err := decidedValue.GetBeaconVote()
	if err != nil {
		return errors.Wrap(err, "decided value is not a beacon vote")
	}

	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     decidedValue.Duty.Slot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for _, duty := range cr.BaseRunner.State.StartingDuty.(*types.CommitteeDuty).BeaconDuties {
		switch duty.Type {
		case types.BNRoleAttester:
			attestationData := constructAttestationData(beaconVote, duty)

			partialMsg, err := cr.BaseRunner.signBeaconObject(cr, duty, attestationData, decidedValue.Duty.Slot, types.DomainAttester)
			if err != nil {
				return errors.Wrap(err, "failed signing attestation data")
			}
			postConsensusMsg.Messages = append(postConsensusMsg.Messages, partialMsg)

		case types.BNRoleSyncCommittee:
			syncCommitteeMessage := ConstructSyncCommittee(beaconVote, duty)
			partialMsg, err := cr.BaseRunner.signBeaconObject(cr, duty, syncCommitteeMessage, decidedValue.Duty.Slot, types.DomainSyncCommittee)
			if err != nil {
				return errors.Wrap(err, "failed signing sync committee message")
			}
			postConsensusMsg.Messages = append(postConsensusMsg.Messages, partialMsg)
		}
	}

	data, err := postConsensusMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		//TODO: The Domain will be updated after new Domain PR... Will be created after this PR is merged
		MsgID: types.NewMsgID(types.GenesisMainnet, cr.GetBaseRunner().QBFTController.Share.ClusterID[:],
			cr.BaseRunner.RunnerRoleType),
		Data: data,
	}

	msgToBroadcast, err := types.SSVMessageToSignedSSVMessage(ssvMsg, cr.BaseRunner.QBFTController.Share.OperatorID,
		cr.operatorSigner.SignSSVMessage)
	if err != nil {
		return errors.Wrap(err, "could not create SignedSSVMessage from SSVMessage")
	}

	if err := cr.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil

}

// TODO finish edge case where some roots may be missing
func (cr CommitteeRunner) ProcessPostConsensus(signedMsg *types.PartialSignatureMessages) error {
	quorum, roots, err := cr.BaseRunner.basePostConsensusMsgProcessing(&cr, signedMsg)

	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}
	attestationMap, committeeMap, beaconObjects, err := cr.expectedPostConsensusRootsAndBeaconObjects()
	if err != nil {
		return errors.Wrap(err, "could not get expected post consensus roots and beacon objects")
	}
	for _, root := range roots {
		role, validators, found := findValidators(root, attestationMap, committeeMap)

		if !found {
			// TODO error?
			continue
		}

		for _, validator := range validators {
			share := cr.BaseRunner.Share[validator]
			pubKey := share.ValidatorPubKey
			sig, err := cr.BaseRunner.State.ReconstructBeaconSig(cr.BaseRunner.State.PostConsensusContainer, root,
				pubKey[:], validator)
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			// TODO should we return an error here? maybe other sigs are fine?
			if err != nil {
				for _, root := range roots {
					cr.BaseRunner.FallBackAndVerifyEachSignature(cr.BaseRunner.State.PostConsensusContainer, root,
						share.Committee, validator)
				}
				return errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
			}
			specSig := phase0.BLSSignature{}
			copy(specSig[:], sig)

			if role == types.BNRoleAttester {
				att := beaconObjects[root].(*phase0.Attestation)
				att.Signature = specSig
				// broadcast
				if err := cr.beacon.SubmitAttestation(att); err != nil {
					return errors.Wrap(err, "could not submit to Beacon chain reconstructed attestation")
				}
			} else if role == types.BNRoleSyncCommittee {
				syncMsg := beaconObjects[root].(*altair.SyncCommitteeMessage)
				syncMsg.Signature = specSig
				if err := cr.beacon.SubmitSyncMessage(syncMsg); err != nil {
					return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed sync committee")
				}
			}
		}

	}
	cr.BaseRunner.State.Finished = true
	return nil
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
			return types.BNRoleSyncCommittee, []phase0.ValidatorIndex{validator}, true
		}
	}
	return types.BNRoleUnknown, nil, false
}

func (cr CommitteeRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	//TODO implement me
	panic("implement me")
}

func (cr CommitteeRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CommitteeRunner) expectedPostConsensusRootsAndBeaconObjects() (attestationMap map[phase0.ValidatorIndex][32]byte,
	syncCommitteeMap map[phase0.ValidatorIndex][32]byte, beaconObjects map[[32]byte]ssz.HashRoot, error error) {
	attestationMap = make(map[phase0.ValidatorIndex][32]byte)
	syncCommitteeMap = make(map[phase0.ValidatorIndex][32]byte)
	duty := cr.BaseRunner.State.StartingDuty
	// TODO DecidedValue should be interface??
	// This should block if consensus is not reached
	beaconVoteData := cr.BaseRunner.State.DecidedValue
	beaconVote, err := types.NewBeaconVote(beaconVoteData)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "could not decode beacon vote")
	}
	beaconVote.Decode(beaconVoteData)
	for _, beaconDuty := range duty.(*types.CommitteeDuty).BeaconDuties {
		if beaconDuty == nil || beaconDuty.IsStopped {
			continue
		}
		switch beaconDuty.Type {
		case types.BNRoleAttester:
			attestationData := constructAttestationData(beaconVote, beaconDuty)
			aggregationBitfield := bitfield.NewBitlist(beaconDuty.CommitteeLength)
			aggregationBitfield.SetBitAt(beaconDuty.ValidatorCommitteeIndex, true)
			unSignedAtt := &phase0.Attestation{
				Data:            attestationData,
				AggregationBits: aggregationBitfield,
			}
			root, _ := attestationData.HashTreeRoot()
			attestationMap[beaconDuty.ValidatorIndex] = root
			beaconObjects[root] = unSignedAtt
		case types.BNRoleSyncCommittee:
			syncCommitteeMessage := ConstructSyncCommittee(beaconVote, beaconDuty)
			root, _ := syncCommitteeMessage.HashTreeRoot()
			syncCommitteeMap[beaconDuty.ValidatorIndex] = root
			beaconObjects[root] = syncCommitteeMessage
		}
	}
	return attestationMap, syncCommitteeMap, beaconObjects, nil
}

func (cr CommitteeRunner) executeDuty(duty types.Duty) error {

	//TODO committeeIndex is 0, is this correct?
	attData, ver, err := cr.GetBeaconNode().GetAttestationData(duty.DutySlot(), 0)
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}

	vote := types.BeaconVote{
		BlockRoot: attData.BeaconBlockRoot,
		Source:    attData.Source,
		Target:    attData.Target,
	}
	voteByts, err := vote.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal attestation data")
	}

	//TODO should duty be empty?
	input := &types.ConsensusData{
		Duty:    types.BeaconDuty{},
		Version: ver,
		DataSSZ: voteByts,
	}

	if err := cr.BaseRunner.decide(cr, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (cr CommitteeRunner) GetSigner() types.BeaconSigner {
	return cr.signer
}

func (cr CommitteeRunner) GetOperatorSigner() types.OperatorSigner {
	return cr.operatorSigner
}

func constructAttestationData(vote *types.BeaconVote, duty *types.BeaconDuty) *phase0.AttestationData {
	return &phase0.AttestationData{
		Slot:            duty.Slot,
		Index:           duty.CommitteeIndex,
		BeaconBlockRoot: vote.BlockRoot,
		Source:          vote.Source,
		Target:          vote.Target,
	}
}
func ConstructSyncCommittee(vote *types.BeaconVote, duty *types.BeaconDuty) *altair.SyncCommitteeMessage {
	return &altair.SyncCommitteeMessage{
		Slot:            duty.Slot,
		BeaconBlockRoot: vote.BlockRoot,
		ValidatorIndex:  duty.ValidatorIndex,
	}
}
