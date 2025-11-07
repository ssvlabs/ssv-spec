package ssv

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CreateRunnerFn func(shareMap map[phase0.ValidatorIndex]*types.Share) Runner

type Committee struct {
	Runners         map[phase0.Slot]Runner
	CommitteeMember types.CommitteeMember
	CreateRunnerFn  CreateRunnerFn
	Share           map[phase0.ValidatorIndex]*types.Share
}

// NewCommittee creates a new cluster
func NewCommittee(
	committeeMember types.CommitteeMember,
	share map[phase0.ValidatorIndex]*types.Share,
	createRunnerFn CreateRunnerFn,
) *Committee {
	c := &Committee{
		Runners:         make(map[phase0.Slot]Runner),
		CommitteeMember: committeeMember,
		CreateRunnerFn:  createRunnerFn,
		Share:           share,
	}
	return c
}

// StartDuty starts a new duty for the given slot
func (c *Committee) StartDuty(duty types.Duty) error {
	slot := duty.DutySlot()
	if _, exists := c.Runners[slot]; exists {
		return fmt.Errorf("Runner for slot %d already exists", slot)
	}

	// Handle different duty types
	switch d := duty.(type) {
	case *types.CommitteeDuty:
		if len(d.ValidatorDuties) == 0 {
			return types.NewError(types.NoBeaconDutiesErrorCode, "no beacon duties")
		}

		// Filter duty and create share map according validators that belong to c.Share
		dutyShares := make(map[phase0.ValidatorIndex]*types.Share)
		filteredDuty := &types.CommitteeDuty{
			Slot: d.Slot,
		}

		for _, bduty := range d.ValidatorDuties {
			if _, exists := c.Share[bduty.ValidatorIndex]; !exists {
				continue
			}
			dutyShares[bduty.ValidatorIndex] = c.Share[bduty.ValidatorIndex]
			filteredDuty.ValidatorDuties = append(filteredDuty.ValidatorDuties, bduty)
		}

		if len(dutyShares) == 0 {
			return types.NewError(types.NoValidatorSharesErrorCode, "no shares for duty's validators")
		}

		c.Runners[slot] = c.CreateRunnerFn(dutyShares)
		return c.Runners[slot].StartNewDuty(filteredDuty, c.CommitteeMember.GetQuorum())

	case *types.AggregatorCommitteeDuty:
		if len(d.ValidatorDuties) == 0 {
			return types.NewError(types.NoBeaconDutiesErrorCode, "no beacon duties")
		}

		// Filter duty and create share map according to validators that belong to c.Share
		dutyShares := make(map[phase0.ValidatorIndex]*types.Share)
		filteredDuty := &types.AggregatorCommitteeDuty{
			Slot: d.Slot,
		}

		for _, bduty := range d.ValidatorDuties {
			if _, exists := c.Share[bduty.ValidatorIndex]; !exists {
				continue
			}
			dutyShares[bduty.ValidatorIndex] = c.Share[bduty.ValidatorIndex]
			filteredDuty.ValidatorDuties = append(filteredDuty.ValidatorDuties, bduty)
		}

		if len(dutyShares) == 0 {
			return types.NewError(types.NoValidatorSharesErrorCode, "no shares for duty's validators")
		}

		c.Runners[slot] = c.CreateRunnerFn(dutyShares)
		return c.Runners[slot].StartNewDuty(filteredDuty, c.CommitteeMember.GetQuorum())

	default:
		return errors.Errorf("unsupported duty type: %T", duty)
	}
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	// Validate message
	if err := signedSSVMessage.Validate(); err != nil {
		return errors.Wrap(err, "invalid SignedSSVMessage")
	}

	// Verify SignedSSVMessage's signature
	if err := types.Verify(signedSSVMessage, c.CommitteeMember.Committee); err != nil {
		return types.WrapError(types.SSVMessageHasInvalidSignatureErrorCode, fmt.Errorf("SignedSSVMessage has an invalid signature: %w", err))
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

		runner, exists := c.Runners[phase0.Slot(qbftMsg.Height)]
		if !exists {
			return types.NewError(types.NoRunnerForSlotErrorCode, "no runner found for message's slot")
		}
		return runner.ProcessConsensus(signedSSVMessage)
	case types.SSVPartialSignatureMsgType:
		pSigMessages := &types.PartialSignatureMessages{}
		if err := pSigMessages.Decode(msg.GetData()); err != nil {
			return errors.Wrap(err, "could not get post consensus Message from network Message")
		}

		// Validate
		if len(signedSSVMessage.OperatorIDs) != 1 {
			return fmt.Errorf("PartialSignatureMessage has more than 1 signer")
		}

		if err := pSigMessages.ValidateForSigner(signedSSVMessage.OperatorIDs[0]); err != nil {
			return errors.Wrap(err, "invalid PartialSignatureMessages")
		}

		runner, exists := c.Runners[pSigMessages.Slot]
		if !exists {
			return types.NewError(types.NoRunnerForSlotErrorCode, "no runner found for message's slot")
		}

		switch pSigMessages.Type {
		case types.PostConsensusPartialSig:
			return runner.ProcessPostConsensus(pSigMessages)
		case types.SelectionProofPartialSig, types.ContributionProofs, types.AggregatorCommitteePartialSig:
			// Pre-consensus messages
			return runner.ProcessPreConsensus(pSigMessages)
		default:
			return errors.Errorf("unknown partial signature message type: %v", pSigMessages.Type)
		}
	default:
		return fmt.Errorf("unknown msg")
	}
}

func (c *Committee) validateMessage(msg *types.SSVMessage) error {
	if !(c.CommitteeMember.CommitteeID.MessageIDBelongs(msg.GetID())) {
		return types.NewError(types.MessageIDCommitteeIDMismatchErrorCode, "msg ID doesn't match committee ID")
	}

	if len(msg.GetData()) == 0 {
		return fmt.Errorf("msg data is invalid")
	}

	return nil
}
