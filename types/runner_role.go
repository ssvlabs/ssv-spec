package types

// https://github.com/ssvlabs/ssv-spec/issues/423 - this can be taken down to one byte on the wire
type RunnerRole int32

const (
	RoleCommittee  RunnerRole = iota // Combines attestation and sync committee duties
	RoleAggregator                   // RoleAggregator will be unused after the Boole fork
	RoleProposer
	RoleSyncCommitteeContribution // RoleSyncCommitteeContribution will be unused after the Boole fork
	RoleValidatorRegistration
	RoleVoluntaryExit
	RoleAggregatorCommittee // Combines aggregator and sync committee contribution duties
	RoleUnknown             = -1
)

// String returns the name of the runner role
func (r RunnerRole) String() string {
	switch r {
	case RoleCommittee:
		return "COMMITTEE_RUNNER"
	case RoleAggregator:
		return "AGGREGATOR_RUNNER"
	case RoleProposer:
		return "PROPOSER_RUNNER"
	case RoleSyncCommitteeContribution:
		return "SYNC_COMMITTEE_CONTRIBUTION_RUNNER"
	case RoleValidatorRegistration:
		return "VALIDATOR_REGISTRATION_RUNNER"
	case RoleVoluntaryExit:
		return "VOLUNTARY_EXIT_RUNNER"
	case RoleAggregatorCommittee:
		return "AGGREGATOR_COMMITTEE_RUNNER"
	default:
		return "UNDEFINED"
	}
}
