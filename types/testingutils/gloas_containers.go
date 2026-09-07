package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// Gloas (ePBS, SIP #94) beacon-container fixtures for the encoding spec tests. All are SSZ containers
// whose hash tree roots are stable under EIP-7688 (unlike the §4 block / §6 envelope family) —
// fixed-size, or with only a variable-length byte-list member (BuilderRequestAuth.Data) that
// merkleizes independently of container layout.

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
