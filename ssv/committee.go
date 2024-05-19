package ssv

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CreateRunnerFn func(shareMap map[spec.ValidatorIndex]*types.Share) *CommitteeRunner

type Committee struct {
	Runners            map[spec.Slot]*CommitteeRunner
	Operator           types.Operator
	SignatureVerifier  types.SignatureVerifier
	CreateRunnerFn     CreateRunnerFn
	Share              map[spec.ValidatorIndex]*types.Share
	HighestDutySlotMap map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot
}

// NewCommittee creates a new cluster
func NewCommittee(
	operator types.Operator,
	verifier types.SignatureVerifier,
	share map[spec.ValidatorIndex]*types.Share,
	createRunnerFn CreateRunnerFn,
) *Committee {
	c := &Committee{
		Runners:            make(map[spec.Slot]*CommitteeRunner),
		Operator:           operator,
		SignatureVerifier:  verifier,
		CreateRunnerFn:     createRunnerFn,
		Share:              share,
		HighestDutySlotMap: make(map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot),
	}
	c.HighestDutySlotMap[types.BNRoleAttester] = make(map[spec.ValidatorIndex]spec.Slot)
	c.HighestDutySlotMap[types.BNRoleSyncCommittee] = make(map[spec.ValidatorIndex]spec.Slot)
	return c
}

// StartDuty starts a new duty for the given slot
func (c *Committee) StartDuty(duty *types.CommitteeDuty) error {
	if _, exists := c.Runners[duty.Slot]; exists {
		return errors.New(fmt.Sprintf("CommitteeRunner for slot %d already exists", duty.Slot))
	}
	c.Runners[duty.Slot] = c.CreateRunnerFn(c.Share)
	var validatorToStopMap map[spec.Slot][]spec.ValidatorIndex
	// Filter old duties based on highest duty slot
	duty, validatorToStopMap, c.HighestDutySlotMap = FilterCommitteeDuty(duty, c.HighestDutySlotMap)
	// Stop validators with old duties
	c.stopDuties(validatorToStopMap)
	c.updateDutySlotMap(duty)
	// TODO: check if there are beacon duties remaining
	return c.Runners[duty.Slot].StartNewDuty(duty)
}

func (c *Committee) stopDuties(validatorToStopMap map[spec.Slot][]spec.ValidatorIndex) {
	for slot, validators := range validatorToStopMap {
		for _, validator := range validators {
			runner, exists := c.Runners[slot]
			if exists {
				runner.StopDuty(validator)
			}
		}
	}
}

// FilterCommitteeDuty filters the committee duties by the slots given per validator.
// It returns the filtered duties, the validators to stop and updated slot map.
func FilterCommitteeDuty(duty *types.CommitteeDuty, dutySlotMap map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot) (
	*types.CommitteeDuty,
	map[spec.Slot][]spec.ValidatorIndex,
	map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot,
) {
	validatorsToStop := make(map[spec.Slot][]spec.ValidatorIndex)

	for i, beaconDuty := range duty.BeaconDuties {

		if _, exists := dutySlotMap[beaconDuty.Type]; !exists {
			dutySlotMap[beaconDuty.Type] = make(map[spec.ValidatorIndex]spec.Slot)
		}

		slot, exists := dutySlotMap[beaconDuty.Type][beaconDuty.ValidatorIndex]
		if exists {
			if slot < beaconDuty.Slot {
				if _, exists := validatorsToStop[slot]; !exists {
					validatorsToStop[slot] = make([]spec.ValidatorIndex, 0)
				}
				validatorsToStop[slot] = append(validatorsToStop[slot], beaconDuty.ValidatorIndex)
				dutySlotMap[beaconDuty.Type][beaconDuty.ValidatorIndex] = beaconDuty.Slot
			} else { // else don't run duty with old slot
				duty.BeaconDuties[i] = nil
			}
		}
	}
	return duty, validatorsToStop, dutySlotMap
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := c.SignatureVerifier.Verify(signedSSVMessage, c.Operator.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	msg := signedSSVMessage.SSVMessage
	if err := c.validateMessage(msg); err != nil {
		return errors.Wrap(err, "Message invalid")
	}

	switch msg.GetType() {
	case types.SSVConsensusMsgType:
		qbftMsg := &qbft.Message{}
		if err := qbftMsg.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}

		if err := qbftMsg.Validate(); err != nil {
			return errors.Wrap(err, "invalid qbft Message")
		}

		runner, exists := c.Runners[spec.Slot(qbftMsg.Height)]
		if !exists {
			return errors.New("no runner found for message's slot")
		}
		return runner.ProcessConsensus(signedSSVMessage)
	case types.SSVPartialSignatureMsgType:
		pSigMessages := &types.PartialSignatureMessages{}
		if err := pSigMessages.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		// Validate
		if len(signedSSVMessage.OperatorIDs) != 1 {
			return errors.New("PartialSignatureMessage has more than 1 signer")
		}

		if err := pSigMessages.ValidateForSigner(signedSSVMessage.OperatorIDs[0]); err != nil {
			return errors.Wrap(err, "invalid PartialSignatureMessages")
		}

		if pSigMessages.Type == types.PostConsensusPartialSig {
			runner, exists := c.Runners[pSigMessages.Slot]
			if !exists {
				return errors.New("no runner found for message's slot")
			}
			return runner.ProcessPostConsensus(pSigMessages)
		}
	default:
		return errors.New("unknown msg")
	}
	return nil

}

