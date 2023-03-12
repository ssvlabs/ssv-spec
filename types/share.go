package types

// Share holds all info about the QBFT/ SSV Committee for msg signing and verification
type Share struct {
	OperatorID            OperatorID
	ValidatorPubKey       ValidatorPK `ssz-size:"48"`
	SharePubKey           []byte      `ssz-size:"48"`
	Committee             []*Operator `ssz-max:"13"`
	Quorum, PartialQuorum uint64
	DomainType            DomainType `ssz-size:"4"`
	FeeRecipientAddress   [20]byte   `ssz-size:"20"`
	Graffiti              []byte     `ssz-size:"32"`
}

// HasQuorum returns true if at least 2f+1 items are present (cnt is the number of items). It assumes nothing about those items, not their type or structure
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L259
func (share *Share) HasQuorum(cnt int) bool {
	return uint64(cnt) >= share.Quorum
}

// HasPartialQuorum returns true if at least f+1 items present (cnt is the number of items). It assumes nothing about those items, not their type or structure.
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L244
func (share *Share) HasPartialQuorum(cnt int) bool {
	return uint64(cnt) >= share.PartialQuorum
}

func (share *Share) Encode() ([]byte, error) {
	return share.MarshalSSZ()
}

func (share *Share) Decode(data []byte) error {
	return share.UnmarshalSSZ(data)
}
