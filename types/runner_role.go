package types

// https://github.com/ssvlabs/ssv-spec/issues/423 - this can be taken down to one byte on the wire
type RunnerRole int32

const (
	RoleCommittee RunnerRole = iota // Combines attestation and sync committee duties
	RoleProposer
	RoleAggregatorCommittee // Combines aggregator and sync committee contribution duties
	RoleValidatorRegistration
	RoleVoluntaryExit
	RoleUnknown = -1
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
