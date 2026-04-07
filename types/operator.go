package types

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID = uint64

// CommitteeMember represents an SSV operator node that is part of a committee
type CommitteeMember struct {
	OperatorID        OperatorID
	CommitteeID       CommitteeID `ssz-size:"32"`
	SSVOperatorPubKey []byte      `ssz-size:"459"`
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

// Validate checks the following rules:
// - OperatorID must be non-zero
// - SSVOperatorPubKey must be 459 bytes
func (o *Operator) Validate() error {
	if o == nil {
		return NewError(InvalidOperatorErrorCode, "nil operator")
	}
	if o.OperatorID == 0 {
		return NewError(InvalidOperatorErrorCode, "operator ID 0 not allowed")
	}
	if len(o.SSVOperatorPubKey) != 459 {
		return NewError(InvalidOperatorErrorCode, "invalid operator public key length")
	}
	return nil
}

// Validate checks the following rules:
// - OperatorID must be non-zero
// - SSVOperatorPubKey must be 459 bytes
// - Committee must be non-empty and within the ssz-max bound
// - FaultyNodes must satisfy the QBFT committee requirement n >= 3f+1
// - CommitteeID must match the ID computed from committee OperatorIDs
// - DomainType must be one of the known SSV domains in this spec
func (cm *CommitteeMember) Validate() error {
	if cm == nil {
		return NewError(InvalidCommitteeMemberErrorCode, "nil committee member")
	}
	if cm.OperatorID == 0 {
		return NewError(InvalidCommitteeMemberErrorCode, "operator ID 0 not allowed")
	}
	if len(cm.SSVOperatorPubKey) != 459 {
		return NewError(InvalidCommitteeMemberErrorCode, "invalid operator public key length")
	}

	if len(cm.Committee) == 0 {
		return NewError(InvalidCommitteeMemberErrorCode, "empty committee")
	}
	if len(cm.Committee) > MaxCommitteeSize {
		return NewError(InvalidCommitteeMemberErrorCode, "committee too large")
	}

	committeeSize := uint64(len(cm.Committee))
	// Enforce the QBFT committee requirement n == 3f+1 without risking uint64 overflow.
	// This also guarantees that quorum calculations using 2f+1 have the expected intersection properties.
	if (committeeSize-1)%3 != 0 || cm.FaultyNodes != (committeeSize-1)/3 {
		return NewError(InvalidCommitteeMemberErrorCode, "invalid faulty nodes bound for committee size")
	}

	seenIDs := make(map[OperatorID]struct{}, len(cm.Committee))
	opIDs := make([]OperatorID, 0, len(cm.Committee))
	containsSelf := false
	for _, op := range cm.Committee {
		if err := op.Validate(); err != nil {
			return NewError(InvalidCommitteeMemberErrorCode, "invalid committee operator")
		}
		if _, exists := seenIDs[op.OperatorID]; exists {
			return NewError(InvalidCommitteeMemberErrorCode, "duplicate operator ID in committee")
		}
		seenIDs[op.OperatorID] = struct{}{}
		opIDs = append(opIDs, op.OperatorID)
		if op.OperatorID == cm.OperatorID {
			containsSelf = true
		}
	}
	if !containsSelf {
		return NewError(InvalidCommitteeMemberErrorCode, "committee does not contain operator")
	}

	if GetCommitteeID(opIDs) != cm.CommitteeID {
		return NewError(InvalidCommitteeMemberErrorCode, "committee ID mismatch")
	}

	if !cm.DomainType.IsKnown() {
		return NewError(InvalidCommitteeMemberErrorCode, "unknown domain type")
	}

	return nil
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
