package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingCBSigningRequest = &types.CBSigningRequest{
	Root: TestingPreconfRoot,
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
// Preconfirmation Duty
// ==================================================

var TestingPreconfDuty = types.ValidatorDuty{
	Type:           types.BNRoleCBSigning,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingPreconfDutyNextEpoch = types.ValidatorDuty{
	Type:           types.BNRoleCBSigning,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot2,
	ValidatorIndex: TestingValidatorIndex,
}
