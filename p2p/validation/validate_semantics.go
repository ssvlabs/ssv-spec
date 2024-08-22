package validation

import (
	"github.com/ssvlabs/ssv-spec/types"
)

// Validates a SignedSSVMessage on general semantic rules
func (mv *MessageValidator) validateGeneralSemantics(msg *types.SignedSSVMessage, topic string) error {

	role := msg.SSVMessage.MsgID.GetRoleType()

	// Rule: If domain is different then self domain
	domain := msg.SSVMessage.MsgID.GetDomain()
	if !mv.Network.ValidDomain(domain) {
		return ErrWrongDomain
	}

	// Rule: If role is invalid
	if !ValidRole(role) {
		return ErrInvalidRole
	}

	senderID := msg.SSVMessage.MsgID.GetDutyExecutorID()
	if role != types.RoleCommittee {
		validatorPK := types.ValidatorPK(senderID)

		// Rule: Validator does not exist
		if !mv.Network.ExistingValidator(validatorPK) {
			return ErrUnknownValidator
		}

		// Rule: If validator is not active
		if !mv.Network.ActiveValidator(validatorPK) {
			return ErrValidatorNotAttesting
		}

		// Rule: If validator is liquidated
		if mv.Network.ValidatorLiquidated(validatorPK) {
			return ErrValidatorLiquidated
		}
	} else {
		// Rule: Cluster does not exist
		if !mv.Network.ExistingCommitteeID(types.CommitteeID(senderID)) {
			return ErrNonExistentCommitteeID
		}
	}

	committeeInfo := mv.Network.GetCommitteeInfo(msg.SSVMessage.MsgID)
	if committeeInfo == nil { // Sanity check. Shouldn't be nil due to above checks
		return ErrNonExistentCommitteeID
	}

	// Rule: Signers must belong to validator committee or CommitteeID
	for _, signer := range msg.OperatorIDs {
		if !SignerBelongsToCommittee(signer, committeeInfo.OperatorIDs()) {
			return ErrSignerNotInCommittee
		}
	}

	// Rule: Check if message was sent in the correct topic
	if !mv.Network.CorrectTopic(committeeInfo.OperatorIDs(), topic) {
		return ErrIncorrectTopic
	}

	if msg.SSVMessage.MsgType == types.DKGMsgType {
		// Rule: DKG message not allowed
		return ErrDKGMessage
	} else if msg.SSVMessage.MsgType != types.SSVConsensusMsgType &&
		msg.SSVMessage.MsgType != types.SSVPartialSignatureMsgType {
		return ErrUnknownSSVMessageType
	}

	return nil
}

// Auxiliary functions

// Check if the signer belongs to the committee
func SignerBelongsToCommittee(signer types.OperatorID, committee []types.OperatorID) bool {
	for _, operator := range committee {
		if signer == operator {
			return true
		}
	}
	return false
}

// Check if the runner role is valid
func ValidRole(role types.RunnerRole) bool {
	switch role {
	case types.RoleCommittee,
		types.RoleProposer,
		types.RoleAggregator,
		types.RoleSyncCommitteeContribution,
		types.RoleVoluntaryExit,
		types.RoleValidatorRegistration:
		return true
	default:
		return false
	}
}
