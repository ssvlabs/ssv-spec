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
// - SharePubKey must be 48 bytes
// - Committee must be non-empty and within the ssz-max bound
// - Committee members must be non-nil, have non-zero signer IDs, and 48-byte SharePubKeys
// - ValidatorPubKey must not be all-zero
// - Graffiti must be exactly 32 bytes
// - DomainType must be one of the known SSV domains in this spec
func (share *Share) Validate() error {
	if share == nil {
		return NewError(InvalidShareErrorCode, "nil share")
	}

	if share.ValidatorPubKey == (ValidatorPK{}) {
		return NewError(InvalidShareErrorCode, "zero validator pubkey not allowed")
	}

	if len(share.SharePubKey) != 48 {
		return NewError(InvalidShareErrorCode, "invalid share public key length")
	}

	if len(share.Committee) == 0 {
		return NewError(InvalidShareErrorCode, "empty committee")
	}
	if len(share.Committee) > MaxCommitteeSize {
		return NewError(InvalidShareErrorCode, "committee too large")
	}

	seenSigners := make(map[OperatorID]struct{}, len(share.Committee))
	for _, member := range share.Committee {
		if member == nil {
			return NewError(InvalidShareErrorCode, "nil committee member")
		}
		if member.Signer == 0 {
			return NewError(InvalidShareErrorCode, "committee member signer ID 0 not allowed")
		}
		if len(member.SharePubKey) != 48 {
			return NewError(InvalidShareErrorCode, "invalid committee member share public key length")
		}
		if _, exists := seenSigners[member.Signer]; exists {
			return NewError(InvalidShareErrorCode, "duplicate committee member signer")
		}
		seenSigners[member.Signer] = struct{}{}
	}

	if len(share.Graffiti) != 32 {
		return NewError(InvalidShareErrorCode, "invalid graffiti length")
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
