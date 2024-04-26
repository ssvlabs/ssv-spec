package types

// TODO since this is on wire no real need to take 32 bits
type RunnerRole int32

const (
	RoleCommittee RunnerRole = iota
	RoleAggregator
	RoleProposer
	RoleSyncCommitteeContribution

	RoleValidatorRegistration
	RoleVoluntaryExit

	RoleUnknown = -1
)
