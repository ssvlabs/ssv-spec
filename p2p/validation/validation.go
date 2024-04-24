package validation

import (
	"context"
	"errors"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	MaxMsgSize                               = (4 + 56 + 8388668) * 1.1
	PartialSignatureSize                     = 48
	MessageSignatureSize                     = 256
	SyncCommitteeSize                        = 512
	MaxSignaturesInSyncCommitteeContribution = 13
)

func GetMessageValidationStateForPeer(peerID peer.ID) *PeerState {
	return nil
}

func GetPublicKeyFromOperator(operatorID types.OperatorID) []byte {
	return nil
}

func (mv *MessageValidation) ValidatePubsubMessage(_ context.Context, _ peer.ID, pmsg *pubsub.Message) pubsub.ValidationResult {

	var peerState *PeerState
	var err error

	peerState = GetMessageValidationStateForPeer(pmsg.ReceivedFrom)

	err = mv.Validate(pmsg)

	// Verify signature if validation chain is successful
	if err == nil {
		err = mv.VerifyMessageSignature(pmsg)
	}

	// Check error
	if err != nil {

		var valErr Error
		if errors.As(err, &valErr) {
			// Update state
			peerState.OnError(valErr)

			if valErr.Reject() {
				// Reject
				return pubsub.ValidationReject
			} else {
				// Ignore
				return pubsub.ValidationIgnore
			}
		} else {
			panic(err)
		}
	} else {
		return pubsub.ValidationAccept
	}
}

func (mv *MessageValidation) VerifyMessageSignature(pmsg *pubsub.Message) error {
	// Already verified
	signedSSVMessage := &types.SignedSSVMessage{}
	_ = signedSSVMessage.Decode(pmsg.Data)

	committeeMembers := make([]*types.CommitteeMember, 0)
	for _, signer := range signedSSVMessage.GetOperatorIDs() {
		signerPublicKey := GetPublicKeyFromOperator(signer)
		committeeMembers = append(committeeMembers, &types.CommitteeMember{
			OperatorID:        signer,
			SSVOperatorPubKey: signerPublicKey,
		})
	}

	err := mv.signatureVerifier.Verify(signedSSVMessage, committeeMembers)
	if err != nil {
		return ErrSignatureVerification
	}

	return nil
}

func (mv *MessageValidation) Validate(pmsg *pubsub.Message) error {
	// Validation Chain: Syntax -> Semantics -> QBFT Semantics | Partial Signature Semantics -> QBFT Logic -> Duty Rules
	return mv.ValidateSyntax(pmsg)
}

func (mv *MessageValidation) ValidateSyntax(pmsg *pubsub.Message) error {

	// Syntax validation

	/*
		Messages structures and checks (->)

			// PubSub Message
			type Message struct {
				*pb.Message
				ID
				ReceivedFrom
				ValidatorData
				Local
			}
			// pb.Message
				type Message struct {
				From
				Data -> Size limit, Decodable to SignedSSVMessage
				Seqno
				Topic
				Signature
				Key
				XXX_NoUnkeyedLiteral
				XXX_unrecognized
				XXX_sizecache
			}
			type SignedSSVMessage struct { -> Validate
				Signatures  [][]byte
				OperatorIDs []OperatorID
				SSVMessage  *SSVMessage
				FullData    []byte
			}
			type SSVMessage struct {
				MsgType MsgType
				MsgID   MessageID
				Data 	[]byte -> Decodable to qbft.Message or types.PartialSignatureMessages
			}
			type Message struct {
				MsgType    				 MessageType
				Height     				 Height
				Round      				 Round
				Identifier 				 []byte
				Root                     [32]byte
				DataRound                Round
				RoundChangeJustification [][]byte -> Decodable to []SignedSSVMessage
				PrepareJustification     [][]byte -> Decodable to []SignedSSVMessage
			}
	*/

	// Rule: Pubsub.Message.Message.Data size upper limit
	if float64(len(pmsg.Data)) > MaxMsgSize {
		return ErrPubSubDataTooBig
	}
	// Rule: Pubsub.Message.Message.Data decoding
	signedSSVMessage := &types.SignedSSVMessage{}
	if err := signedSSVMessage.Decode(pmsg.Data); err != nil {
		return ErrUndecodableData
	}

	// Rule: Invalid SignedSSVMessage
	if err := signedSSVMessage.Validate(); err != nil {
		return ErrInvalidSignedMessage
	}

	switch signedSSVMessage.SSVMessage.MsgType {
	case types.SSVConsensusMsgType:
		// Rule: SSVMessage.Data decoding
		var qbftMessage qbft.Message
		if err := qbftMessage.Decode(signedSSVMessage.SSVMessage.Data); err != nil {
			return ErrUndecodableData
		}

		// Rule: Message.RoundChangeJustification or Message.PrepareJustification decoding
		if _, err := qbftMessage.GetPrepareJustifications(); err != nil {
			return ErrUndecodableData
		}
		if _, err := qbftMessage.GetRoundChangeJustifications(); err != nil {
			return ErrUndecodableData
		}

	case types.SSVPartialSignatureMsgType:
		// Rule: SSVMessage.Data decoding
		var partialSignatureMessages types.PartialSignatureMessages
		if err := partialSignatureMessages.Decode(signedSSVMessage.SSVMessage.Data); err != nil {
			return ErrUndecodableData
		}
	}

	return mv.ValidateSemantics(pmsg.ReceivedFrom, signedSSVMessage, pmsg.GetTopic())
}

