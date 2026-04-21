package types

import "math"

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

const runnerRoleUnknownWireUint32 = uint32(math.MaxUint32)

// WireUint32 returns the canonical on-wire representation for RunnerRole.
//
// Note: RunnerRole is encoded as a 4-byte little-endian uint32 in MessageID.
// RoleUnknown uses the all-ones sentinel (MaxUint32), consistent with other enums
// (e.g. BeaconRole uses MaxUint64).
func (r RunnerRole) WireUint32() uint32 {
	if r == RoleUnknown || r < 0 {
		return runnerRoleUnknownWireUint32
	}
	return uint32(r)
}

// RunnerRoleFromWireUint32 converts a 4-byte uint32 MessageID role field into RunnerRole.
// Values outside the int32 range are treated as unknown to avoid surprising negative roles.
func RunnerRoleFromWireUint32(v uint32) RunnerRole {
	// Any value outside the int32 range (including the RoleUnknown sentinel) is treated as unknown.
	if v > uint32(math.MaxInt32) {
		return RoleUnknown
	}
	return RunnerRole(int32(v))
}

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
