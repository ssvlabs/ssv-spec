package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	bitfield "github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// Gloas (ePBS, SIP #94) beacon-container fixtures for the encoding spec tests.

var TestPayloadAttestationData = &gloas.PayloadAttestationData{
	BeaconBlockRoot:   phase0.Root{0x01, 0x02},
	Slot:              42,
	PayloadPresent:    true,
	BlobDataAvailable: false,
}

var TestPayloadAttestationMessage = &gloas.PayloadAttestationMessage{
	ValidatorIndex: 7,
	Data:           TestPayloadAttestationData,
	Signature:      phase0.BLSSignature{0xbb, 0xcc},
}

var TestProposerPreferences = &gloas.ProposerPreferences{
	DependentRoot:  phase0.Root{0x01, 0x02},
	ProposalSlot:   42,
	ValidatorIndex: 7,
	FeeRecipient:   bellatrix.ExecutionAddress{0xaa, 0xbb},
	TargetGasLimit: 36_000_000,
}

var TestSignedProposerPreferences = &gloas.SignedProposerPreferences{
	Message:   TestProposerPreferences,
	Signature: phase0.BLSSignature{0xbb, 0xcc},
}

var TestBuilderRequestAuth = &gloas.BuilderRequestAuth{
	Data: []byte{0x01, 0x02, 0x03},
	Slot: 42,
}

var TestSignedBuilderRequestAuth = &gloas.SignedBuilderRequestAuth{
	Message:   TestBuilderRequestAuth,
	Signature: phase0.BLSSignature{0xbb, 0xcc},
}

// TestSignedBeaconBlock populates the ePBS-specific body fields: the payload bid, an aggregated payload
// attestation, and the parent execution requests (SIP #94 §4).
var TestSignedBeaconBlock = &gloas.SignedBeaconBlock{
	Message: &gloas.BeaconBlock{
		Slot:          7,
		ProposerIndex: 3,
		Body: &gloas.BeaconBlockBody{
			ETH1Data:      &phase0.ETH1Data{BlockHash: make([]byte, 32)},
			SyncAggregate: &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{
				BuilderIndex:       gloas.BuilderIndexSelfBuild,
				BlobKZGCommitments: []deneb.KZGCommitment{{0x01}},
			}},
			PayloadAttestations: []*gloas.PayloadAttestation{{
				AggregationBits: bitfield.NewBitvector512(),
				Data:            &gloas.PayloadAttestationData{Slot: 6, PayloadPresent: true},
			}},
			ParentExecutionRequests: &gloas.ExecutionRequests{},
		},
	},
	Signature: phase0.BLSSignature{0xbb, 0xcc},
}

var TestSignedExecutionPayloadBid = &gloas.SignedExecutionPayloadBid{
	Message: &gloas.ExecutionPayloadBid{
		BlockHash:          [32]byte{0xaa},
		BuilderIndex:       gloas.BuilderIndexSelfBuild,
		Value:              123,
		BlobKZGCommitments: []deneb.KZGCommitment{{0x01}, {0x02}},
	},
	Signature: phase0.BLSSignature{0xbb, 0xcc},
}

var TestBlindedExecutionPayloadEnvelope = &gloas.BlindedExecutionPayloadEnvelope{
	PayloadRoot:           phase0.Root{0x01},
	ExecutionRequestsRoot: phase0.Root{0x04},
	BuilderIndex:          gloas.BuilderIndexSelfBuild,
	BeaconBlockRoot:       phase0.Root{0x02},
	ParentBeaconBlockRoot: phase0.Root{0x03},
}
