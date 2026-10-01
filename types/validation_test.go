package types

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSSVMessageValidate(t *testing.T) {
	t.Run("invalid msg type", func(t *testing.T) {
		var committeeID CommitteeID
		msg := &SSVMessage{
			MsgType: MsgType(99),
			MsgID:   NewCommitteeMsgID(GenesisMainnet, committeeID, RoleCommittee),
		}
		err := msg.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, MessageTypeInvalidErrorCode, err.(*Error).Code)
	})

	t.Run("invalid role encoding in MsgID", func(t *testing.T) {
		var committeeID CommitteeID
		msg := &SSVMessage{
			MsgType: SSVConsensusMsgType,
			MsgID:   NewCommitteeMsgID(GenesisMainnet, committeeID, RoleUnknown),
		}
		err := msg.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, SSVMessageInvalidRoleErrorCode, err.(*Error).Code)
	})

}

func TestSignedSSVMessageValidateDelegatesToSSVMessageValidate(t *testing.T) {
	msg := &SignedSSVMessage{
		Signatures:  [][]byte{{1}},
		OperatorIDs: []OperatorID{1},
		SSVMessage: &SSVMessage{
			MsgType: SSVConsensusMsgType,
			MsgID:   NewCommitteeMsgID(GenesisMainnet, CommitteeID{}, RoleUnknown),
		},
	}

	err := msg.Validate()
	require.Error(t, err)
	require.ErrorIs(t, err, &Error{})
	require.Equal(t, SSVMessageInvalidRoleErrorCode, err.(*Error).Code)
}

func TestBeaconVoteValidate(t *testing.T) {
	t.Run("nil beacon vote", func(t *testing.T) {
		var bv *BeaconVote
		err := bv.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, BeaconVoteNilCheckpointErrorCode, err.(*Error).Code)
	})

	t.Run("zero-value beacon vote", func(t *testing.T) {
		err := (&BeaconVote{}).Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, BeaconVoteNilCheckpointErrorCode, err.(*Error).Code)
	})

	t.Run("nil source only", func(t *testing.T) {
		bv := &BeaconVote{
			Target: &phase0.Checkpoint{Epoch: 1},
		}
		err := bv.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, BeaconVoteNilCheckpointErrorCode, err.(*Error).Code)
	})

	t.Run("source epoch must be less than target epoch", func(t *testing.T) {
		bv := &BeaconVote{
			Source: &phase0.Checkpoint{Epoch: 2},
			Target: &phase0.Checkpoint{Epoch: 1},
		}

		err := bv.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, AttestationSourceNotLessThanTargetErrorCode, err.(*Error).Code)
	})
}

func TestCommitteeDutyValidate(t *testing.T) {
	t.Run("nil committee duty", func(t *testing.T) {
		var duty *CommitteeDuty
		err := duty.Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, InvalidCommitteeDutyErrorCode, err.(*Error).Code)
	})

	t.Run("nil validator duty", func(t *testing.T) {
		err := (&CommitteeDuty{
			Slot: 1,
			ValidatorDuties: []*ValidatorDuty{
				nil,
			},
		}).Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, InvalidCommitteeDutyErrorCode, err.(*Error).Code)
	})

	t.Run("mismatched validator duty slot", func(t *testing.T) {
		var pubKey phase0.BLSPubKey
		pubKey[0] = 1

		err := (&CommitteeDuty{
			Slot: 1,
			ValidatorDuties: []*ValidatorDuty{
				{
					Type:                    BNRoleAttester,
					PubKey:                  pubKey,
					Slot:                    2,
					CommitteeLength:         1,
					ValidatorCommitteeIndex: 0,
				},
			},
		}).Validate()
		require.Error(t, err)
		require.ErrorIs(t, err, &Error{})
		require.Equal(t, InvalidCommitteeDutyErrorCode, err.(*Error).Code)
	})
}
