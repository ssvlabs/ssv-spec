package types

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID = uint64

// Operator represents an SSV operator node
type Operator struct {
	OperatorID  OperatorID
	SharePubKey []byte `ssz-size:"48"`
	SSVPubKey   []byte `ssz-size:"294"`
}

// GetSharePublicKey returns the share public key
func (n *Operator) GetSharePublicKey() []byte {
	return n.SharePubKey
}

// GetSSVPublicKey returns the ssv public key with which the node is identified with
func (n *Operator) GetSSVPublicKey() []byte {
	return n.SSVPubKey
}

// GetID returns the node's ID
func (n *Operator) GetID() OperatorID {
	return n.OperatorID
}
