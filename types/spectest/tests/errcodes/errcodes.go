package errcodes

import "strings"

// Code represents a spec test error code.
type Code uint32

const (
	ErrNonUniqueSigner Code = iota + 1
	ErrInconsistentSigners
	ErrMessageSignerZero
	ErrNoPartialSignatureMessages
	ErrIncorrectEndOfOffset
	ErrUnknownDutyRole
	ErrIncorrectSize
	ErrVoluntaryExitNoConsensusData
	ErrValidatorRegistrationNoConsensusData
	ErrNilSSVMessage
	ErrNoSigners
	ErrSignersSignaturesLengthMismatch
	ErrSignerIDZero
	ErrNoSignatures
	ErrEmptySignature
	ErrUnknownBlockVersion
)

var errorStringToCode = map[string]Code{
	"non unique signer":                                        ErrNonUniqueSigner,
	"inconsistent signers":                                     ErrInconsistentSigners,
	"message invalid: signer ID 0 not allowed":                 ErrMessageSignerZero,
	"no PartialSignatureMessages messages":                     ErrNoPartialSignatureMessages,
	"could not unmarshal ssz: incorrect end of offset: 12 0":   ErrIncorrectEndOfOffset,
	"unknown duty role":                                        ErrUnknownDutyRole,
	"could not unmarshal ssz: incorrect size":                  ErrIncorrectSize,
	"voluntary exit has no consensus data":                     ErrVoluntaryExitNoConsensusData,
	"validator registration has no consensus data":             ErrValidatorRegistrationNoConsensusData,
	"nil SSVMessage":                                           ErrNilSSVMessage,
	"no signers":                                               ErrNoSigners,
	"number of signatures is different than number of signers": ErrSignersSignaturesLengthMismatch,
	"signer ID 0 not allowed":                                  ErrSignerIDZero,
	"no signatures":                                            ErrNoSignatures,
	"empty signature":                                          ErrEmptySignature,
}

// FromError returns the error code for a given error.
func FromError(err error) Code {
    if err == nil {
        return 0
    }
    msg := err.Error()
    if code, ok := errorStringToCode[msg]; ok {
        return code
    }
    if strings.HasPrefix(msg, "unknown block version") {
        return ErrUnknownBlockVersion
    }
    // Allow substring match for wrapped errors used across suites
    for key, code := range errorStringToCode {
        if strings.Contains(msg, key) {
            return code
        }
    }
    return 0
}
