package types

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID = uint64

// CommitteeMember represents an SSV operator node that is part of a committee
type CommitteeMember struct {
	OperatorID        OperatorID
	CommitteeID       CommitteeID `ssz-size:"32"`
	SSVOperatorPubKey []byte      `ssz-size:"294"`
	// FaultyNodes is the number of nodes that are considered faulty or malicious in the operator's committee
	FaultyNodes uint64
	// All the members of the committee
	Committee  []*Operator `ssz-max:"13"`
	DomainType DomainType  `ssz-size:"4"`
}

// Operator represents a node in the network that holds an ID and a public key
type Operator struct {
	OperatorID        OperatorID
	SSVOperatorPubKey []byte `ssz-size:"459"`
}

// HasQuorum returns true if at least 2f+1 items are present (cnt is the number of items). It assumes nothing about those items, not their type or structure
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L259

func (cm *CommitteeMember) HasQuorum(cnt int) bool {
	return uint64(cnt) >= cm.GetQuorum()
}

func (cm *CommitteeMember) GetQuorum() uint64 {
	return 2*cm.FaultyNodes + 1
}

// HasPartialQuorum returns true if at least f+1 items present (cnt is the number of items). It assumes nothing about those items, not their type or structure.
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L244

func (cm *CommitteeMember) HasPartialQuorum(cnt int) bool {
	return uint64(cnt) >= cm.FaultyNodes+1
}

func (cm *CommitteeMember) Encode() ([]byte, error) {
	return cm.MarshalSSZ()
}

func (cm *CommitteeMember) Decode(data []byte) error {
	return cm.UnmarshalSSZ(data)
}
