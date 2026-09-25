package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// Proposer Preferences (SIP #94 §5)
// ==================================================

// TestingProposerDutiesDependentRoot is the dependent root the testing beacon node reports for every
// epoch's proposer duties.
var TestingProposerDutiesDependentRoot = phase0.Root{0xde, 0x9e, 0x4d, 0xe1}

var TestingProposerPreferences = func(proposalSlot phase0.Slot) *gloas.ProposerPreferences {
	return &gloas.ProposerPreferences{
		DependentRoot:  TestingProposerDutiesDependentRoot,
		ProposalSlot:   proposalSlot,
		ValidatorIndex: TestingValidatorIndex,
		FeeRecipient:   TestingFeeRecipient,
		TargetGasLimit: types.DefaultGasLimit,
	}
}

// TestingSignedProposerPreferences is the reconstructed object the runner submits on quorum.
var TestingSignedProposerPreferences = func(ks *TestKeySet, proposalSlot phase0.Slot) *gloas.SignedProposerPreferences {
	preferences := TestingProposerPreferences(proposalSlot)
	return &gloas.SignedProposerPreferences{
		Message:   preferences,
		Signature: signBeaconObject(preferences, types.DomainProposerPreferences, ks),
	}
}

var TestingProposerPreferencesDuty = func() *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:           types.BNRoleProposerPreferences,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlotGloas,
		ValidatorIndex: TestingValidatorIndex,
	}
}

// TestingProposerPreferencesSecondDuty is a second lookahead proposal slot, concurrently active with
// TestingProposerPreferencesDuty.
var TestingProposerPreferencesSecondDuty = func() *types.ValidatorDuty {
	duty := TestingProposerPreferencesDuty()
	duty.Slot = TestingDutySlotGloas + 1
	return duty
}

var TestingProposerPreferencesNextEpochDuty = func() *types.ValidatorDuty {
	duty := TestingProposerPreferencesDuty()
	duty.Slot = TestingDutySlotGloasNextEpoch
	return duty
}

// TestingBuilderRequestAuth is the frozen §5 builder-request-auth for one entry's data at a proposal slot.
var TestingBuilderRequestAuth = func(data []byte, proposalSlot phase0.Slot) *gloas.BuilderRequestAuth {
	return &gloas.BuilderRequestAuth{
		Data: data,
		Slot: proposalSlot,
	}
}

// TestingSignedBuilderRequestAuth is the reconstructed object the runner submits on per-root quorum.
var TestingSignedBuilderRequestAuth = func(ks *TestKeySet, data []byte, proposalSlot phase0.Slot) *gloas.SignedBuilderRequestAuth {
	auth := TestingBuilderRequestAuth(data, proposalSlot)
	return &gloas.SignedBuilderRequestAuth{
		Message:   auth,
		Signature: signBeaconObject(auth, types.DomainBuilderRequestAuth, ks),
	}
}

// ProposerDutiesDependentRoot returns the fixture dependent root for any epoch
func (bn *TestingBeaconNode) ProposerDutiesDependentRoot(epoch phase0.Epoch) (phase0.Root, error) {
	return TestingProposerDutiesDependentRoot, nil
}

// SubmitProposerPreferences records the signed proposer preferences' root
func (bn *TestingBeaconNode) SubmitProposerPreferences(preferences *gloas.SignedProposerPreferences) error {
	r, err := preferences.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}

// SubmitBuilderRequestAuth records the signed builder-request-auth's root (SIP #94 §5)
func (bn *TestingBeaconNode) SubmitBuilderRequestAuth(auth *gloas.SignedBuilderRequestAuth) error {
	r, err := auth.HashTreeRoot()
	if err != nil {
		return err
	}
	bn.BroadcastedRoots = append(bn.BroadcastedRoots, r)
	return nil
}
