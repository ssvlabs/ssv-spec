package testingutils

import (
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingOperatorID = types.OperatorID(1)

var BaseValidator = func(keySet *TestKeySet) *ssv.Validator {
	return ssv.NewValidator(
		NewTestingNetwork(TestingOperatorID, keySet.OperatorKeys[TestingOperatorID]),
		NewTestingBeaconNode(),
		TestingCommitteeMember(keySet),
		TestingShare(keySet, TestingValidatorIndex),
		NewTestingKeyManager(),
		NewOperatorSigner(keySet, TestingOperatorID),
		map[types.RunnerRole]ssv.Runner{
			types.RoleCommittee:                 CommitteeRunner(keySet),
			types.RoleProposer:                  ProposerRunner(keySet),
			types.RoleAggregator:                AggregatorRunner(keySet),
			types.RoleSyncCommitteeContribution: SyncCommitteeContributionRunner(keySet),
			types.RoleAggregatorCommittee:       AggregatorCommitteeRunner(keySet),
			types.RoleValidatorRegistration:     ValidatorRegistrationRunner(keySet),
			types.RoleVoluntaryExit:             VoluntaryExitRunner(keySet),
		},
	)
}
