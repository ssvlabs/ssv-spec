package ssv

import (
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CreateRunnerFn func(shareMap map[spec.ValidatorIndex]*types.Share) *CommitteeRunner

type Committee struct {
	CommitteeRunners map[spec.Slot]*CommitteeRunner
	CommitteeMember  types.CommitteeMember
	CreateRunnerFn   CreateRunnerFn
	Validators       map[spec.ValidatorIndex]*Validator
}

// NewCommittee creates a new cluster
func NewCommittee(
	committeeMember types.CommitteeMember,
	validators map[spec.ValidatorIndex]*Validator,
	createRunnerFn CreateRunnerFn,
) *Committee {
	c := &Committee{
		CommitteeRunners: make(map[spec.Slot]*CommitteeRunner),
		CommitteeMember:  committeeMember,
		CreateRunnerFn:   createRunnerFn,
		Validators:       validators,
	}
	return c
}

// StartDuty starts a duty redirecting it according to its type (committee / validator duty)
func (c *Committee) StartDuty(duty types.Duty) error {
	switch duty := duty.(type) {
	case *types.CommitteeDuty:
		return c.StartCommitteeDuty(duty)
	case *types.ValidatorDuty:
		return c.StartValidatorDuty(duty)
	default:
		return errors.New("unknown duty object")
	}
}

// StartCommitteeDuty starts a new validator duty
func (c *Committee) StartValidatorDuty(duty *types.ValidatorDuty) error {
	if _, exists := c.Validators[duty.ValidatorIndex]; !exists {
		return errors.New("unknown duty's validator")
	}
	return c.Validators[duty.ValidatorIndex].StartDuty(duty)
}

// StartCommitteeDuty starts a new committee duty for the given slot
func (c *Committee) StartCommitteeDuty(duty *types.CommitteeDuty) error {
	if len(duty.ValidatorDuties) == 0 {
		return errors.New("no beacon duties")
	}
	if _, exists := c.CommitteeRunners[duty.Slot]; exists {
		return errors.New(fmt.Sprintf("CommitteeRunner for slot %d already exists", duty.Slot))
	}

	// Filter duty and create share map according validators that belong to c.Share
	dutyShares := make(map[spec.ValidatorIndex]*types.Share)
	filteredDuty := &types.CommitteeDuty{
		Slot: duty.Slot,
	}

	for _, bduty := range duty.ValidatorDuties {
		if _, exists := c.Validators[bduty.ValidatorIndex]; !exists {
			continue
		}
		dutyShares[bduty.ValidatorIndex] = c.Validators[bduty.ValidatorIndex].Share
		filteredDuty.ValidatorDuties = append(filteredDuty.ValidatorDuties, bduty)
	}

	if len(dutyShares) == 0 {
		return errors.New("no shares for duty's validators")
	}

	c.CommitteeRunners[filteredDuty.Slot] = c.CreateRunnerFn(dutyShares)

	return c.CommitteeRunners[filteredDuty.Slot].StartNewDuty(filteredDuty, c.CommitteeMember.GetQuorum())
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(msg *types.SignedSSVMessage) error {

	// Validate message
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := types.Verify(msg, c.CommitteeMember.Committee); err != nil {
		return errors.Wrap(err, "SignedSSVMessage has an invalid signature")
	}

	// Process message
	if c.isMessageForCommittee(msg) {
		return c.ProcessMessageForCommitteeDuty(msg)
	} else {
		for _, validator := range c.Validators {
			if validator.isMessageForValidator(msg) {
				return validator.ProcessMessage(msg)
			}
		}
	}
	return errors.New("message doesn't belong to committee or one of its validators")
}

// ProcessMessage processes a message of all types for a committee duty
func (c *Committee) ProcessMessageForCommitteeDuty(msg *types.SignedSSVMessage) error {

	switch msg.SSVMessage.MsgType {
	case types.SSVConsensusMsgType:
		// Get inner message to get slot
		qbftMsg := &qbft.Message{}
		if err := qbftMsg.Decode(msg.SSVMessage.Data); err != nil {
			return errors.Wrap(err, "could not get consensus Message from network Message")
		}
		// Get runner
		runner, exists := c.CommitteeRunners[spec.Slot(qbftMsg.Height)]
		if !exists {
			return errors.New("no runner found for message's slot")
		}
		// Process message
		return RunnerProcessMessage(runner, msg)
	case types.SSVPartialSignatureMsgType:
		// Get inner message to get slot
		pSigMessages := &types.PartialSignatureMessages{}
		if err := pSigMessages.Decode(msg.SSVMessage.Data); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}
		// Get runner
		runner, exists := c.CommitteeRunners[pSigMessages.Slot]
		if !exists {
			return errors.New("no runner found for message's slot")
		}
		// Process message
		return RunnerProcessMessage(runner, msg)
	default:
		return errors.New("unknown msg")
	}
}

// Returns true if message is intended to the committee according to its MessageID
func (c *Committee) isMessageForCommittee(msg *types.SignedSSVMessage) bool {
	return c.CommitteeMember.CommitteeID.MessageIDBelongs(msg.SSVMessage.MsgID)
}
