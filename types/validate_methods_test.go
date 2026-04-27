package types

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestDomainTypeIsKnown(t *testing.T) {
	require.True(t, GenesisMainnet.IsKnown())
	require.True(t, JatoTestnet.IsKnown())
	require.False(t, DomainType{0xaa, 0xbb, 0xcc, 0xdd}.IsKnown())
}

func TestOperatorValidate(t *testing.T) {
	valid := func() *Operator {
		return &Operator{
			OperatorID:        1,
			SSVOperatorPubKey: make([]byte, 459),
		}
	}

	t.Run("valid", func(t *testing.T) {
		require.NoError(t, valid().Validate())
	})

	t.Run("zero operator id", func(t *testing.T) {
		op := valid()
		op.OperatorID = 0
		err := op.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidOperatorErrorCode, err.(*Error).Code)
	})

}

func TestCommitteeMemberValidate(t *testing.T) {
	valid := func() *CommitteeMember {
		ops := []*Operator{
			{OperatorID: 1, SSVOperatorPubKey: make([]byte, 459)},
			{OperatorID: 2, SSVOperatorPubKey: make([]byte, 459)},
			{OperatorID: 3, SSVOperatorPubKey: make([]byte, 459)},
			{OperatorID: 4, SSVOperatorPubKey: make([]byte, 459)},
		}
		opIDs := []OperatorID{1, 2, 3, 4}
		return &CommitteeMember{
			OperatorID:        1,
			CommitteeID:       GetCommitteeID(opIDs),
			SSVOperatorPubKey: make([]byte, 459),
			FaultyNodes:       1, // n=4 == 3f+1
			Committee:         ops,
			DomainType:        JatoTestnet,
		}
	}

	t.Run("valid", func(t *testing.T) {
		require.NoError(t, valid().Validate())
	})

	t.Run("faulty nodes mismatch committee size", func(t *testing.T) {
		cm := valid()
		cm.FaultyNodes = 0
		err := cm.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidCommitteeMemberErrorCode, err.(*Error).Code)
	})

	t.Run("duplicate operator IDs", func(t *testing.T) {
		cm := valid()
		cm.Committee[1].OperatorID = 1
		err := cm.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidCommitteeMemberErrorCode, err.(*Error).Code)
	})

	t.Run("committee id mismatch", func(t *testing.T) {
		cm := valid()
		cm.CommitteeID = CommitteeID{}
		err := cm.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidCommitteeMemberErrorCode, err.(*Error).Code)
	})
}

func TestShareValidate(t *testing.T) {
	valid := func() *Share {
		var pk ValidatorPK
		pk[0] = 1
		return &Share{
			ValidatorIndex:  1,
			ValidatorPubKey: pk,
			SharePubKey:     make([]byte, 48),
			Committee: []*ShareMember{
				{SharePubKey: make([]byte, 48), Signer: 1},
				{SharePubKey: make([]byte, 48), Signer: 2},
				{SharePubKey: make([]byte, 48), Signer: 3},
				{SharePubKey: make([]byte, 48), Signer: 4},
			},
			DomainType:          JatoTestnet,
			FeeRecipientAddress: [20]byte{},
			Graffiti:            make([]byte, 32),
		}
	}

	t.Run("valid", func(t *testing.T) {
		require.NoError(t, valid().Validate())
	})

	t.Run("zero validator pubkey", func(t *testing.T) {
		s := valid()
		s.ValidatorPubKey = ValidatorPK{}
		err := s.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidShareErrorCode, err.(*Error).Code)
	})
}

func TestValidatorDutyValidate(t *testing.T) {
	valid := func() *ValidatorDuty {
		var pk phase0.BLSPubKey
		pk[0] = 1
		return &ValidatorDuty{
			Type:                    BNRoleAttester,
			PubKey:                  pk,
			Slot:                    1,
			ValidatorIndex:          1,
			CommitteeIndex:          1,
			CommitteeLength:         4,
			CommitteesAtSlot:        1,
			ValidatorCommitteeIndex: 0,
		}
	}

	t.Run("valid", func(t *testing.T) {
		require.NoError(t, valid().Validate())
	})

	t.Run("unknown duty type", func(t *testing.T) {
		d := valid()
		d.Type = BNRoleUnknown
		err := d.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidValidatorDutyErrorCode, err.(*Error).Code)
	})

	t.Run("committee index out of bounds", func(t *testing.T) {
		d := valid()
		d.ValidatorCommitteeIndex = d.CommitteeLength
		err := d.Validate()
		require.Error(t, err)
		require.Equal(t, InvalidValidatorDutyErrorCode, err.(*Error).Code)
	})

	t.Run("sync committee allows zero committee metadata", func(t *testing.T) {
		d := valid()
		d.Type = BNRoleSyncCommittee
		d.CommitteeLength = 0
		d.ValidatorCommitteeIndex = 0
		require.NoError(t, d.Validate())
	})

	t.Run("sync committee contribution allows zero committee metadata", func(t *testing.T) {
		d := valid()
		d.Type = BNRoleSyncCommitteeContribution
		d.CommitteeLength = 0
		d.ValidatorCommitteeIndex = 0
		require.NoError(t, d.Validate())
	})
}
