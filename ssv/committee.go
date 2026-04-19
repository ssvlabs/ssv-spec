package ssv

import (
	"fmt"
	"sync"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

type CreateRunnerFn func(shareMap map[phase0.ValidatorIndex]*types.Share) Runner

type Committee struct {
	CommitteeRunners                  map[phase0.Slot]Runner
	AggregatorCommitteeRunners        map[phase0.Slot]Runner
	CommitteeMember                   types.CommitteeMember
	CreateCommitteeRunnerFn           CreateRunnerFn
	CreateAggregatorCommitteeRunnerFn CreateRunnerFn
	Share                             map[phase0.ValidatorIndex]*types.Share

	validateOnce sync.Once
	validateErr  error
}

// NewCommittee creates a new cluster
func NewCommittee(
	committeeMember types.CommitteeMember,
	share map[phase0.ValidatorIndex]*types.Share,
	createRunnerFn CreateRunnerFn,
	createAggregatorCommitteeRunnerFn CreateRunnerFn,
) *Committee {
	c := &Committee{
		CommitteeRunners:                  make(map[phase0.Slot]Runner),
		AggregatorCommitteeRunners:        make(map[phase0.Slot]Runner),
		CommitteeMember:                   committeeMember,
		CreateCommitteeRunnerFn:           createRunnerFn,
		CreateAggregatorCommitteeRunnerFn: createAggregatorCommitteeRunnerFn,
		Share:                             share,
	}
	return c
}

func (c *Committee) validateInvariants() error {
	c.validateOnce.Do(func() {
		if err := (&c.CommitteeMember).Validate(); err != nil {
			c.validateErr = errors.Wrap(err, "invalid committee member")
			return
		}
		if err := validateShareMap(c.Share); err != nil {
			c.validateErr = errors.Wrap(err, "invalid share map")
			return
		}
	})
	return c.validateErr
}

// StartDuty starts a new duty for the given slot
func (c *Committee) StartDuty(duty types.Duty) error {
	if err := c.validateInvariants(); err != nil {
		return err
	}

	// Get objects according to duty type
	var slot phase0.Slot
	var runnerMap *map[phase0.Slot]Runner
	var createFn *CreateRunnerFn
	var validatorDuties []*types.ValidatorDuty
	switch d := duty.(type) {
	case *types.CommitteeDuty:
		slot = phase0.Slot(d.Slot)
		runnerMap = &c.CommitteeRunners
		createFn = &c.CreateCommitteeRunnerFn
		validatorDuties = d.ValidatorDuties
	case *types.AggregatorCommitteeDuty:
		if d == nil {
			return types.NewError(types.InvalidAggregatorCommitteeDutyErrorCode, "nil aggregator committee duty")
		}
		slot = phase0.Slot(d.Slot)
		runnerMap = &c.AggregatorCommitteeRunners
		createFn = &c.CreateAggregatorCommitteeRunnerFn
		validatorDuties = d.ValidatorDuties
	default:
		return errors.Errorf("unsupported duty type: %T", duty)
	}

	if _, exists := (*runnerMap)[slot]; exists {
		return fmt.Errorf("Runner for slot %d already exists", slot)
	}

	if len(validatorDuties) == 0 {
		return types.NewError(types.NoBeaconDutiesErrorCode, "no beacon duties")
	}

	// Filter duty and create share map according validators that belong to c.Share
	dutyShares := make(map[phase0.ValidatorIndex]*types.Share)
	filteredValidatorDuties := make([]*types.ValidatorDuty, 0)

	for _, bduty := range validatorDuties {
		if bduty == nil {
			// Preserve the previous best-effort behavior for malformed foreign entries:
			// only duties that can be matched to this committee's shares are enforced.
			continue
		}
		if _, exists := c.Share[bduty.ValidatorIndex]; !exists {
			continue
		}
		if err := bduty.Validate(); err != nil {
			return errors.Wrap(err, "invalid validator duty")
		}
		dutyShares[bduty.ValidatorIndex] = c.Share[bduty.ValidatorIndex]
		filteredValidatorDuties = append(filteredValidatorDuties, bduty)
	}

	if len(dutyShares) == 0 {
		return types.NewError(types.NoValidatorSharesErrorCode, "no shares for duty's validators")
	}

	var filteredDuty types.Duty
	switch duty.(type) {
	case *types.CommitteeDuty:
		committeeDuty := &types.CommitteeDuty{
			Slot:            slot,
			ValidatorDuties: filteredValidatorDuties,
		}
		if err := committeeDuty.Validate(); err != nil {
			return errors.Wrap(err, "invalid committee duty")
		}
		filteredDuty = committeeDuty
	case *types.AggregatorCommitteeDuty:
		aggregatorDuty := &types.AggregatorCommitteeDuty{
			Slot:            slot,
			ValidatorDuties: filteredValidatorDuties,
		}
		validatorIndex := make(map[phase0.ValidatorIndex]struct{}, len(dutyShares))
		for validatorIdx := range dutyShares {
			validatorIndex[validatorIdx] = struct{}{}
		}
		if err := aggregatorDuty.Validate(validatorIndex); err != nil {
			return errors.Wrap(err, "invalid aggregator committee duty")
		}
		filteredDuty = aggregatorDuty
	default:
		return errors.Errorf("unsupported duty type: %T", duty)
	}

	(*runnerMap)[slot] = (*createFn)(dutyShares)
	return (*runnerMap)[slot].StartNewDuty(filteredDuty, c.CommitteeMember.GetQuorum())
}

// ProcessMessage processes Network Message of all types
func (c *Committee) ProcessMessage(signedSSVMessage *types.SignedSSVMessage) error {
	if err := c.validateInvariants(); err != nil {
		return err
	}

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

	// Get runner map according to message role
	var runnerMap *map[phase0.Slot]Runner
	role := msg.MsgID.GetRoleType()
	switch role {
	case types.RoleCommittee:
		runnerMap = &c.CommitteeRunners
	case types.RoleAggregatorCommittee:
		runnerMap = &c.AggregatorCommitteeRunners
	default:
		return types.NewError(types.CommitteeWrongRoleErrorCode, "msg role is invalid")
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

		runner, exists := (*runnerMap)[phase0.Slot(qbftMsg.Height)]
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

		runner, exists := (*runnerMap)[pSigMessages.Slot]
		if !exists {
			return types.NewError(types.NoRunnerForSlotErrorCode, "no runner found for message's slot")
		}

		switch pSigMessages.Type {
		case types.PostConsensusPartialSig:
			return runner.ProcessPostConsensus(pSigMessages)
		case types.AggregatorCommitteePartialSig:
			if role != types.RoleAggregatorCommittee {
				return errors.Errorf("invalid aggregator partial sig msg for commmittee role")
			}
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

	role := msg.GetID().GetRoleType()
	if role != types.RoleCommittee && role != types.RoleAggregatorCommittee {
		return types.NewError(types.CommitteeWrongRoleErrorCode, "msg role is invalid")
	}

	if len(msg.GetData()) == 0 {
		return fmt.Errorf("msg data is invalid")
	}

	return nil
}
