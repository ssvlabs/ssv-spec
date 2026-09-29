package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWireConstants pins every consensus-critical wire constant to its exact on-the-wire
// value. Domains, and the iota-assigned BeaconRole / RunnerRole / PartialSigMsgType blocks,
// must match consensus-specs / builder-specs and every other SSV implementation (e.g. Anchor):
// they are SSZ-encoded or used as signing domains, so any drift — reordering an iota block,
// swapping two entries, editing a domain byte — must fail loudly here rather than silently
// diverge on the wire. Each type is pinned as a name→value table so existing values are held
// just as firmly as the Gloas (ePBS, SIP #94) additions.
func TestWireConstants(t *testing.T) {
	t.Run("domains", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			got  [4]byte
			want [4]byte
		}{
			{"DomainProposer", DomainProposer, [4]byte{0x00, 0x00, 0x00, 0x00}},
			{"DomainAttester", DomainAttester, [4]byte{0x01, 0x00, 0x00, 0x00}},
			{"DomainRandao", DomainRandao, [4]byte{0x02, 0x00, 0x00, 0x00}},
			{"DomainDeposit", DomainDeposit, [4]byte{0x03, 0x00, 0x00, 0x00}},
			{"DomainVoluntaryExit", DomainVoluntaryExit, [4]byte{0x04, 0x00, 0x00, 0x00}},
			{"DomainSelectionProof", DomainSelectionProof, [4]byte{0x05, 0x00, 0x00, 0x00}},
			{"DomainAggregateAndProof", DomainAggregateAndProof, [4]byte{0x06, 0x00, 0x00, 0x00}},
			{"DomainSyncCommittee", DomainSyncCommittee, [4]byte{0x07, 0x00, 0x00, 0x00}},
			{"DomainSyncCommitteeSelectionProof", DomainSyncCommitteeSelectionProof, [4]byte{0x08, 0x00, 0x00, 0x00}},
			{"DomainContributionAndProof", DomainContributionAndProof, [4]byte{0x09, 0x00, 0x00, 0x00}},
			{"DomainApplicationBuilder", DomainApplicationBuilder, [4]byte{0x00, 0x00, 0x00, 0x01}},
			// Gloas (ePBS) beacon domains, plus builder-specs' DOMAIN_BUILDER_REQUEST_AUTH which
			// shares the 0x0b first byte with DomainBeaconBuilder and differs only in the last byte.
			{"DomainBeaconBuilder", DomainBeaconBuilder, [4]byte{0x0b, 0x00, 0x00, 0x00}},
			{"DomainPTCAttester", DomainPTCAttester, [4]byte{0x0c, 0x00, 0x00, 0x00}},
			{"DomainProposerPreferences", DomainProposerPreferences, [4]byte{0x0d, 0x00, 0x00, 0x00}},
			{"DomainBuilderRequestAuth", DomainBuilderRequestAuth, [4]byte{0x0b, 0x00, 0x00, 0x01}},
		} {
			require.Equal(t, tc.want, tc.got, tc.name)
		}
	})

	t.Run("beacon_roles", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			got  BeaconRole
			want BeaconRole
		}{
			{"BNRoleAttester", BNRoleAttester, 0},
			{"BNRoleAggregator", BNRoleAggregator, 1},
			{"BNRoleProposer", BNRoleProposer, 2},
			{"BNRoleSyncCommittee", BNRoleSyncCommittee, 3},
			{"BNRoleSyncCommitteeContribution", BNRoleSyncCommitteeContribution, 4},
			{"BNRoleValidatorRegistration", BNRoleValidatorRegistration, 5},
			{"BNRoleVoluntaryExit", BNRoleVoluntaryExit, 6},
			{"BNRolePTCAttester", BNRolePTCAttester, 7},
			{"BNRoleProposerPreferences", BNRoleProposerPreferences, 8},
			{"BNRoleUnknown", BNRoleUnknown, math.MaxUint64},
		} {
			require.Equal(t, tc.want, tc.got, tc.name)
		}
	})

	t.Run("runner_roles", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			got  RunnerRole
			want RunnerRole
		}{
			{"RoleCommittee", RoleCommittee, 0},
			{"RoleProposer", RoleProposer, 2},
			{"RoleValidatorRegistration", RoleValidatorRegistration, 4},
			{"RoleVoluntaryExit", RoleVoluntaryExit, 5},
			{"RoleAggregatorCommittee", RoleAggregatorCommittee, 6},
			{"RolePTCAttester", RolePTCAttester, 7},
			{"RoleProposerPreferences", RoleProposerPreferences, 8},
			{"RoleUnknown", RoleUnknown, -1},
		} {
			require.Equal(t, tc.want, tc.got, tc.name)
		}
	})

	t.Run("partial_sig_types", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			got  PartialSigMsgType
			want PartialSigMsgType
		}{
			{"PostConsensusPartialSig", PostConsensusPartialSig, 0},
			{"RandaoPartialSig", RandaoPartialSig, 1},
			{"ValidatorRegistrationPartialSig", ValidatorRegistrationPartialSig, 4},
			{"VoluntaryExitPartialSig", VoluntaryExitPartialSig, 5},
			{"AggregatorCommitteePartialSig", AggregatorCommitteePartialSig, 6},
			{"PTCAttesterPartialSig", PTCAttesterPartialSig, 7},
			{"ProposerPreferencesPartialSig", ProposerPreferencesPartialSig, 8},
			{"RequestAuthPartialSig", RequestAuthPartialSig, 9},
		} {
			require.Equal(t, tc.want, tc.got, tc.name)
		}
	})

	// MapDutyToRunnerRole is consensus-critical routing: it decides which runner handles a duty.
	t.Run("duty_to_runner", func(t *testing.T) {
		for _, tc := range []struct {
			beacon BeaconRole
			runner RunnerRole
		}{
			{BNRoleAttester, RoleCommittee},
			{BNRoleSyncCommittee, RoleCommittee},
			{BNRoleProposer, RoleProposer},
			{BNRoleAggregator, RoleAggregatorCommittee},
			{BNRoleSyncCommitteeContribution, RoleAggregatorCommittee},
			{BNRoleValidatorRegistration, RoleValidatorRegistration},
			{BNRoleVoluntaryExit, RoleVoluntaryExit},
			{BNRolePTCAttester, RolePTCAttester},
			{BNRoleProposerPreferences, RoleProposerPreferences},
			{BNRoleUnknown, RoleUnknown},
		} {
			require.Equal(t, tc.runner, MapDutyToRunnerRole(tc.beacon), tc.beacon.String())
		}
	})

	// String() names surface in generated spec-test vectors, so pin them too.
	t.Run("role_strings", func(t *testing.T) {
		require.Equal(t, "ATTESTER", BNRoleAttester.String())
		require.Equal(t, "AGGREGATOR", BNRoleAggregator.String())
		require.Equal(t, "PROPOSER", BNRoleProposer.String())
		require.Equal(t, "SYNC_COMMITTEE", BNRoleSyncCommittee.String())
		require.Equal(t, "SYNC_COMMITTEE_CONTRIBUTION", BNRoleSyncCommitteeContribution.String())
		require.Equal(t, "VALIDATOR_REGISTRATION", BNRoleValidatorRegistration.String())
		require.Equal(t, "VOLUNTARY_EXIT", BNRoleVoluntaryExit.String())
		require.Equal(t, "PTC_ATTESTER", BNRolePTCAttester.String())
		require.Equal(t, "PROPOSER_PREFERENCES", BNRoleProposerPreferences.String())

		require.Equal(t, "COMMITTEE_RUNNER", RoleCommittee.String())
		require.Equal(t, "PROPOSER_RUNNER", RoleProposer.String())
		require.Equal(t, "VALIDATOR_REGISTRATION_RUNNER", RoleValidatorRegistration.String())
		require.Equal(t, "VOLUNTARY_EXIT_RUNNER", RoleVoluntaryExit.String())
		require.Equal(t, "AGGREGATOR_COMMITTEE_RUNNER", RoleAggregatorCommittee.String())
		require.Equal(t, "PTC_ATTESTER_RUNNER", RolePTCAttester.String())
		require.Equal(t, "PROPOSER_PREFERENCES_RUNNER", RoleProposerPreferences.String())
	})
}
