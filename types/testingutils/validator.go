package testingutils

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var BaseValidator = func(keySet *TestKeySet) *ssv.Validator {
	return ssv.NewValidator(
		NewTestingNetwork(1, keySet.OperatorKeys[1]),
		NewTestingBeaconNode(),
		TestingCommitteeMember(keySet),
		TestingShare(keySet, TestingValidatorIndex),
		NewTestingKeyManager(),
		NewTestingOperatorSigner(keySet, 1),
		map[types.RunnerRole]ssv.Runner{
			types.RoleCommittee:                 CommitteeRunner(keySet),
			types.RoleProposer:                  ProposerRunner(keySet),
			types.RoleAggregator:                AggregatorRunner(keySet),
			types.RoleSyncCommitteeContribution: SyncCommitteeContributionRunner(keySet),
			types.RoleValidatorRegistration:     ValidatorRegistrationRunner(keySet),
			types.RoleVoluntaryExit:             VoluntaryExitRunner(keySet),
		},
		NewTestingVerifier(),
	)
}