func (mv *MessageValidation) ValidateSemantics(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage, topic string) error {

	/*
		Messages structures and checks (->)

			type SignedSSVMessage struct {
				Signatures  [][]byte -> Signature size
				OperatorIDs []OperatorID -> Check if signers are in committee
				SSVMessage  *SSVMessage
				FullData    []byte
			}
			type SSVMessage struct {
				MsgType MsgType -> Can't be event, Can't be DKG, Must be valid (known)
				MsgID   MessageID -> Equal domain, Valid role (known), Active Validator, Validator Liquidated, ClusterID exists, Correct topic
				Data 	[]byte
			}
	*/

	// Rule: Signature size
	for _, sig := range signedSSVMessage.Signatures {
		if len(sig) != MessageSignatureSize {
			return ErrWrongSignatureSize
		}
	}

	// Rule: Signers must belong to validator committee or clusterID -
	for _, signer := range signedSSVMessage.GetOperatorIDs() {
		if !mv.SignerBelongsToCommittee(signer, signedSSVMessage.SSVMessage.MsgID) {
			return ErrSignerNotInCommittee
		}
	}

	if signedSSVMessage.SSVMessage.MsgType == types.SSVEventMsgType {
		// Rule: Event message
		return ErrEventMessage
	} else if signedSSVMessage.SSVMessage.MsgType == types.DKGMsgType {
		// Rule: DKG message
		return ErrDKGMessage
	} else if signedSSVMessage.SSVMessage.MsgType != types.SSVConsensusMsgType && signedSSVMessage.SSVMessage.MsgType != types.SSVPartialSignatureMsgType {
		// Rule: Unknown msg type
		return ErrUnknownSSVMessageType
	}

	// Rule: If domain is different then self domain
	domain := signedSSVMessage.SSVMessage.MsgID.GetDomain()
	if !mv.ValidDomain(domain) {
		return ErrWrongDomain
	}

	// Rule: If role is invalid
	role := signedSSVMessage.SSVMessage.MsgID.GetRoleType()
	if !mv.ValidRole(role) {
		return ErrInvalidRole
	}

	senderID := signedSSVMessage.SSVMessage.MsgID.GetSenderID()
	if role != types.RoleCommittee {
		validatorPK := senderID

		// Rule: If validator is not active
		if !mv.ActiveValidator(validatorPK) {
			return ErrValidatorNotAttesting
		}

		// Rule: If validator is liquidated
		if mv.ValidatorLiquidated(validatorPK) {
			return ErrValidatorLiquidated
		}
	} else {
		// Rule: Cluster does not exist
		if !mv.ExistingClusterID(senderID) {
			return ErrNonExistentClusterID
		}
	}

	// Rule: Check if message was sent in the correct topic
	if !mv.CorrectTopic(signedSSVMessage.SSVMessage.MsgID, topic) {
		return ErrIncorrectTopic
	}

	// Call specific semantics rules
	if signedSSVMessage.SSVMessage.MsgType == types.SSVConsensusMsgType {
		return mv.ValidateConsensusMessageSemantics(peerID, signedSSVMessage)
	} else if signedSSVMessage.SSVMessage.MsgType == types.SSVPartialSignatureMsgType {
		return mv.ValidatePartialSignatureMessageSemantics(peerID, signedSSVMessage)
	} else {
		panic("type should have been checked to be either consensus or partial signature")
	}

	return nil
}

