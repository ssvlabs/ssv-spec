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

func (share *Share) Encode() ([]byte, error) {
	return share.MarshalSSZ()
}

func (share *Share) Decode(data []byte) error {
	return share.UnmarshalSSZ(data)
}
