package testingutils

import (
	"github.com/ssvlabs/ssv-spec/ssv"
)

var BaseValidatorCommitBoost = func(keySet *TestKeySet) *ssv.ValidatorCommitBoost {
	return ssv.NewValidatorCommitBoost(
		NewTestingBeaconNode().GetBeaconNetwork(),
		NewTestingNetwork(1, keySet.OperatorKeys[1]),
		NewTestingBeaconNode(),
		TestingCommitteeMember(keySet),
		TestingShare(keySet, TestingValidatorIndex),
		NewTestingKeyManager(),
		NewOperatorSigner(keySet, 1),
	)
}
