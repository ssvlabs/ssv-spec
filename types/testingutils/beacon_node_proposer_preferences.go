package testingutils

import (
	"encoding/binary"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// Proposer Preferences (SIP #94 §5)
// ==================================================

var TestingFeeRecipient = bellatrix.ExecutionAddress(ethAddressFromHex("535953b5a6040074948cf185eaa7d2abbd66808f"))

// TestingProposerDutiesDependentRoot is the dependent root the testing beacon node reports for the proposer
// duties of every epoch before TestingEpochDependentFrom.
var TestingProposerDutiesDependentRoot = phase0.Root{0xde, 0x9e, 0x4d, 0xe1}

// TestingProposerDutiesDependentRootAt is the dependent root the testing beacon node reports for an epoch's
// proposer duties: from TestingEpochDependentFrom on, each epoch has its own.
func TestingProposerDutiesDependentRootAt(epoch phase0.Epoch) phase0.Root {
	root := TestingProposerDutiesDependentRoot
	if epoch >= TestingEpochDependentFrom {
		binary.LittleEndian.PutUint64(root[8:], uint64(epoch))
	}
	return root
}

var TestingProposerPreferences = func(proposalSlot phase0.Slot) *gloas.ProposerPreferences {
	return &gloas.ProposerPreferences{
		DependentRoot:  TestingProposerDutiesDependentRootAt(types.BeaconTestNetwork.EstimatedEpochAtSlot(proposalSlot)),
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
		Signature: signBeaconObjectAtSlot(preferences, types.DomainProposerPreferences, ks, proposalSlot),
	}
}

var TestingProposerPreferencesDuty = func() *types.ValidatorDuty {
	return TestingProposerPreferencesDutyForSlot(TestingDutySlotGloas)
}

// TestingProposerPreferencesDutyForSlot is the proposer-preferences duty of a proposal slot.
var TestingProposerPreferencesDutyForSlot = func(proposalSlot phase0.Slot) *types.ValidatorDuty {
	return &types.ValidatorDuty{
		Type:           types.BNRoleProposerPreferences,
		PubKey:         TestingValidatorPubKey,
		Slot:           proposalSlot,
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

// ProposerDutiesDependentRoot returns the epoch's fixture dependent root (TestingProposerDutiesDependentRootAt),
// or fails under SetFailProposerDutiesDependentRoot.
func (bn *TestingBeaconNode) ProposerDutiesDependentRoot(epoch phase0.Epoch) (phase0.Root, error) {
	if bn.failProposerDutiesDependentRoot {
		return phase0.Root{}, fmt.Errorf("forced proposer duties dependent root failure")
	}
	return TestingProposerDutiesDependentRootAt(epoch), nil
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
