package qbft_test

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgContainer_AddIfDoesntExist(t *testing.T) {
	ks := testingutils.Testing4SharesSet()

	t.Run("same msg and signers", func(t *testing.T) {

		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		added, err := c.AddFirstMsgForSignerAndRound(testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)))
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddFirstMsgForSignerAndRound(testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)))
		require.NoError(t, err)
		require.False(t, added)
	})

	t.Run("same msg different signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		added, err := c.AddFirstMsgForSignerAndRound(testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)))
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddFirstMsgForSignerAndRound(testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(2)))
		require.NoError(t, err)
		require.True(t, added)
	})

	t.Run("same msg common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		m := testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))
		m.Signers = []types.OperatorID{1, 2, 3, 4}
		added, err := c.AddFirstMsgForSignerAndRound(m)
		require.NoError(t, err)
		require.True(t, added)

		m = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))
		m.Signers = []types.OperatorID{1, 5, 6, 7}
		added, err = c.AddFirstMsgForSignerAndRound(m)
		require.NoError(t, err)
		require.True(t, added)
	})
}

func TestMsgContainer_Marshaling(t *testing.T) {
	ks := testingutils.Testing4SharesSet()

	c := &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.SignedMessage{},
	}
	c.Msgs[1] = []*qbft.SignedMessage{testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))}

	byts, err := c.Encode()
	require.NoError(t, err)

	decoded := &qbft.MsgContainer{}
	require.NoError(t, decoded.Decode(byts))

	decodedByts, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, decodedByts)
}

func TestMsgContainer_AddMsg(t *testing.T) {
	t.Run("same message, one with more signers", func(t *testing.T) {

		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.AddMsg(&qbft.SignedMessage{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Round: 1, Root: [32]byte{1, 2, 3, 4}}})
		c.AddMsg(&qbft.SignedMessage{Signers: []types.OperatorID{1}, Message: qbft.Message{Round: 1, Root: [32]byte{1, 2, 3, 4}}})
		c.AddMsg(&qbft.SignedMessage{Signers: []types.OperatorID{1, 2, 3, 4}, Message: qbft.Message{Round: 1, Root: [32]byte{1, 2, 3, 4}}})

		cnt, msgs := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 4}, cnt)
		require.Len(t, msgs, 1)
	})
}

func TestMsgContainer_UniqueSignersSetForRoundAndValue(t *testing.T) {
	t.Run("multi common signers with different values", func(t *testing.T) {

		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{1, 2}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 5}}},
			{Signers: []types.OperatorID{4}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 6}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3}, cnt)

		cnt, _ = c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 6})
		require.EqualValues(t, []types.OperatorID{4}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{1, 2}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{4}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 4}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{1, 2, 5}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{4}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 4}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{1, 2, 5}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{4}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{3, 7, 8, 9, 10}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 5, 4, 3, 7, 8, 9, 10}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{1}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3}, cnt)
	})

	t.Run("no common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{6}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{4, 7}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 6, 4, 7}, cnt)
	})

	t.Run("no round", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.SignedMessage{},
		}

		c.Msgs[1] = []*qbft.SignedMessage{
			{Signers: []types.OperatorID{1, 2, 3}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{6}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
			{Signers: []types.OperatorID{4, 7}, Message: qbft.Message{Root: [32]byte{1, 2, 3, 4}}},
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(2, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{}, cnt)
	})
}
