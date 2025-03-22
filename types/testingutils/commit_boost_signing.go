package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingCBSigningRequest = &types.CBSigningRequest{
	Root: phase0.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2},
}

var TestingCBSigningRequestWrong = &types.CBSigningRequest{
	Root: phase0.Root{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 10, 9},
}

// ==================================================
// Signed CBSigningRequest Object
// ==================================================

var TestingSignedCBSigningRequest = func(ks *TestKeySet) phase0.BLSSignature {
	return signBeaconObject(TestingCBSigningRequest, types.DomainCommitBoost, ks)
}

// ==================================================
// Commit-Boost Signing Duty
// ==================================================

var TestingCBSigningDuty = types.CBSigningDuty{
	Request: *TestingCBSigningRequest,
	Duty: types.ValidatorDuty{
		Type:           types.BNRoleCBSigning,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlot,
		ValidatorIndex: TestingValidatorIndex,
	},
}

var TestingCBSigningDutyNextEpoch = types.CBSigningDuty{
	Request: *TestingCBSigningRequest,
	Duty: types.ValidatorDuty{
		Type:           types.BNRoleCBSigning,
		PubKey:         TestingValidatorPubKey,
		Slot:           TestingDutySlot2,
		ValidatorIndex: TestingValidatorIndex,
	},
}
