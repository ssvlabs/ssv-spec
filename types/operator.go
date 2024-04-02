package types

// OperatorID is a unique ID for the node, used to create shares and verify msgs
type OperatorID = uint64

// Operator represents an SSV operator node
type Operator struct {
	OperatorID    OperatorID
	SharePubKey   []byte `ssz-size:"48"`
	NetworkPubKey []byte `ssz-size:"294"`
}

// GetSharePublicKey returns the share public key
func (n *Operator) GetSharePublicKey() []byte {
	return n.SharePubKey
}

// GetNetworkPublicKey returns the network public key with which the node is identified with
func (n *Operator) GetNetworkPublicKey() []byte {
	return n.NetworkPubKey
}

// GetID returns the node's ID
func (n *Operator) GetID() OperatorID {
	return n.OperatorID
}
