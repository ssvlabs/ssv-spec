package testingutils

import (
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Validator Registration Object
// ==================================================

var TestingFeeRecipient = bellatrix.ExecutionAddress(ethAddressFromHex("535953b5a6040074948cf185eaa7d2abbd66808f"))

var TestingValidatorRegistration = &v1.ValidatorRegistration{
	FeeRecipient: TestingFeeRecipient,
	GasLimit:     types.DefaultGasLimit,
	Timestamp:    types.PraterNetwork.EpochStartTime(TestingDutyEpoch),
	Pubkey:       TestingValidatorPubKey,
}

var TestingValidatorRegistrationWrong = &v1.ValidatorRegistration{
	FeeRecipient: TestingFeeRecipient,
	GasLimit:     5,
	Timestamp:    types.PraterNetwork.EpochStartTime(TestingDutyEpoch),
	Pubkey:       TestingValidatorPubKey,
}

func TestingValidatorRegistrationBySlot(slot phase0.Slot) *v1.ValidatorRegistration {
	epoch := types.PraterNetwork.EstimatedEpochAtSlot(slot)
	return &v1.ValidatorRegistration{
		FeeRecipient: TestingFeeRecipient,
		GasLimit:     types.DefaultGasLimit,
		Timestamp:    types.PraterNetwork.EpochStartTime(epoch),
		Pubkey:       TestingValidatorPubKey,
	}
}

// ==================================================
// Signed Validator Registration Object
// ==================================================

var TestingSignedValidatorRegistration = func(ks *TestKeySet) *v1.SignedValidatorRegistration {
	vr := TestingValidatorRegistration
	sig := signBeaconObject(vr, types.DomainApplicationBuilder, ks)
	return &v1.SignedValidatorRegistration{
		Message:   vr,
		Signature: sig,
	}
}

var TestingSignedValidatorRegistrationWrong = func(ks *TestKeySet) *v1.SignedValidatorRegistration {
	vr := TestingValidatorRegistrationWrong
	sig := signBeaconObject(vr, types.DomainApplicationBuilder, ks)
	return &v1.SignedValidatorRegistration{
		Message:   vr,
		Signature: sig,
	}
}

var TestingSignedValidatorRegistrationBySlot = func(ks *TestKeySet, slot phase0.Slot) *v1.SignedValidatorRegistration {
	vr := TestingValidatorRegistrationBySlot(slot)
	sig := signBeaconObject(vr, types.DomainApplicationBuilder, ks)
	return &v1.SignedValidatorRegistration{
		Message:   vr,
		Signature: sig,
	}
}

// ==================================================
// Validator Registration Duty
// ==================================================

var TestingValidatorRegistrationDuty = types.ValidatorDuty{
	Type:           types.BNRoleValidatorRegistration,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingValidatorRegistrationDutyNextEpoch = types.ValidatorDuty{
	Type:           types.BNRoleValidatorRegistration,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot2,
	ValidatorIndex: TestingValidatorIndex,
}
