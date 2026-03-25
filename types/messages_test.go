package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignedSSVMessageDeepCopyDoesNotAliasSlices(t *testing.T) {
	original := &SignedSSVMessage{
		OperatorIDs: []OperatorID{1, 2},
		Signatures: [][]byte{
			{1, 2, 3},
			{4, 5, 6},
		},
		SSVMessage: &SSVMessage{
			MsgType: 1,
			MsgID:   [56]byte{9, 8, 7},
			Data:    []byte{10, 11, 12},
		},
		FullData: []byte{13, 14, 15},
	}

	copied := original.DeepCopy()
	copied.OperatorIDs[0] = 99
	copied.Signatures[0][0] = 42
	copied.SSVMessage.Data[0] = 43
	copied.FullData[0] = 44

	require.Equal(t, OperatorID(1), original.OperatorIDs[0])
	require.Equal(t, byte(1), original.Signatures[0][0])
	require.Equal(t, byte(10), original.SSVMessage.Data[0])
	require.Equal(t, byte(13), original.FullData[0])
	require.NotSame(t, original.SSVMessage, copied.SSVMessage)
	require.NotSame(t, original.SSVMessage, copied.SSVMessage)
	// Aliasing already verified above via mutation (copied.Signatures[0][0] = 42 did not affect original)
}
