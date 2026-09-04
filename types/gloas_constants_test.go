package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGloasWireConstants pins the ePBS (Gloas, SIP #94) wire constants to their exact
// on-the-wire values. These are consensus-critical: they must match consensus-specs /
// builder-specs and every other SSV implementation (e.g. Anchor), so any drift —
// renumbering an iota block, editing a domain byte — must fail loudly here.
func TestGloasWireConstants(t *testing.T) {
	// Beacon domains (consensus-specs Gloas), plus builder-specs' DOMAIN_BUILDER_REQUEST_AUTH
	// which shares the 0x0b first byte and differs only in the last byte.
	require.Equal(t, [4]byte{0x0b, 0x00, 0x00, 0x00}, DomainBeaconBuilder)
	require.Equal(t, [4]byte{0x0c, 0x00, 0x00, 0x00}, DomainPTCAttester)
	require.Equal(t, [4]byte{0x0d, 0x00, 0x00, 0x00}, DomainProposerPreferences)
	require.Equal(t, [4]byte{0x0b, 0x00, 0x00, 0x01}, DomainBuilderRequestAuth)

	// Beacon roles are iota-assigned, so inserting or reordering entries above them
	// would silently renumber them.
	require.Equal(t, BeaconRole(7), BNRolePTCAttester)
	require.Equal(t, BeaconRole(8), BNRoleProposerPreferences)
	require.Equal(t, BeaconRole(9), BNRoleEnvelopeProposer)

	require.Equal(t, RunnerRole(7), RolePTCAttester)
	require.Equal(t, RunnerRole(8), RoleProposerPreferences)
	require.Equal(t, RunnerRole(9), RoleEnvelopeProposer)

	require.Equal(t, PartialSigMsgType(7), PTCAttesterPartialSig)
	require.Equal(t, PartialSigMsgType(8), ProposerPreferencesPartialSig)
	require.Equal(t, PartialSigMsgType(9), RequestAuthPartialSig)
	require.Equal(t, PartialSigMsgType(10), EnvelopePartialSig)

	// SSVMessage.MsgType is iota-assigned, so appending or reordering the class above the
	// envelope-dissemination class would silently renumber it (SIP #94 §6).
	require.Equal(t, MsgType(3), SSVEnvelopeDisseminationMsgType)

	require.Equal(t, RolePTCAttester, MapDutyToRunnerRole(BNRolePTCAttester))
	require.Equal(t, RoleProposerPreferences, MapDutyToRunnerRole(BNRoleProposerPreferences))
	require.Equal(t, RoleEnvelopeProposer, MapDutyToRunnerRole(BNRoleEnvelopeProposer))

	require.Equal(t, "PTC_ATTESTER", BNRolePTCAttester.String())
	require.Equal(t, "PROPOSER_PREFERENCES", BNRoleProposerPreferences.String())
	require.Equal(t, "ENVELOPE_PROPOSER", BNRoleEnvelopeProposer.String())
	require.Equal(t, "PTC_ATTESTER_RUNNER", RolePTCAttester.String())
	require.Equal(t, "PROPOSER_PREFERENCES_RUNNER", RoleProposerPreferences.String())
	require.Equal(t, "ENVELOPE_PROPOSER_RUNNER", RoleEnvelopeProposer.String())
}
