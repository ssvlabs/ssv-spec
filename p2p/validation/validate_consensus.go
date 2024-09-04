package validation

import (
	"bytes"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Validates a consensus SignedSSVMessage on semantic rules
func (mv *MessageValidator) ValidateConsensusMessageSemantics(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	// Decode to qbft.Message
	qbftMessage := &qbft.Message{}
	err := qbftMessage.Decode(signedSSVMessage.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	signers := signedSSVMessage.OperatorIDs

	if len(signers) > 1 {
		// Rule: Decided msg with different type than Commit
		if qbftMessage.MsgType != qbft.CommitMsgType {
			return ErrNonDecidedWithMultipleSigners
		}

		// Rule: Number of signers must be >= quorum size
		if !ValidSignersLengthForCommitMessage(signers, mv.Network.GetCommitteeInfo(signedSSVMessage.SSVMessage.MsgID).OperatorIDs()) {
			return ErrDecidedNotEnoughSigners
		}
	}

	if len(signedSSVMessage.FullData) > 0 {
		// Rule: Prepare or commit messages must not have full data
		if (qbftMessage.MsgType == qbft.PrepareMsgType) ||
			(qbftMessage.MsgType == qbft.CommitMsgType && len(signers) == 1) {
			return ErrPrepareOrCommitWithFullData
		}

		// Rule: Full data hash must match root
		if !ValidFullDataRoot(signedSSVMessage.FullData, qbftMessage.Root) {
			return ErrInvalidHash
		}
	}

	// Rule: Consensus message type must be valid
	if !ValidConsensusMessageType(qbftMessage.MsgType) {
		return ErrUnknownQBFTMessageType
	}

	// Rule: Round must not be zero
	if qbftMessage.Round == qbft.NoRound {
		return ErrZeroRound
	}

	// Rule: consensus message must have the same identifier as the ssv message's identifier
	if !MatchedIdentifiers(qbftMessage.Identifier, signedSSVMessage.SSVMessage.MsgID[:]) {
		return ErrMismatchedIdentifier
	}

	return nil
}

// Validates a consensus SignedSSVMessage on QBFT logic rules
func (mv *MessageValidator) ValidateConsensusMessageQBFTLogic(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage, receivedAt time.Time) error {

	// Decode to qbft.Message
	qbftMessage := &qbft.Message{}
	err := qbftMessage.Decode(signedSSVMessage.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	peerState := mv.GetPeerState(peerID)
	dutyState := peerState.GetDutyState(signedSSVMessage.SSVMessage.MsgID, phase0.Slot(qbftMessage.Height))
	signers := signedSSVMessage.OperatorIDs

	if qbftMessage.MsgType == qbft.ProposalMsgType {
		// Rule: Signer must be the leader
		if !mv.IsLeader(signers[0], qbftMessage.Height, qbftMessage.Round) {
			return ErrSignerNotLeader
		}
	}

	if len(signers) > 1 {
		// Rule: Decided msg can't have the same signers as previously sent before for the same duty
		if HasSentSameDecidedSigners(dutyState, signedSSVMessage.OperatorIDs) {
			return ErrDecidedWithSameSigners
		}
	}

	if len(signers) == 1 {

		signerState := dutyState.GetSignerState(signers[0])

		if qbftMessage.Round == signerState.Round {
			// Rule: Peer must not send two proposals with different data
			if PeerHasSentProposalWithDifferentData(signerState, qbftMessage.MsgType, signedSSVMessage.FullData) {
				return ErrDuplicatedProposalWithDifferentData
			}

			// Rule: Peer must send only 1 proposal, 1 prepare, 1 commit and 1 round-change per round
			if err := ValidConsensusMessageCount(signerState, qbftMessage.MsgType, len(signers)); err != nil {
				return err
			}
		}

		// Rule: Round must not be smaller then current peer's round -1 or +1. Only for non-decided messages
		if err := mv.RoundBelongToAllowedSpread(phase0.Slot(qbftMessage.Height), qbftMessage.Round, receivedAt); err != nil {
			return err
		}

		// Rule: Ignore if peer already advanced to a later round. Only for non-decided messages
		if SignerAlreadyAdvancedRound(signerState, qbftMessage.Round) {
			return ErrRoundAlreadyAdvanced
		}
	}

	// Rule: Can only exist for Proposal or Round-Change messages
	if err := ValidRoundChangeJustificationForMessageType(qbftMessage); err != nil {
		return err
	}

	// Rule: Can only exist for Proposal messages
	if err := ValidPrepareJustificationForMessageType(qbftMessage); err != nil {
		return err
	}

	return nil
}

// Validates a consensus SignedSSVMessage on duty logic rules
func (mv *MessageValidator) ValidateConsensusMessageDutyLogic(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage, receivedAt time.Time) error {

	// Decode to qbft.Message
	qbftMessage := &qbft.Message{}
	err := qbftMessage.Decode(signedSSVMessage.SSVMessage.Data)
	if err != nil {
		return ErrUndecodableData
	}

	role := signedSSVMessage.SSVMessage.MsgID.GetRoleType()

	// Rule: Duty role has consensus (true except for ValidatorRegistration and VoluntaryExit)
	if !ValidRoleForConsensus(role) {
		return ErrUnexpectedConsensusMessage
	}

	// Rule: Round cut-offs for roles:
	// - 12 (committee and aggregation)
	// - 6 (other types)
	if !ValidRoundForRole(qbftMessage.Round, role) {
		return ErrRoundTooHigh
	}

	return nil
}

// Auxiliary functions

// Check if signer is the leader
func (mv *MessageValidator) IsLeader(signer types.OperatorID, height qbft.Height, round qbft.Round) bool {
	expectedProposer := mv.Config.ProposerForHeightAndRound(height, round)
	return (signer == expectedProposer)
}

// Check if the number of signers is valid (i.e. belongs to the [quorum, committee size] interval)
func ValidSignersLengthForCommitMessage(signers []types.OperatorID, committee []types.OperatorID) bool {
	f := int((len(committee) - 1) / 3)
	quorum := 2*f + 1
	return (len(signers) >= quorum && len(signers) <= len(committee))
}

// Check if the root is correct for the full data
func ValidFullDataRoot(fullData []byte, root [32]byte) bool {
	expectedRoot, err := qbft.HashDataRoot(fullData)
	if err != nil {
		return false
	}
	return bytes.Equal(root[:], expectedRoot[:])
}

// Check if the consensus message type is valid
func ValidConsensusMessageType(msgType qbft.MessageType) bool {
	switch msgType {
	case qbft.ProposalMsgType,
		qbft.PrepareMsgType,
		qbft.CommitMsgType,
		qbft.RoundChangeMsgType:
		return true
	default:
		return false
	}
}

// Check if the identifiers are equal
func MatchedIdentifiers(id1 []byte, id2 []byte) bool {
	return bytes.Equal(id1, id2)
}

// Check if the signers list of the decided message already exists in the state
func HasSentSameDecidedSigners(dutyState *DutyState, signers []types.OperatorID) bool {
	return dutyState.HasDecided(signers)
}

// Check if there's already another proposal stored with a different proposed data
func PeerHasSentProposalWithDifferentData(signerState *SignerState, msgType qbft.MessageType, fullData []byte) bool {
	if msgType != qbft.ProposalMsgType {
		return false
	}
	if signerState.MessageCounts.Proposal == 0 {
		return false
	}
	return !bytes.Equal(signerState.ProposalData, fullData)
}

// Check if the couting of consensus message is valid
func ValidConsensusMessageCount(signerState *SignerState, msgType qbft.MessageType, numSigners int) error {
	return signerState.MessageCounts.ValidateConsensusMessage(msgType, numSigners)
}

// Returns the estimated round for a given time duration after the slot start time
func currentEstimatedRound(sinceSlotStart time.Duration) qbft.Round {
	// Check if it's in the quick timeout
	if currentQuickRound := qbft.FirstRound + qbft.Round(sinceSlotStart/qbft.QuickTimeout); currentQuickRound <= qbft.QuickTimeoutThreshold {
		return currentQuickRound
	}

	// Compute the round considering the slow timeout phase
	sinceFirstSlowRound := sinceSlotStart - (time.Duration(qbft.QuickTimeoutThreshold) * qbft.QuickTimeout)
	estimatedRound := qbft.QuickTimeoutThreshold + qbft.FirstRound + qbft.Round(sinceFirstSlowRound/qbft.SlowTimeout)
	return estimatedRound
}

// Check if round belongs to allowed interval (i.e. [FirstRound, estimated round + allowed spread])
func (mv *MessageValidator) RoundBelongToAllowedSpread(slot phase0.Slot, round qbft.Round, receivedAt time.Time) error {

	// Get estimated round
	slotStartTime := time.Unix(mv.Beacon.EstimatedTimeAtSlot(slot), 0)
	sinceSlotStart := time.Duration(0)
	estimatedRound := qbft.FirstRound
	if receivedAt.After(slotStartTime) {
		sinceSlotStart = receivedAt.Sub(slotStartTime)
		estimatedRound = currentEstimatedRound(sinceSlotStart)
	}

	// Allowed interval
	lowestAllowed := qbft.FirstRound
	highestAllowed := estimatedRound + AllowedRoundsInFuture

	// Check
	if round < lowestAllowed || round > highestAllowed {
		return ErrRoundNotAllowedInTimeSpread
	}
	return nil
}

// Check if, in the peer's view, the signer has already advanced to a later round
func SignerAlreadyAdvancedRound(signerState *SignerState, round qbft.Round) bool {
	return (signerState.Round > round)
}

// Check if round justification is valid according to successfull mashalling and msg type
func ValidRoundChangeJustificationForMessageType(msg *qbft.Message) error {
	rcj, err := msg.GetRoundChangeJustifications()
	if err != nil {
		return ErrUndecodableData
	}

	// Rule: Can only exist for Proposal or Round-Change messages
	if len(rcj) != 0 && msg.MsgType != qbft.ProposalMsgType && msg.MsgType != qbft.RoundChangeMsgType {
		return ErrUnexpectedRoundChangeJustifications
	}

	return nil
}

// Check if prepare justification is valid according to successfull mashalling and msg type
func ValidPrepareJustificationForMessageType(msg *qbft.Message) error {
	pj, err := msg.GetPrepareJustifications()
	if err != nil {
		return ErrUndecodableData
	}

	// Rule: Can only exist for Proposal messages
	if len(pj) != 0 && msg.MsgType != qbft.ProposalMsgType {
		return ErrUnexpectedPrepareJustifications
	}

	return nil
}

// Check if role has a consensus phase
func ValidRoleForConsensus(role types.RunnerRole) bool {
	switch role {
	case types.RoleValidatorRegistration, types.RoleVoluntaryExit:
		return false
	default:
		return true
	}
}

// Check if round is valid for role
func ValidRoundForRole(round qbft.Round, role types.RunnerRole) bool {
	switch role {
	case types.RoleCommittee, types.RoleAggregator:
		return (round <= 12)
	case types.RoleProposer, types.RoleSyncCommitteeContribution:
		return (round <= 6)
	default:
		return false
	}
}
