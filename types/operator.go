package types

import "github.com/bloxapp/ssv-spec/ssv"

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID = uint64

// Operator represents an SSV operator node that is part of a committee
type Operator struct {
	OperatorID        OperatorID
	ClusterID         ssv.ClusterID
	SSVOperatorPubKey []byte `ssz-size:"294"`
	// TODO: change with one parameter F
	Quorum, PartialQuorum uint64
	// All the members of the committee
	Committee []*CommitteeMember `ssz-max:"13"`
}

// CommitteeMember represents all data in order to verify a committee member's identity
type CommitteeMember struct {
	OperatorID        OperatorID
	SSVOperatorPubKey []byte
}

// HasQuorum returns true if at least 2f+1 items are present (cnt is the number of items). It assumes nothing about those items, not their type or structure
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L259
func (operator *Operator) HasQuorum(cnt int) bool {
	return uint64(cnt) >= operator.Quorum
}

// HasPartialQuorum returns true if at least f+1 items present (cnt is the number of items). It assumes nothing about those items, not their type or structure.
// https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/main/dafny/spec/L1/node_auxiliary_functions.dfy#L244
func (operator *Operator) HasPartialQuorum(cnt int) bool {
	return uint64(cnt) >= operator.PartialQuorum
}

func (operator *Operator) Encode() ([]byte, error) {
	return operator.MarshalSSZ()
}

func (operator *Operator) Decode(data []byte) error {
	return operator.UnmarshalSSZ(data)
}
