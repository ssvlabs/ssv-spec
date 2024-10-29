package validation

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/ssvlabs/ssv-spec/types"
)

// Validates a SignedSSVMessage on duty logic rules
func (mv *MessageValidator) ValidateMessageDutyLogic(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage, receivedAt time.Time) error {

	msgID := signedSSVMessage.SSVMessage.MsgID
	role := signedSSVMessage.SSVMessage.MsgID.GetRoleType()

	slot, err := getMessageSlot(signedSSVMessage)
	if err != nil {
		return err
	}

	// Rule: Height must not be "old". I.e., signer must not have already advanced to a later slot.
	if role != types.RoleCommittee { // Rule only for validator runners
		for _, signer := range signedSSVMessage.OperatorIDs {
			if MessageFromOldSlot(mv.GetPeerState(peerID), signer, msgID, slot) {
				return ErrSlotAlreadyAdvanced
			}
		}
	}

	// Rule: For proposal and sync committee aggregation duties, we check if the validator is assigned to it
	if err := mv.ValidBeaconDuty(msgID, slot); err != nil {
		return err
	}

	// Rule: current slot(height) must be between duty's starting slot and:
	// - duty's starting slot + 34 (committee and aggregation)
	// - duty's starting slot + 3 (other types)
	if err := mv.ValidDutySlot(slot, role, receivedAt); err != nil {
		// Err should be ErrEarlySlotMessage or ErrLateSlotMessage
		return err
	}

	// Rule: valid number of duties per epoch:
	// - 2 for aggregation, voluntary exit and validator registration
	// - 2*V for Committee duty (where V is the number of validators in the cluster) (if no validator is doing sync committee in this epoch)
	// - else, accept
	if err := mv.ValidNumberOfDutiesPerEpoch(peerID, msgID, slot); err != nil {
		return err
	}

	return nil
}

// Auxiliary functions

// Check if, in the peer's view, the signer has already advanced to a later slot
func MessageFromOldSlot(peerState *PeerState, signer types.OperatorID, msgID types.MessageID, slot phase0.Slot) bool {
	return (peerState.GetHighestSlotForSigner(msgID, signer) > slot)
}

// Check if the current time is valid according to the duty's slot.
// I.e. if it's not too early or to late considering the duty's slot
func (mv *MessageValidator) ValidDutySlot(slot phase0.Slot, role types.RunnerRole, receivedAt time.Time) error {
	if earliness := mv.messageEarliness(slot, receivedAt); earliness > ClockErrorTolerance {
		return ErrEarlySlotMessage
	}

	if lateness := mv.messageLateness(slot, role, receivedAt); lateness > ClockErrorTolerance {
		return ErrLateSlotMessage
	}

	return nil
}

// Returns how early message is or 0 if it's not
func (mv *MessageValidator) messageEarliness(slot phase0.Slot, receivedAt time.Time) time.Duration {
	return time.Unix(mv.Beacon.EstimatedTimeAtSlot(slot), 0).Sub(receivedAt)
}

// Returns how late message is (compared to the latest possible message) or 0 if it's not
func (mv *MessageValidator) messageLateness(slot phase0.Slot, role types.RunnerRole, receivedAt time.Time) time.Duration {
	var ttl phase0.Slot
	switch role {
	case types.RoleProposer, types.RoleSyncCommitteeContribution:
		ttl = 1 + LateSlotAllowance
	case types.RoleCommittee, types.RoleAggregator:
		ttl = phase0.Slot(mv.Beacon.SlotsPerEpoch()) + LateSlotAllowance
	case types.RoleValidatorRegistration, types.RoleVoluntaryExit:
		return 0
	}

	deadline := time.Unix(mv.Beacon.EstimatedTimeAtSlot(slot+ttl), 0).
		Add(LateMessageMargin)

	return receivedAt.Sub(deadline)
}

// Get the epoch duty limit according to the role
func (mv *MessageValidator) DutyLimitForRole(role types.RunnerRole, committeeInfo *CommitteeInfo, slot phase0.Slot) int {
	switch role {
	case types.RoleAggregator, types.RoleValidatorRegistration, types.RoleVoluntaryExit:
		return 2
	case types.RoleCommittee:
		// If at least one validator is in sync committee, the limit is 32 (all slots in the epoch)
		for _, validator := range committeeInfo.Validators {
			if mv.DutyFetcher.HasSyncCommitteeDuty(validator, slot) {
				return 32
			}
		}
		// min(2*V, 32), where V is the number of validators in the committee
		limit := 2 * len(committeeInfo.Validators)
		if limit > 32 {
			limit = 32
		}
		return limit
	default:
		return 0
	}
}

// Check if the number of duties per epoch is valid for a msg (considering its identifier (validator/committee and role) and slot)
// according to the duty counter state
func (mv *MessageValidator) ValidNumberOfDutiesPerEpoch(peerID peer.ID, msgID types.MessageID, slot phase0.Slot) error {

	role := msgID.GetRoleType()

	// For proposer and sync committee contribution, we don't count the number of duties
	// but rather check with the Beacon node if it exists
	if role == types.RoleProposer || role == types.RoleSyncCommitteeContribution {
		return nil
	}

	dutyCounter := mv.GetDutyCounter(peerID)
	epoch := mv.Beacon.EstimatedEpochAtSlot(slot)
	hasDuty := dutyCounter.HasDuty(msgID, epoch, slot)
	numDuties := dutyCounter.CountDutiesForEpoch(msgID, epoch)
	dutyLimit := mv.DutyLimitForRole(role, mv.Network.GetCommitteeInfo(msgID), slot)

	if hasDuty {
		// Duty already validated
		// Just do a sanity check
		if numDuties > dutyLimit {
			return ErrTooManyDutiesPerEpoch
		}
	} else {
		// New duty, so we compare the limit against the incremented value
		newNumDuties := numDuties + 1
		if newNumDuties > dutyLimit {
			return ErrTooManyDutiesPerEpoch
		}
	}

	return nil
}

// Check if the Beacon duty is valid according to the duty fetcher
func (mv *MessageValidator) ValidBeaconDuty(msgID types.MessageID, slot phase0.Slot) error {

	committeeInfo := mv.Network.GetCommitteeInfo(msgID)
	role := msgID.GetRoleType()

	// Rule: For a proposal duty message, we check if the validator is assigned to it
	if role == types.RoleProposer {
		// Non-committee roles always have one validator index.
		validatorIndex := committeeInfo.Validators[0]
		if !mv.DutyFetcher.HasProposerDuty(validatorIndex, slot) {
			return ErrNoDuty
		}
	}

	// Rule: For a sync committee aggregation duty message, we check if the validator is assigned to it
	if role == types.RoleSyncCommitteeContribution {
		// Non-committee roles always have one validator index.
		validatorIndex := committeeInfo.Validators[0]
		if !mv.DutyFetcher.HasSyncCommitteeContributionDuty(validatorIndex, slot) {
			return ErrNoDuty
		}
	}

	return nil
}
