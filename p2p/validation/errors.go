package validation

import (
	"strings"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
)

type Error struct {
	text   string
	reject bool
}

// Implement error interface
func (e Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.text)

	return sb.String()
}

func handleValidationError(err error) pubsub.ValidationResult {
	var valErr Error
	if !errors.As(err, &valErr) {
		return pubsub.ValidationIgnore
	}
	if !valErr.reject {
		return pubsub.ValidationIgnore
	}
	return pubsub.ValidationReject
}

// List of errors
var (
	// PubSub errors
	ErrPubSubMessageHasNoData = Error{text: "pub-sub message has no data", reject: true}
	ErrPubSubDataTooBig       = Error{text: "pub-sub message data too big"}
	ErrMalformedPubSubMessage = Error{text: "pub-sub message is malformed", reject: true}

	// Syntax errors
	ErrNilSignedSSVMessage                     = Error{text: "decoded SignedSSVMessage is nil", reject: true}
	ErrNoSigners                               = Error{text: "no signers", reject: true}
	ErrNoSignatures                            = Error{text: "no signatures", reject: true}
	ErrWrongRSASignatureSize                   = Error{text: "wrong RSA signature size", reject: true}
	ErrSignersNotSorted                        = Error{text: "signers are not sorted", reject: true}
	ErrZeroSigner                              = Error{text: "zero signer ID", reject: true}
	ErrDuplicatedSigner                        = Error{text: "signer is duplicated", reject: true}
	ErrSignersAndSignaturesWithDifferentLength = Error{text: "signers and signatures with different length", reject: true}
	ErrNilSSVMessage                           = Error{text: "SSVMessage is nil", reject: true}
	ErrEmptyData                               = Error{text: "SSVMessage.Data is empty", reject: true}
	ErrSSVMessageDataTooBig                    = Error{text: "SSVMessage.Data is too big"}
	ErrUndecodableData                         = Error{text: "undecodable data", reject: true}

	// General semantic rules
	ErrSignerNotInCommittee   = Error{text: "signer is not in committee", reject: true}
	ErrWrongDomain            = Error{text: "wrong domain"}
	ErrInvalidRole            = Error{text: "invalid role", reject: true}
	ErrUnknownValidator       = Error{text: "validator does not exist"}
	ErrValidatorNotAttesting  = Error{text: "validator is not attesting"}
	ErrValidatorLiquidated    = Error{text: "validator is liquidated"}
	ErrNonExistentCommitteeID = Error{text: "non existent committee ID"}
	ErrIncorrectTopic         = Error{text: "incorrect topic"}
	ErrDKGMessage             = Error{text: "DKG messages are not supported", reject: true}
	ErrUnknownSSVMessageType  = Error{text: "unknown SSV message type", reject: true}

	// Consensus semantic rules
	ErrNonDecidedWithMultipleSigners = Error{text: "non-decided with multiple signers", reject: true}
	ErrDecidedNotEnoughSigners       = Error{text: "decided signers size is less than quorum size", reject: true}
	ErrPrepareOrCommitWithFullData   = Error{text: "prepare or commit with full data", reject: true}
	ErrInvalidHash                   = Error{text: "root doesn't match full data hash", reject: true}
	ErrUnknownQBFTMessageType        = Error{text: "unknown QBFT message type", reject: true}
	ErrZeroRound                     = Error{text: "round is zero", reject: true}
	ErrMismatchedIdentifier          = Error{text: "message ID mismatched", reject: true}

	// Consensus QBFT Logic rules
	ErrSignerNotLeader                     = Error{text: "signer is not leader", reject: true}
	ErrDecidedWithSameSigners              = Error{text: "decided with the same signers as sent before"}
	ErrDuplicatedProposalWithDifferentData = Error{text: "duplicated proposal with different data", reject: true}
	ErrDuplicatedMessage                   = Error{text: "message is duplicated", reject: true}
	ErrRoundNotAllowedInTimeSpread         = Error{text: "round is out of the interval [FirstRound, estimated round + 1]"}
	ErrRoundAlreadyAdvanced                = Error{text: "signer has already advanced to a later round"}
	ErrUnexpectedRoundChangeJustifications = Error{text: "message has a round-change justification but it's not a proposal or round-change", reject: true}
	ErrUnexpectedPrepareJustifications     = Error{text: "message has a prepare justification but it's not a proposal", reject: true}

	// Consensus duty Logic rules
	ErrUnexpectedConsensusMessage = Error{text: "unexpected consensus message for this role", reject: true}
	ErrRoundTooHigh               = Error{text: "round is too high for this role", reject: true}

	// Partial signature semantic rules
	ErrPartialSigMultipleSigners        = Error{text: "partial signature message with len(signers) != 1", reject: true}
	ErrFullDataNotInConsensusMessage    = Error{text: "full data in message different than consensus", reject: true}
	ErrInvalidPartialSignatureType      = Error{text: "invalid partial signature type", reject: true}
	ErrPartialSignatureTypeRoleMismatch = Error{text: "partial signature type and role don't match", reject: true}
	ErrNoPartialSignatureMessages       = Error{text: "no partial signature messages", reject: true}
	ErrWrongBLSSignatureSize            = Error{text: "wrong BLS signature size", reject: true}
	ErrInconsistentSigners              = Error{text: "inconsistent signers", reject: true}
	ErrValidatorIndexMismatch           = Error{text: "validator index mismatch"}

	// Partial signature duty logic rules
	ErrInvalidPartialSignatureTypeCount        = Error{text: "sent more partial signature messages of a certain type than allowed", reject: true}
	ErrTooManyPartialSignatureMessages         = Error{text: "too many signatures for cluster in partial signature message"}
	ErrTripleValidatorIndexInPartialSignatures = Error{text: "validator index appear 3 times in partial signature message", reject: true}

	// General duty logic rules
	ErrSlotAlreadyAdvanced   = Error{text: "signer already advanced to later slot"}
	ErrEarlySlotMessage      = Error{text: "message was sent before slot starts"}
	ErrLateSlotMessage       = Error{text: "current time is above duty's start +34(committee and aggregator) or +3(else) slots"}
	ErrNoDuty                = Error{text: "no duty for this epoch"}
	ErrTooManyDutiesPerEpoch = Error{text: "too many duties per epoch"}

	ErrSignatureVerification = Error{text: "signature verification", reject: true}
)
