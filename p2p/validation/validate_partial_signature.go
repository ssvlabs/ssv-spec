package validation

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/ssvlabs/ssv-spec/types"
)

func (mv *MessageValidator) ValidatePartialSignatureMessageSemantics(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	// Decode to partial signature message
	partialSignatureMessages := &types.PartialSignatureMessages{}
	err := partialSignatureMessages.Decode(signedSSVMessage.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	msgID := signedSSVMessage.SSVMessage.MsgID

	// Rule: Partial Signature message must have 1 signer
	signers := signedSSVMessage.OperatorIDs
	if len(signers) != 1 {
		return ErrPartialSigMultipleSigners
	}
	signer := signers[0]

	// Rule: Partial signature message must not have full data
	if len(signedSSVMessage.FullData) > 0 {
		return ErrFullDataNotInConsensusMessage
	}

	// Rule: Valid signature type
	if !ValidPartialSignatureType(partialSignatureMessages.Type) {
		return ErrInvalidPartialSignatureType
	}

	// Rule: Partial signature type must match expected type:
	// - PostConsensusPartialSig, for Committee duty
	// - RandaoPartialSig or PostConsensusPartialSig for Proposer
	// - SelectionProofPartialSig or PostConsensusPartialSig for Aggregator
	// - SelectionProofPartialSig or PostConsensusPartialSig for Sync committee contribution
	// - ValidatorRegistrationPartialSig for Validator Registration
	// - VoluntaryExitPartialSig for Voluntary Exit
	if !ExpectedPartialSignatureTypeForRole(partialSignatureMessages.Type, msgID.GetRoleType()) {
		return ErrPartialSignatureTypeRoleMismatch
	}

	// Rule: Partial signature message must have at least one signature
	if len(partialSignatureMessages.Messages) == 0 {
		return ErrNoPartialSignatureMessages
	}

	for _, psigMsg := range partialSignatureMessages.Messages {
		// Rule: Partial signature must have expected length
		if len(psigMsg.PartialSignature) != PartialSignatureSize {
			return ErrWrongBLSSignatureSize
		}
		// Rule: Partial signature signer must be consistent
		if psigMsg.Signer != signer {
			return ErrInconsistentSigners
		}
		// Rule: (only for Validator duties) Validator index must match with validatorPK
		// For Committee duties, we don't assume that operators are synced on the validators set
		// So, we can't make this assertion
		if msgID.GetRoleType() != types.RoleCommittee {
			if !ValidatorIndexBelongsToCommittee(psigMsg.ValidatorIndex, mv.Network.GetCommitteeInfo(msgID)) {
				return ErrValidatorIndexMismatch
			}
		}
	}

	return nil
}

func (mv *MessageValidator) ValidatePartialSignatureMessageDutyLogic(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage, receivedAt time.Time) error {

	// Decode to partial signature message
	partialSignatureMessages := &types.PartialSignatureMessages{}
	err := partialSignatureMessages.Decode(signedSSVMessage.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	signer := signedSSVMessage.OperatorIDs[0]
	msgID := signedSSVMessage.SSVMessage.MsgID
	role := msgID.GetRoleType()

	peerState := mv.GetPeerState(peerID)
	dutyState := peerState.GetDutyState(msgID, partialSignatureMessages.Slot)
	signerState := dutyState.GetSignerState(signer)

	// Rule: peer must send only:
	// - 1 PostConsensusPartialSig, for Committee duty
	// - 1 RandaoPartialSig and 1 PostConsensusPartialSig for Proposer
	// - 1 SelectionProofPartialSig and 1 PostConsensusPartialSig for Aggregator
	// - 1 SelectionProofPartialSig and 1 PostConsensusPartialSig for Sync committee contribution
	// - 1 ValidatorRegistrationPartialSig for Validator Registration
	// - 1 VoluntaryExitPartialSig for Voluntary Exit
	if err := mv.ValidPartialSigMessageCount(signerState, partialSignatureMessages); err != nil {
		return err
	}

	if role == types.RoleCommittee {

		// Rule: The number of signatures must be <= min(2*V, V + SYNC_COMMITTEE_SIZE) where V is the number of validators assigned to the cluster
		if !mv.ValidNumberOfSignaturesForCommitteeDuty(mv.Network.GetCommitteeInfo(msgID), partialSignatureMessages) {
			return ErrTooManyPartialSignatureMessages
		}

		// Rule: a ValidatorIndex can't appear more than 2 times in the []*PartialSignatureMessage list
		if mv.HasTripleValidatorOccurrence(partialSignatureMessages) {
			return ErrTripleValidatorIndexInPartialSignatures
		}
	} else if role == types.RoleSyncCommitteeContribution {
		// Rule: The number of signatures must be <= MaxSignaturesInSyncCommitteeContribution for the sync comittee contribution duty
		if len(partialSignatureMessages.Messages) > MaxSignatures {
			return ErrTooManyPartialSignatureMessages
		}
	} else {
		// Rule: The number of signatures must be 1 for the other types of duties
		if len(partialSignatureMessages.Messages) > 1 {
			return ErrTooManyPartialSignatureMessages
		}
	}

	return nil
}

// Auxiliary functions

// Check if the partial signature type is valid
func ValidPartialSignatureType(msgType types.PartialSigMsgType) bool {
	switch msgType {
	case types.PostConsensusPartialSig,
		types.RandaoPartialSig,
		types.SelectionProofPartialSig,
		types.ContributionProofs,
		types.ValidatorRegistrationPartialSig,
		types.VoluntaryExitPartialSig:
		return true
	default:
		return false
	}
}

// Check if the partial signature type matches a possible one for the role
func ExpectedPartialSignatureTypeForRole(msgType types.PartialSigMsgType, role types.RunnerRole) bool {
	switch role {
	case types.RoleCommittee:
		return msgType == types.PostConsensusPartialSig
	case types.RoleAggregator:
		return msgType == types.PostConsensusPartialSig || msgType == types.SelectionProofPartialSig
	case types.RoleProposer:
		return msgType == types.PostConsensusPartialSig || msgType == types.RandaoPartialSig
	case types.RoleSyncCommitteeContribution:
		return msgType == types.PostConsensusPartialSig || msgType == types.ContributionProofs
	case types.RoleValidatorRegistration:
		return msgType == types.ValidatorRegistrationPartialSig
	case types.RoleVoluntaryExit:
		return msgType == types.VoluntaryExitPartialSig
	default:
		return false
	}
}

// Check if a validator index belongs to a committee
func ValidatorIndexBelongsToCommittee(validatorIndex phase0.ValidatorIndex, committeeInfo *CommitteeInfo) bool {
	for _, committeeValidator := range committeeInfo.Validators {
		if validatorIndex == committeeValidator {
			return true
		}
	}
	return false
}

// Check if the counting of partial signature messages remains valid
func (mv *MessageValidator) ValidPartialSigMessageCount(signerState *SignerState, pSigMsgs *types.PartialSignatureMessages) error {
	return signerState.MessageCounts.ValidatePartialSignatureMessage(pSigMsgs.Type)
}

// Check if the number of signatures is valid, considering the committee role
func (mv *MessageValidator) ValidNumberOfSignaturesForCommitteeDuty(committeeInfo *CommitteeInfo, pSigMsgs *types.PartialSignatureMessages) bool {
	numValidators := len(committeeInfo.Validators)

	// Limit: min(2 * numValidators, numValidators + syncCommitteeSize)
	signaturesLimit := 2 * numValidators
	if signaturesLimit > numValidators+SyncCommitteeSize {
		signaturesLimit = numValidators + SyncCommitteeSize
	}

	return len(pSigMsgs.Messages) <= signaturesLimit
}

// Check if the partial signature message has a validator index related to 3 signatures
func (mv *MessageValidator) HasTripleValidatorOccurrence(pSigMsgs *types.PartialSignatureMessages) bool {
	validatorIndexCount := make(map[phase0.ValidatorIndex]int)
	for _, message := range pSigMsgs.Messages {
		validatorIndexCount[message.ValidatorIndex]++
		if validatorIndexCount[message.ValidatorIndex] > 2 {
			return true
		}
	}
	return false
}
