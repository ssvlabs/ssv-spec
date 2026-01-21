package maxmsgsize

import (
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	MaxSizePhase0Attestation  = 2276
	MaxSizeElectraAttestation = 131308
)

// Wrapper types to allow defining Encode/Decode methods on external types
type Phase0AttestationWrapper struct {
	*phase0.Attestation
}

func (w *Phase0AttestationWrapper) Encode() ([]byte, error) {
	return w.MarshalSSZ()
}

func (w *Phase0AttestationWrapper) Decode(data []byte) error {
	return w.UnmarshalSSZ(data)
}

type ElectraAttestationWrapper struct {
	*electra.Attestation
}

func (w *ElectraAttestationWrapper) Encode() ([]byte, error) {
	return w.MarshalSSZ()
}

func (w *ElectraAttestationWrapper) Decode(data []byte) error {
	return w.UnmarshalSSZ(data)
}

func maxPhase0Attestation() *Phase0AttestationWrapper {
	aggbits := [2048]byte{1}
	return &Phase0AttestationWrapper{
		Attestation: &phase0.Attestation{
			AggregationBits: bitfield.Bitlist(aggbits[:]),
			Data: &phase0.AttestationData{
				Slot:            1,
				Index:           0,
				BeaconBlockRoot: [32]byte{1},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
			},
			Signature: phase0.BLSSignature([96]byte{1}),
		},
	}
}

func maxElectraAttestation() *ElectraAttestationWrapper {
	aggbits := [131072]byte{1}
	return &ElectraAttestationWrapper{
		Attestation: &electra.Attestation{
			AggregationBits: bitfield.Bitlist(aggbits[:]),
			Data: &phase0.AttestationData{
				Slot:            1,
				Index:           0,
				BeaconBlockRoot: [32]byte{1},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
				Target: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{1},
				},
			},
			Signature:     phase0.BLSSignature([96]byte{1}),
			CommitteeBits: bitfield.NewBitvector64(),
		},
	}
}

func MaxPhase0Attestation() *StructureSizeTest {
	return NewStructureSizeTest(
		"max Phase0Attestation",
		testdoc.StructureSizeTestMaxPhase0AttestationDoc,
		maxPhase0Attestation(),
		MaxSizePhase0Attestation,
		true,
	)
}

func MaxElectraAttestation() *StructureSizeTest {
	return NewStructureSizeTest(
		"max ElectraAttestation",
		testdoc.StructureSizeTestMaxElectraAttestationDoc,
		maxElectraAttestation(),
		MaxSizeElectraAttestation,
		true,
	)
}
