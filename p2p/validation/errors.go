package validation

import (
	"strings"
)

type Error struct {
	text   string
	reject bool
}

func (e Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.text)

	return sb.String()
}

func (e Error) Reject() bool {
	return e.reject
}

func (e Error) Text() string {
	return e.text
}

var (
	ErrPubSubMessageHasNoData = Error{text: "pub-sub message has no data", reject: true}
	ErrPubSubDataTooBig       = Error{text: "pub-sub message data too big"}
	ErrUndecodableData        = Error{text: "undecodable data", reject: true}
	ErrSignatureVerification  = Error{text: "signature verification", reject: true}

	ErrIncorrectTopic                          = Error{text: "incorrect topic"}
	ErrWrongDomain                             = Error{text: "wrong domain"}
	ErrNonExistentClusterID                    = Error{text: "non existent cluster ID"}
	ErrValidatorLiquidated                     = Error{text: "validator is liquidated"}
	ErrValidatorNotAttesting                   = Error{text: "validator is not attesting"}
	ErrInvalidRole                             = Error{text: "invalid role", reject: true}
	ErrNoSigners                               = Error{text: "no signers", reject: true}
	ErrMoreSignersThanCommitteeSize            = Error{text: "more signers than committee size", reject: true}
	ErrZeroSigner                              = Error{text: "zero signer ID", reject: true}
	ErrDuplicatedSigner                        = Error{text: "signer is duplicated", reject: true}
	ErrSignersNotSorted                        = Error{text: "signers are not sorted", reject: true}
	ErrSignersAndSignaturesWithDifferentLength = Error{text: "signers and signatures with different length", reject: true}
	ErrNoSignatures                            = Error{text: "no signatures", reject: true}
	ErrWrongRSASignatureSize                   = Error{text: "wrong RSA signature size", reject: true}
	ErrWrongBLSSignatureSize                   = Error{text: "wrong BLS signature size", reject: true}
	ErrSignerNotInCommittee                    = Error{text: "signer is not in committee", reject: true}
	ErrUnknownSSVMessageType                   = Error{text: "unknown SSV message type", reject: true}
	ErrEventMessage                            = Error{text: "event messages are not broadcast", reject: true}
	ErrDKGMessage                              = Error{text: "DKG messages are not supported", reject: true}

	ErrUnknownQBFTMessageType              = Error{text: "unknown QBFT message type", reject: true}
	ErrRoundTooHigh                        = Error{text: "round is too high for this role"}
	ErrDecidedWithSameNumberOfSigners      = Error{text: "decided with the same amount of signers sent before"}
	ErrPrepareOrCommitWithFullData         = Error{text: "prepare or commit with full data", reject: true}
	ErrFullDataNotInConsensusMessage       = Error{text: "full data in message different than consensus", reject: true}
	ErrMismatchedIdentifier                = Error{text: "message ID mismatched"}
	ErrUnexpectedConsensusMessage          = Error{text: "unexpected consensus message for this role", reject: true}
	ErrSignerNotLeader                     = Error{text: "signer is not leader", reject: true}
	ErrInvalidHash                         = Error{text: "root doesn't match full data hash", reject: true}
	ErrNonDecidedWithMultipleSigners       = Error{text: "non-decided with multiple signers", reject: true}
	ErrWrongSignersLength                  = Error{text: "decided signers size is not between quorum and committee size", reject: true}
	ErrRoundAlreadyAdvanced                = Error{text: "signer has already advanced to a later round"}
	ErrDuplicatedProposalWithDifferentData = Error{text: "duplicated proposal with different data", reject: true}
	ErrInvalidRound                        = Error{text: "invalid round", reject: true}

	ErrPartialSigOneSigner              = Error{text: "partial signature message with len(signers) != 1", reject: true}
	ErrTooManySignatures                = Error{text: "too many signatures for cluster in partial signature message"}
	ErrTripleValidatorIndexInSignatures = Error{text: "validator index appear 3 times in partial signature message", reject: true}
	ErrNoPartialSignatureMessages       = Error{text: "no partial signature messages", reject: true}
	ErrInconsistentSigners              = Error{text: "inconsistent signers", reject: true}
	ErrValidatorIndexMismatch           = Error{text: "validator index mismatch"}
	ErrInvalidPartialSignatureType      = Error{text: "invalid partial signature type", reject: true}
	ErrPartialSignatureTypeRoleMismatch = Error{text: "partial signature type and role don't match", reject: true}

	ErrTooManyDutiesPerEpoch = Error{text: "too many duties per epoch"}
	ErrNoDuty                = Error{text: "no duty for this epoch"}
)
