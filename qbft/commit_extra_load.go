package qbft

import "github.com/bloxapp/ssv-spec/types"

// This structure is appended in commit messages carrying extra information.
// It's useful for the underlying duty but not necessary for the consensus termination
type CommitExtraLoad struct {
	// List of signatures. We must use list instead of a unique element due to the sync committee aggregation duty.
	Signatures []types.Signature `ssz-max:"4"` // Maximum number of sync committee subnets as defined in https://github.com/ethereum/consensus-specs/blob/dev/specs/altair/validator.md
}

func (c *CommitExtraLoad) GetSignatures() []types.Signature {
	return c.Signatures
}

// Encode returns a CommitExtraLoad encoded bytes or error
func (c *CommitExtraLoad) Encode() ([]byte, error) {
	return c.MarshalSSZ()
}

// Decode returns error if decoding failed
func (c *CommitExtraLoad) Decode(data []byte) error {
	return c.UnmarshalSSZ(data)
}

// Upon commit messages, the Instance shall use its CommitExtraLoadManagerI to:
// - Validate() the CommitExtraLoad
// - Process()
// To create the CommitExtraLoad object, the Instance shall call Create()
type CommitExtraLoadManagerI interface {
	Validate(signedMessage *SignedMessage) error
	Process(signedMessage *SignedMessage) error
	Create(fullData []byte) (CommitExtraLoad, error)
}