func (mv *MessageValidation) ValidateConsensusMessageSemantics(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	/*
		Messages structures and checks (->)

			type SignedSSVMessage struct {
				Signatures  [][]byte
				OperatorIDs []OperatorID -> Valid length for decided, Length 1 for non-decided
				SSVMessage  *SSVMessage
				FullData    []byte -> Must be empty for prepare and commit (with 1 signer)
			}
			type SSVMessage struct {
				MsgType MsgType
				MsgID   MessageID
				Data 	[]byte
			}
			type Message struct {
				MsgType    				 MessageType -> Valid (known), Must be of type commit for decided message
				Height     				 Height
				Round      				 Round
				Identifier 				 []byte -> Match SSVMessage.MsgID
				Root                     [32]byte -> Must be hash of FullData
				DataRound                Round
				RoundChangeJustification [][]byte
				PrepareJustification     [][]byte
			}
	*/

	// Already verified
	var qbftMessage qbft.Message
	_ = qbftMessage.Decode(signedSSVMessage.SSVMessage.Data)

	signers := signedSSVMessage.GetOperatorIDs()

	if len(signers) > 1 {
		// Rule: Decided msg with different type than Commit
		if qbftMessage.MsgType != qbft.CommitMsgType {
			return ErrNonDecidedWithMultipleSigners
		}

		// Rule: Number of signers must be between quorum and committee size
		if !mv.ValidSignersLengthForCommitMessage(signers) {
			return ErrWrongSignersLength
		}
	} else if len(signers) != 1 {
		// Rule: Non-decided messages must have one signer
		return ErrNoSigners
	}

	if len(signedSSVMessage.FullData) > 0 {
		// Rule: Prepare or commit messages must not have full data
		if (qbftMessage.MsgType == qbft.PrepareMsgType) &&
			(qbftMessage.MsgType == qbft.CommitMsgType && len(signers) == 1) {
			return ErrPrepareOrCommitWithFullData
		}

		// Rule: Full data hash must match root
		if !mv.ValidFullDataRoot(signedSSVMessage.FullData, qbftMessage.Root) {
			return ErrInvalidHash
		}
	}

	// Rule: Consensus message type must be valid
	if !mv.ValidConsensusMessageType(qbftMessage.MsgType) {
		return ErrUnknownQBFTMessageType
	}

	// Rule: consensus message must have the same identifier as the ssv message's identifier
	if !mv.MatchedIdentifiers(qbftMessage.Identifier, signedSSVMessage.SSVMessage.MsgID[:]) {
		return ErrMismatchedIdentifier
	}

	return mv.ValidateQBFTLogic(peerID, signedSSVMessage)
}

