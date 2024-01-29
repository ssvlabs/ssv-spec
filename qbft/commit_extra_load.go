package qbft

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
)

type PostConsensusSignature struct {
	Signature   types.Signature `ssz-size:"96"`
	SigningRoot phase0.Root     `ssz-size:"32"`
}

// Encode encoded bytes or error
func (p *PostConsensusSignature) Encode() ([]byte, error) {
	return p.MarshalSSZ()
}

// Decode returns error if decoding failed
func (p *PostConsensusSignature) Decode(data []byte) error {
	return p.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (p *PostConsensusSignature) GetRoot() ([32]byte, error) {
	return p.HashTreeRoot()
}

// This structure is appended in commit messages carrying extra information.
// It's useful for the underlying duty but not necessary for the consensus termination
type CommitExtraLoad struct {
	PostConsensusSignatures []*PostConsensusSignature `ssz-max:"4"` // Maximum number of sync committee subnets as defined in https://github.com/ethereum/consensus-specs/blob/dev/specs/altair/validator.md
}

func (c *CommitExtraLoad) GetPostConsensusSignatures() []*PostConsensusSignature {
	return c.PostConsensusSignatures
}

// Encode returns a CommitExtraLoad encoded bytes or error
func (c *CommitExtraLoad) Encode() ([]byte, error) {
	return c.MarshalSSZ()
}

// Decode returns error if decoding failed
func (c *CommitExtraLoad) Decode(data []byte) error {
	return c.UnmarshalSSZ(data)
}

// GetRoot returns the root used for signing and verification
func (c *CommitExtraLoad) GetRoot() ([32]byte, error) {
	return c.HashTreeRoot()
}

// Upon commit messages, the Instance shall use its CommitExtraLoadManagerI to:
// - Validate() the CommitExtraLoad
// - Process()
// To create the CommitExtraLoad object, the Instance shall call Create()
type CommitExtraLoadManagerI interface {
	Validate(signedMessage *SignedMessage, fullData []byte) error
	Process(signedMessage *SignedMessage) error
	Create(fullData []byte) (CommitExtraLoad, error)
}
