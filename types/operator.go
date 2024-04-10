package types

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID = uint64

// Operator represents an SSV operator node that is part of a committee
type Operator struct {
	OperatorID        OperatorID
	SSVOperatorPubKey []byte `ssz-size:"294"`
	// TODO: change with one parameter F
	Quorum, PartialQuorum uint32
	Committee             []*CommitteeMember `ssz-max:"13"`
}

// CommitteeMember represents all data in order to verify a committee member's identity
type CommitteeMember struct {
	OperatorID        OperatorID
	SSVOperatorPubKey []byte
}
