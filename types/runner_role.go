package types

// https://github.com/ssvlabs/ssv-spec/issues/423 - this can be taken down to one byte on the wire
type RunnerRole int32

const (
	RoleCommittee             = RunnerRole(0) // Combines attestation and sync committee duties
	RoleProposer              = RunnerRole(2)
	RoleValidatorRegistration = RunnerRole(4)
	RoleVoluntaryExit         = RunnerRole(5)
	RoleAggregatorCommittee   = RunnerRole(6) // Combines aggregator and sync committee contribution duties
	RoleUnknown               = RunnerRole(-1)
)

// String returns the name of the runner role
func (r RunnerRole) String() string {
	switch r {
	case RoleCommittee:
		return "COMMITTEE_RUNNER"
	case RoleProposer:
		return "PROPOSER_RUNNER"
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