func (mv *MessageValidation) ValidateQBFTLogic(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	/*
		Messages structures and checks (->)

			type SignedSSVMessage struct {
				Signatures  [][]byte
				OperatorIDs []OperatorID -> Must be leader if message is Proposal, Decided msg must have more signers than a previous decided
				SSVMessage  *SSVMessage
				FullData    []byte
			}
			type SSVMessage struct {
				MsgType MsgType
				MsgID   MessageID
				Data 	[]byte
			}
			type Message struct {
				MsgType    				 MessageType -> Message count rules
				Height     				 Height
				Round      				 Round -> Must belong to round spread
				Identifier 				 []byte
				Root                     [32]byte
				DataRound                Round
				RoundChangeJustification [][]byte
				PrepareJustification     [][]byte
			}
	*/

	// Already verified
	var qbftMessage qbft.Message
	_ = qbftMessage.Decode(signedSSVMessage.SSVMessage.Data)

	signers := signedSSVMessage.GetOperatorIDs()

	if qbftMessage.MsgType == qbft.ProposalMsgType {
		// Rule: Signer must be the leader
		if !mv.IsLeader(signers[0], qbftMessage.Height, qbftMessage.Round) {
			return ErrSignerNotLeader
		}
	}

	if len(signers) > 1 {
		// Rule: Decided must have more signers than previous messages
		if !mv.NewDecidedHaveMoreSigners(peerID, signedSSVMessage) {
			return ErrDecidedWithLessSignersThanPrevious
		}
	}

	if len(signers) == 1 {
		// Rule: peer must send only 1 proposal, 1 prepare, 1 commit and 1 round-change per round
		if err := mv.ValidConsensusMessageCount(peerID, signedSSVMessage.SSVMessage.MsgID, qbftMessage.Height, qbftMessage.MsgType, qbftMessage.Round); err != nil {
			return err
		}
	}

	// Rule: Round must not be smaller then current peer's round -1 or +1
	if err := mv.RoundBelongToAllowedSpread(peerID, qbftMessage.Round, signedSSVMessage.SSVMessage.MsgID); err != nil {
		return err
	}

	return mv.ValidateQBFTMessageByDutyRules(peerID, signedSSVMessage)
}

func (mv *MessageValidation) ValidateQBFTMessageByDutyRules(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	/*
		Messages structures and checks (->)

			type SignedSSVMessage struct {
				Signatures  [][]byte
				OperatorIDs []OperatorID
				SSVMessage  *SSVMessage
				FullData    []byte
			}
			type SSVMessage struct {
				MsgType MsgType
				MsgID   MessageID -> Role must have consensus, Must be assigned to duty if role is Proposal or Sync Committee Aggregator
				Data 	[]byte
			}
			type Message struct {
				MsgType    				 MessageType
				Height     				 Height -> Must belong to spread for role, Satisfies a maximum number of duties per epoch for role
				Round      				 Round -> Must be below cut-off for role
				Identifier 				 []byte
				Root                     [32]byte
				DataRound                Round
				RoundChangeJustification [][]byte
				PrepareJustification     [][]byte
			}
	*/

	// Already verified
	var qbftMessage qbft.Message
	_ = qbftMessage.Decode(signedSSVMessage.SSVMessage.Data)

	// Rule: Duty role has consensus (true except for ValidatorRegistration and VoluntaryExit)
	if !mv.ValidRoleForConsensus(signedSSVMessage.SSVMessage.MsgID.GetRoleType()) {
		return ErrUnexpectedConsensusMessage
	}

	// Rule: For a proposal duty message, we check if the validator is assigned to it
	if signedSSVMessage.SSVMessage.MsgID.GetRoleType() == types.RoleProposer {
		if !mv.HasProposerDuty(signedSSVMessage.SSVMessage.MsgID.GetSenderID(), phase0.Slot(qbftMessage.Height)) {
			return ErrNoDuty
		}
	}

	// Rule: For a sync committee aggregation duty message, we check if the validator is assigned to it
	if signedSSVMessage.SSVMessage.MsgID.GetRoleType() == types.RoleSyncCommitteeContribution {
		if !mv.HasSyncCommitteeDuty(signedSSVMessage.SSVMessage.MsgID.GetSenderID(), phase0.Slot(qbftMessage.Height)) {
			return ErrNoDuty
		}
	}

	// Rule: current slot(height) must be between duty's starting slot and:
	// - duty's starting slot + 34 (committee and aggregation)
	// - duty's starting slot + 3 (other types)
	if err := mv.ValidDutySlot(peerID, phase0.Slot(qbftMessage.Height), signedSSVMessage.SSVMessage.MsgID.GetRoleType()); err != nil {
		return err
	}

	// Rule: valid number of duties per epoch:
	// - 2 for aggregation, voluntary exit and validator registration
	// - 2*V for Committee duty (where V is the number of validators in the cluster) (if no validator is doing sync committee in this epoch)
	// - else, accept
	if !mv.ValidNumberOfCommitteeDutiesPerEpoch(peerID, signedSSVMessage.SSVMessage.MsgID, phase0.Slot(qbftMessage.Height)) {
		return ErrTooManyDutiesPerEpoch
	}

	// Rule: Round cut-offs for roles:
	// - 12 (committee and aggregation)
	// - 6 (other types)
	if !mv.ValidRoundForRole(qbftMessage.Round, signedSSVMessage.SSVMessage.MsgID.GetRoleType()) {
		return ErrRoundTooHigh
	}

	return nil
}

