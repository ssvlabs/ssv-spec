package types

import "github.com/attestantio/go-eth2-client/spec/phase0"

// PreconfRequest is used as the data to be agreed on consensus for the PreconfRunner
// https://github.com/Commit-Boost/commit-boost-client/blob/main/crates/common/src/commit/request.rs#L82
type PreconfRequest struct {
	Root phase0.Root `ssz-size:"32"`
}

// Encode the PreconfRequest object
func (p *PreconfRequest) Encode() ([]byte, error) {
	return p.MarshalSSZ()
}

// Decode the PreconfRequest object
func (p *PreconfRequest) Decode(data []byte) error {
	return p.UnmarshalSSZ(data)
}

type SignedPreconfRequest struct {
	Root      phase0.Root         `ssz-size:"32"`
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// Encode the SignedPreconfRequest object
func (p *SignedPreconfRequest) Encode() ([]byte, error) {
	return p.MarshalSSZ()
}

// Decode the SignedPreconfRequest object
func (p *SignedPreconfRequest) Decode(data []byte) error {
	return p.UnmarshalSSZ(data)
}
