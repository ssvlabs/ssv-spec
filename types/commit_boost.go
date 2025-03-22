package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// CBSigningRequest is used as the data to be agreed on consensus for the CBSigningRunner
// https://github.com/Commit-Boost/commit-boost-client/blob/main/crates/common/src/commit/request.rs#L82
type CBSigningRequest struct {
	Root phase0.Root `ssz-size:"32"`
}

// Encode the CBSigningRequest object
func (r *CBSigningRequest) Encode() ([]byte, error) {
	return r.MarshalSSZ()
}

// Decode the CBSigningRequest object
func (r *CBSigningRequest) Decode(data []byte) error {
	return r.UnmarshalSSZ(data)
}

type CBPartialSignature struct {
	RequestRoot phase0.Root `ssz-size:"32"`
	PartialSig  PartialSignatureMessages
}

// Encode the CBPartialSignature object
func (p *CBPartialSignature) Encode() ([]byte, error) {
	return p.MarshalSSZ()
}

// Decode the CBPartialSignature object
func (p *CBPartialSignature) Decode(data []byte) error {
	return p.UnmarshalSSZ(data)
}

type CBSigningDuty struct {
	Request CBSigningRequest
	Slot    phase0.Slot `ssz-size:"8"`
}

func (d *CBSigningDuty) DutySlot() phase0.Slot {
	return d.Slot
}

func (d *CBSigningDuty) RunnerRole() RunnerRole {
	return RoleCBSigning
}