func (mv *MessageValidation) ValidatePartialSignatureMessageSemantics(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	/*
		Messages structures and checks (->)

			type SignedSSVMessage struct {
				Signatures  [][]byte
				OperatorIDs []OperatorID -> Must have len = 1
				SSVMessage  *SSVMessage
				FullData    []byte -> Must be empty
			}
			type SSVMessage struct {
				MsgType MsgType
				MsgID   MessageID
				Data 	[]byte
			}
			type PartialSignatureMessages struct {
				Type     PartialSigMsgType -> Aligned to MsgID.Role
				Slot     phase0.Slot
				Messages []*PartialSignatureMessage -> Mut not be empty
			}

			type PartialSignatureMessage struct {
				PartialSignature Signature -> Correct size
				SigningRoot      [32]byte
				Signer         	 OperatorID -> Consistent with SignedSSVMessage.OperatorIDs[0]
				ValidatorIndex 	 phase0.ValidatorIndex -> Must belongs to MsgID.Committee or Validator
			}
	*/

	// Already verified
	var partialSignatureMessages types.PartialSignatureMessages
	_ = partialSignatureMessages.Decode(signedSSVMessage.SSVMessage.Data)

	// Rule: Partial Signature message must have 1 signer
	signers := signedSSVMessage.GetOperatorIDs()
	if len(signers) != 1 {
		return ErrPartialSigOneSigner
	}
	signer := signers[0]

	// Rule: Partial signature message must not have full data
	if len(signedSSVMessage.FullData) > 0 {
		return ErrFullDataNotInConsensusMessage
	}

	// Rule: Partial signature type must match expected type:
	// - PostConsensusPartialSig, for Committee duty
	// - RandaoPartialSig or PostConsensusPartialSig for Proposer
	// - SelectionProofPartialSig or PostConsensusPartialSig for Aggregator
	// - SelectionProofPartialSig or PostConsensusPartialSig for Sync committee contribution
	// - ValidatorRegistrationPartialSig for Validator Registration
	// - VoluntaryExitPartialSig for Voluntary Exit
	if !mv.ExpectedPartialSignatureTypeForRole(partialSignatureMessages.Type, signedSSVMessage.SSVMessage.MsgID) {
		return ErrPartialSignatureTypeRoleMismatch
	}

	// Rule: Partial signature message must have at least one signature
	if len(partialSignatureMessages.Messages) == 0 {
		return ErrEmptySignatures
	}

	for _, psigMsg := range partialSignatureMessages.Messages {
		// Rule: Partial signature must have expected length
		if len(psigMsg.PartialSignature) != PartialSignatureSize {
			return ErrWrongSignatureSize
		}
		// Rule: Partial signature signer must be consistent
		if psigMsg.Signer != signer {
			return ErrInconsistentSigners
		}
		// Rule: Validator index must match with validatorPK or one of clusterID's validators
		if !mv.ValidatorIndexBelongsToCommittee(psigMsg.ValidatorIndex, signedSSVMessage.SSVMessage.MsgID) {
			return ErrValidatorIndexMismatch
		}
	}

	return mv.ValidatePartialSigMessagesByDutyRules(peerID, signedSSVMessage)
}

