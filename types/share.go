package types

import "github.com/attestantio/go-eth2-client/spec/phase0"

// Share holds all info about the validator share
type Share struct {
	ValidatorIndex      phase0.ValidatorIndex
	ValidatorPubKey     ValidatorPK      `ssz-size:"48"`
	SharePubKey         ShareValidatorPK `ssz-size:"48"`
	Committee           []*ShareMember   `ssz-max:"13"`
	DomainType          DomainType       `ssz-size:"4"`
	FeeRecipientAddress [20]byte         `ssz-size:"20"`
	Graffiti            []byte           `ssz-size:"32"`
}

// ShareMember holds ShareValidatorPK and ValidatorIndex
type ShareMember struct {
	SharePubKey ShareValidatorPK `ssz-size:"48"`
	Signer      OperatorID
}

// Validate checks the following rules:
// - Committee must be non-empty
// - Committee members must have non-zero signer IDs with no duplicates
// - ValidatorPubKey must not be all-zero
// - DomainType must be one of the known SSV domains in this spec
func (share *Share) Validate() error {
	if share == nil {
		return NewError(InvalidShareErrorCode, "nil share")
	}

	if share.ValidatorPubKey == (ValidatorPK{}) {
		return NewError(InvalidShareErrorCode, "zero validator pubkey not allowed")
	}

	if len(share.Committee) == 0 {
		return NewError(InvalidShareErrorCode, "empty committee")
	}

	seenSigners := make(map[OperatorID]struct{}, len(share.Committee))
	for _, member := range share.Committee {
		var signer OperatorID
		if member != nil {
			signer = member.Signer
		}
		if signer == 0 {
			return NewError(InvalidShareErrorCode, "committee member signer ID 0 not allowed")
		}
		if _, exists := seenSigners[signer]; exists {
			return NewError(InvalidShareErrorCode, "duplicate committee member signer")
		}
		seenSigners[signer] = struct{}{}
	}

	if !share.DomainType.IsKnown() {
		return NewError(InvalidShareErrorCode, "unknown domain type")
	}

	return nil
}

func (share *Share) Encode() ([]byte, error) {
	return share.MarshalSSZ()
}

func (share *Share) Decode(data []byte) error {
	return share.UnmarshalSSZ(data)
}
