package validation

import (
	"slices"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Validates the syntax of a SignedSSVMessage
func (mv *MessageValidator) validateSyntax(msg *types.SignedSSVMessage) error {

	// Rule: SignedSSVMessage cannot be nil
	if msg == nil {
		return ErrNilSignedSSVMessage
	}

	// Rule: Must have at least one signer
	if len(msg.OperatorIDs) == 0 {
		return ErrNoSigners
	}

	// Rule: Must have at least one signature
	if len(msg.Signatures) == 0 {
		return ErrNoSignatures
	}

	// Rule: Signature size
	for _, sig := range msg.Signatures {
		if len(sig) != RsaSignatureSize {
			return ErrWrongRSASignatureSize
		}
	}

	// Rule: Signers must be sorted
	if !slices.IsSorted(msg.OperatorIDs) {
		return ErrSignersNotSorted
	}

	// Rule: Signer can't be zero
	for _, signer := range msg.OperatorIDs {
		if signer == 0 {
			return ErrZeroSigner
		}
	}

	// Rule: Signers must be unique
	var prevSigner types.OperatorID
	for _, signer := range msg.OperatorIDs {
		if signer == prevSigner {
			return ErrDuplicatedSigner
		}
		prevSigner = signer
	}

	// Rule: Len(Signers) must be equal to Len(Signatures)
	if len(msg.OperatorIDs) != len(msg.Signatures) {
		return ErrSignersAndSignaturesWithDifferentLength
	}

	// Rule: SSVMessage cannot be nil
	if msg.SSVMessage == nil {
		return ErrNilSSVMessage
	}

	// Rule: SSVMessage.Data must not be empty
	if len(msg.SSVMessage.Data) == 0 {
		return ErrEmptyData
	}

	// SSVMessage.Data must respect the size limit
	if len(msg.SSVMessage.Data) > MaxSSVMessageData {
		return ErrSSVMessageDataTooBig
	}

	switch msg.SSVMessage.MsgType {
	case types.SSVConsensusMsgType:
		// Rule: SSVMessage.Data decoding
		var qbftMessage qbft.Message
		if err := qbftMessage.Decode(msg.SSVMessage.Data); err != nil {
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
		if err := partialSignatureMessages.Decode(msg.SSVMessage.Data); err != nil {
			return ErrUndecodableData
		}
	}
	return nil
}
