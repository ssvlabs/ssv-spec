package types

type ValidatorCommittee struct {
	CommitteeID     CommitteeID       `ssz-size:"32"`
	ValidatorShares []*ValidatorShare `ssz-max:"13"`
}
