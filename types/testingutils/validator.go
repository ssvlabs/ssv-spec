package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var BaseValidator = func(keySet *TestKeySet) *ssv.Validator {
	return BaseValidatorWithIndex(keySet, TestingValidatorIndex)
}

var BaseValidatorWithIndex = func(keySet *TestKeySet, validatorIndex phase0.ValidatorIndex) *ssv.Validator {
	return ssv.NewValidator(
		TestingCommitteeMember(keySet),
		TestingShare(keySet, validatorIndex),
		map[types.RunnerRole]ssv.Runner{
			types.RoleCommittee:                 CommitteeRunner(keySet),
			types.RoleProposer:                  ProposerRunner(keySet),
			types.RoleAggregator:                AggregatorRunner(keySet),
			types.RoleSyncCommitteeContribution: SyncCommitteeContributionRunner(keySet),
			types.RoleValidatorRegistration:     ValidatorRegistrationRunner(keySet),
			types.RoleVoluntaryExit:             VoluntaryExitRunner(keySet),
		},
	)
}
