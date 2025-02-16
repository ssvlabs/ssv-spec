package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// Voluntary Exit Object
// ==================================================

var TestingVoluntaryExit = &phase0.VoluntaryExit{
	Epoch:          0,
	ValidatorIndex: TestingValidatorIndex,
}
var TestingVoluntaryExitWrong = &phase0.VoluntaryExit{
	Epoch:          1,
	ValidatorIndex: TestingValidatorIndex,
}

func TestingVoluntaryExitBySlot(slot phase0.Slot) *phase0.VoluntaryExit {
	epoch := types.PraterNetwork.EstimatedEpochAtSlot(slot)
	return &phase0.VoluntaryExit{
		Epoch:          epoch,
		ValidatorIndex: TestingValidatorIndex,
	}
}

// ==================================================
// Signed Voluntary Exit Object
// ==================================================

var TestingSignedVoluntaryExit = func(ks *TestKeySet) *phase0.SignedVoluntaryExit {
	return &phase0.SignedVoluntaryExit{
		Message:   TestingVoluntaryExit,
		Signature: signBeaconObject(TestingVoluntaryExit, types.DomainVoluntaryExit, ks),
	}
}

// ==================================================
// Voluntary Exit Duty
// ==================================================

var TestingVoluntaryExitDuty = types.ValidatorDuty{
	Type:           types.BNRoleVoluntaryExit,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot,
	ValidatorIndex: TestingValidatorIndex,
}

var TestingVoluntaryExitDutyNextEpoch = types.ValidatorDuty{
	Type:           types.BNRoleVoluntaryExit,
	PubKey:         TestingValidatorPubKey,
	Slot:           TestingDutySlot2,
	ValidatorIndex: TestingValidatorIndex,
}
