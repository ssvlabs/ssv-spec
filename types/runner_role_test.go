package types

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunnerRoleWireMapping(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint32(math.MaxUint32), RoleUnknown.WireUint32())
	require.Equal(t, RoleUnknown, RunnerRoleFromWireUint32(uint32(math.MaxUint32)))

	// Arbitrary negative role (not RoleUnknown) should also map to sentinel.
	require.Equal(t, uint32(math.MaxUint32), RunnerRole(-42).WireUint32())

	for _, role := range []RunnerRole{
		RoleCommittee,
		RoleProposer,
		RoleValidatorRegistration,
		RoleVoluntaryExit,
		RoleAggregatorCommittee,
	} {
		require.Equal(t, role, RunnerRoleFromWireUint32(role.WireUint32()))
	}

	// Values outside the int32 range must not produce surprising negative roles.
	require.Equal(t, RoleUnknown, RunnerRoleFromWireUint32(0x80000000))
}

func TestMessageIDRoleTypeDecoding(t *testing.T) {
	t.Parallel()

	dutyExecutorID := make([]byte, dutyExecutorIDSize)
	for i := range dutyExecutorID {
		dutyExecutorID[i] = byte(i)
	}

	var validatorPK ValidatorPK
	copy(validatorPK[:], dutyExecutorID)
	msgID := NewValidatorMsgID(GenesisMainnet, validatorPK, RoleUnknown)
	require.Equal(t, RoleUnknown, msgID.GetRoleType())

	roleWire := binary.LittleEndian.Uint32(msgID[roleTypeStartPos : roleTypeStartPos+roleTypeSize])
	require.Equal(t, uint32(math.MaxUint32), roleWire)
}
