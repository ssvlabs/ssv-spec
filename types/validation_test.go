package types

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSSVMessageValidate(t *testing.T) {
	t.Run("invalid msg type", func(t *testing.T) {
		msg := &SSVMessage{
			MsgType: MsgType(99),
			MsgID:   NewMsgID(GenesisMainnet, make([]byte, 48), RoleCommittee),
		}
		err := msg.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, MessageTypeInvalidErrorCode, err.(*Error).Code)
	})

	t.Run("invalid role encoding in MsgID", func(t *testing.T) {
		msg := &SSVMessage{
			MsgType: SSVConsensusMsgType,
			MsgID:   NewMsgID(GenesisMainnet, make([]byte, 48), RoleUnknown),
		}
		err := msg.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, SSVMessageInvalidRoleErrorCode, err.(*Error).Code)
	})

	t.Run("data too large", func(t *testing.T) {
		msg := &SSVMessage{
			MsgType: SSVConsensusMsgType,
			MsgID:   NewMsgID(GenesisMainnet, make([]byte, 48), RoleCommittee),
			Data:    make([]byte, maxSSVMessageDataSize+1),
		}
		err := msg.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, SSVMessageDataTooLargeErrorCode, err.(*Error).Code)
	})
}

func TestSignedSSVMessageValidateDelegatesToSSVMessageValidate(t *testing.T) {
	msg := &SignedSSVMessage{
		Signatures:  [][]byte{{1}},
		OperatorIDs: []OperatorID{1},
		SSVMessage: &SSVMessage{
			MsgType: SSVConsensusMsgType,
			MsgID:   NewMsgID(GenesisMainnet, make([]byte, 48), RoleUnknown),
		},
	}

	err := msg.Validate()
	require.Error(t, err)
	require.ErrorIs(t, err, &Error{})
	require.Equal(t, SSVMessageInvalidRoleErrorCode, err.(*Error).Code)
}

func TestBeaconVoteValidate(t *testing.T) {
	t.Run("nil checkpoints", func(t *testing.T) {
		bv := &BeaconVote{Source: nil, Target: nil}
		err := bv.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, BeaconVoteNilCheckpointErrorCode, err.(*Error).Code)
	})

	bv := &BeaconVote{
		Source: &phase0.Checkpoint{Epoch: 2},
		Target: &phase0.Checkpoint{Epoch: 1},
	}

	err := bv.Validate()
	require.Error(t, err)
	require.ErrorIs(t, err, &Error{})
	require.Equal(t, AttestationSourceNotLessThanTargetErrorCode, err.(*Error).Code)
}