func (mv *MessageValidation) ValidatePartialSigMessagesByDutyRules(peerID peer.ID, signedSSVMessage *types.SignedSSVMessage) error {

	/*
		Messages structures and checks (->)

			type SignedSSVMessage struct {
				Signatures  [][]byte
				OperatorIDs []OperatorID
				SSVMessage  *SSVMessage
				FullData    []byte
			}
			type SSVMessage struct {
				MsgType MsgType
				MsgID   MessageID -> Must be assigned to duty if role is Proposal or Sync Committee Aggregator
				Data 	[]byte
			}
			type PartialSignatureMessages struct {
				Type     PartialSigMsgType -> Message count rules
				Slot     phase0.Slot -> Must belong to allowed spread, Satisfies a maximum number of duties per epoch for role
				Messages []*PartialSignatureMessage -> Valid number of signatures (3 cases: committee duty, sync committee contribution, others)
			}

			type PartialSignatureMessage struct {
				PartialSignature Signature
				SigningRoot      [32]byte
				Signer         	 OperatorID
				ValidatorIndex 	 phase0.ValidatorIndex -> Can't appear more than 2 times for role Committee
			}
	*/

	// Already verified
	var partialSignatureMessages types.PartialSignatureMessages
	_ = partialSignatureMessages.Decode(signedSSVMessage.SSVMessage.Data)

	// Rule: For a proposal duty message, we check if the validator is assigned to it
	if signedSSVMessage.SSVMessage.MsgID.GetRoleType() == types.RoleProposer {
		if !mv.HasProposerDuty(signedSSVMessage.SSVMessage.MsgID.GetSenderID(), partialSignatureMessages.Slot) {
			return ErrNoDuty
		}
	}

	// Rule: For a sync committee aggregation duty message, we check if the validator is assigned to it
	if signedSSVMessage.SSVMessage.MsgID.GetRoleType() == types.RoleSyncCommitteeContribution {
		if !mv.HasSyncCommitteeDuty(signedSSVMessage.SSVMessage.MsgID.GetSenderID(), partialSignatureMessages.Slot) {
			return ErrNoDuty
		}
	}

	// Rule: peer must send only:
	// - 1 PostConsensusPartialSig, for Committee duty
	// - 1 RandaoPartialSig and 1 PostConsensusPartialSig for Proposer
	// - 1 SelectionProofPartialSig and 1 PostConsensusPartialSig for Aggregator
	// - 1 SelectionProofPartialSig and 1 PostConsensusPartialSig for Sync committee contribution
	// - 1 ValidatorRegistrationPartialSig for Validator Registration
	// - 1 VoluntaryExitPartialSig for Voluntary Exit
	if err := mv.ValidPartialSigMessageCount(peerID, signedSSVMessage.SSVMessage.MsgID.GetRoleType(), &partialSignatureMessages); err != nil {
		return err
	}

	// Rule: current slot must be between duty's starting slot and:
	// - duty's starting slot + 34 (committee and aggregation)
	// - duty's starting slot + 3 (other duties)
	if err := mv.ValidDutySlot(peerID, partialSignatureMessages.Slot, signedSSVMessage.SSVMessage.MsgID.GetRoleType()); err != nil {
		return err
	}

	// Rule: valid number of duties per epoch:
	// - 2 for aggregation, voluntary exit and validator registration
	// - 2*V for Committee duty (where V is the number of validators in the cluster) (if no validator is doing sync committee in this epoch)
	// - else, accept
	if !mv.ValidNumberOfCommitteeDutiesPerEpoch(peerID, signedSSVMessage.SSVMessage.MsgID, partialSignatureMessages.Slot) {
		return ErrTooManyDutiesPerEpoch
	}

	if signedSSVMessage.SSVMessage.MsgID.GetRoleType() == types.RoleCommittee {

		// Rule: The number of signatures must be <= min(2*V, V + SYNC_COMMITTEE_SIZE) where V is the number of validators assigned to the cluster
		if !mv.ValidNumberOfSignaturesForCommitteeDuty(signedSSVMessage.SSVMessage.MsgID.GetSenderID(), &partialSignatureMessages) {
			return ErrTooManySignatures
		}

		// Rule: a ValidatorIndex can't appear more than 2 times in the []*PartialSignatureMessage list
		if !mv.NoTripleValidatorOccurrence(&partialSignatureMessages) {
			return ErrTripleValidatorIndexInSignatures
		}
	} else if signedSSVMessage.SSVMessage.MsgID.GetRoleType() == types.RoleSyncCommitteeContribution {
		// Rule: The number of signatures must be <= MaxSignaturesInSyncCommitteeContribution for the sync comittee contribution duty
		if len(partialSignatureMessages.Messages) > MaxSignaturesInSyncCommitteeContribution {
			return ErrTooManySignatures
		}
	} else {
		// Rule: The number of signatures must be 1 for the other types of duties
		if len(partialSignatureMessages.Messages) > 1 {
			return ErrTooManySignatures
		}
	}

	return nil
}
