package testingutils

import (
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/ssv"
	"github.com/MatheusFranco99/ssv-spec-AleaBFT/types"
)

var BaseValidator = func(keySet *TestKeySet) *ssv.Validator {
	return ssv.NewValidator(
		NewTestingNetwork(),
		NewTestingBeaconNode(),
		TestingShare(keySet),
		NewTestingKeyManager(),
		map[types.BeaconRole]ssv.Runner{
			types.BNRoleAttester:                  AttesterRunner(keySet),
			types.BNRoleProposer:                  ProposerRunner(keySet),
			types.BNRoleAggregator:                AggregatorRunner(keySet),
			types.BNRoleSyncCommittee:             SyncCommitteeRunner(keySet),
			types.BNRoleSyncCommitteeContribution: SyncCommitteeContributionRunner(keySet),
		},
	)
}