// updateAttestingSlotMap updates the highest attesting slot map from beacon duties
func (c *Committee) updateDutySlotMap(duty *types.CommitteeDuty) {
	for _, beaconDuty := range duty.BeaconDuties {

		if _, exists := c.HighestDutySlotMap[beaconDuty.Type]; !exists {
			c.HighestDutySlotMap[beaconDuty.Type] = make(map[spec.ValidatorIndex]spec.Slot)
		}

		if _, ok := c.HighestDutySlotMap[beaconDuty.Type][beaconDuty.ValidatorIndex]; !ok {
			c.HighestDutySlotMap[beaconDuty.Type][beaconDuty.ValidatorIndex] = beaconDuty.Slot
		}
		if c.HighestDutySlotMap[beaconDuty.Type][beaconDuty.ValidatorIndex] < beaconDuty.Slot {
			c.HighestDutySlotMap[beaconDuty.Type][beaconDuty.ValidatorIndex] = beaconDuty.Slot
		}
	}
}

func (c *Committee) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Committee) Decode(data []byte) error {
	return json.Unmarshal(data, &c)
}

// GetRoot returns the state's deterministic root
func (c *Committee) GetRoot() ([32]byte, error) {
	marshaledRoot, err := c.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode state")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}

func (c *Committee) MarshalJSON() ([]byte, error) {

	type CommitteeAlias struct {
		Runners            map[spec.Slot]*CommitteeRunner
		Operator           types.Operator
		Share              map[spec.ValidatorIndex]*types.Share
		HighestDutySlotMap map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot
	}

	// Create object and marshal
	alias := &CommitteeAlias{
		Runners:            c.Runners,
		Operator:           c.Operator,
		Share:              c.Share,
		HighestDutySlotMap: c.HighestDutySlotMap,
	}

	byts, err := json.Marshal(alias)

	return byts, err
}

func (c *Committee) UnmarshalJSON(data []byte) error {

	type CommitteeAlias struct {
		Runners            map[spec.Slot]*CommitteeRunner
		Operator           types.Operator
		Share              map[spec.ValidatorIndex]*types.Share
		HighestDutySlotMap map[types.BeaconRole]map[spec.ValidatorIndex]spec.Slot
	}

	// Unmarshal the JSON data into the auxiliary struct
	aux := &CommitteeAlias{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Assign fields
	c.Runners = aux.Runners
	c.Operator = aux.Operator
	c.Share = aux.Share
	c.HighestDutySlotMap = aux.HighestDutySlotMap

	return nil
}

func (c *Committee) validateMessage(msg *types.SSVMessage) error {
	if !(c.Operator.CommitteeID.MessageIDBelongs(msg.GetID())) {
		return errors.New("msg ID doesn't match committee ID")
	}

	if len(msg.GetData()) == 0 {
		return errors.New("msg data is invalid")
	}

	return nil
}
