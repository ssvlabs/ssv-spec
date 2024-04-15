package types

// Share holds all info about the validator share
type Share struct {
	ValidatorPubKey     ValidatorPK      `ssz-size:"48"`
	SharePubKey         ShareValidatorPK `ssz-size:"48"`
	Committee           []ShareMember    `ssz-max:"13"`
	Quorum              uint64
	DomainType          DomainType `ssz-size:"4"`
	FeeRecipientAddress [20]byte   `ssz-size:"20"`
	Graffiti            []byte     `ssz-size:"32"`
}

// ShareMember holds ShareValidatorPK and ValidatorIndex
type ShareMember struct {
	SharePubKey ShareValidatorPK `ssz-size:"48"`
	signer      OperatorID
}

// HasQuorum returns true if at least 2f+1 items are present (cnt is the number of items). It assumes nothing about those items, not their type or structure
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L259
func (share *Share) HasQuorum(cnt int) bool {
	return uint64(cnt) >= share.Quorum
}

func (share *Share) Encode() ([]byte, error) {
	return share.MarshalSSZ()
}

func (share *Share) Decode(data []byte) error {
	return share.UnmarshalSSZ(data)
}
