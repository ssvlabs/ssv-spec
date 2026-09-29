package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// Versioned Beacon Block
// ==================================================

// SupportedBlockVersions is a list of supported beacon block versions by spec.
var SupportedBlockVersions = []spec.DataVersion{gloas.DataVersionGloas}

// TestingBeaconBlockBytesV is the decided value's DataSSZ: the {block, payload_root} proposal data (SIP #94 §4).
var TestingBeaconBlockBytesV = func(version spec.DataVersion) []byte {
	if version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	return TestingGloasProposalDataBytes(TestingDutySlotV(version))
}

// TestingSignedBeaconBlockV is the bid-only block signed under DomainProposer (SIP #94 §4).
var TestingSignedBeaconBlockV = func(ks *TestKeySet, version spec.DataVersion) types.HashRoot {
	if version != gloas.DataVersionGloas {
		panic("unsupported version")
	}
	blk := gloas.TestingBeaconBlock(TestingDutySlotV(version))
	return &gloas.SignedBeaconBlock{
		Message:   blk,
		Signature: signBeaconObject(blk, types.DomainProposer, ks),
	}
}

// ==================================================
// Versioned Proposer Duty
// ==================================================

var TestingProposerDutyV = func(version spec.DataVersion) *types.ValidatorDuty {
	duty := &types.ValidatorDuty{
		Type:           types.BNRoleProposer,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlotV(version),
		ValidatorIndex: TestingValidatorIndex,
		// ISSUE 233: We are initializing unused struct fields here
		CommitteeIndex:          3,
		CommitteesAtSlot:        36,
		CommitteeLength:         128,
		ValidatorCommitteeIndex: 11,
	}

	return duty
}

var TestingProposerDutyNextEpochV = func(version spec.DataVersion) *types.ValidatorDuty {
	duty := TestingProposerDutyV(version)
	duty.Slot = TestingDutySlotNextEpochV(version)
	return duty
}

var TestingProposerDutyFirstSlot = types.ValidatorDuty{
	Type:           types.BNRoleProposer,
	PubKey:         TestingValidatorPubKey,
	Slot:           0,
	ValidatorIndex: TestingValidatorIndex,
}
