package testingutils

import "github.com/ssvlabs/ssv-spec/types"

var TestingPreconfRequest = &types.PreconfRequest{
	Root: TestingPreconfRoot,
}

// ==================================================
// Preconfirmation Duty
// ==================================================

var TestingPreconfirmationDuty = types.ValidatorDuty{
	Type:           types.BNRolePreconfirmation,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot,
	ValidatorIndex: TestingValidatorIndex,
}
