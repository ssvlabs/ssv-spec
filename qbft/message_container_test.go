package qbft_test

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestMsgContainer_AddIfDoesntExist(t *testing.T) {
	ks := testingutils.Testing4SharesSet()

	t.Run("same msg and signers", func(t *testing.T) {

		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		added, err := c.AddFirstMsgForSignerAndRound(testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))))
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddFirstMsgForSignerAndRound(testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))))
		require.NoError(t, err)
		require.False(t, added)
	})

	t.Run("same msg different signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		added, err := c.AddFirstMsgForSignerAndRound(testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1))))
		require.NoError(t, err)
		require.True(t, added)

		added, err = c.AddFirstMsgForSignerAndRound(testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(2))))
		require.NoError(t, err)
		require.True(t, added)
	})

	t.Run("same msg common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		m := testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))
		m.SignedMessage.OperatorIDs = []types.OperatorID{1, 2, 3, 4}
		added, err := c.AddFirstMsgForSignerAndRound(m)
		require.NoError(t, err)
		require.True(t, added)

		m = testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))
		m.SignedMessage.OperatorIDs = []types.OperatorID{1, 5, 6, 7}
		added, err = c.AddFirstMsgForSignerAndRound(m)
		require.NoError(t, err)
		require.True(t, added)
	})
}

func TestMsgContainer_Marshaling(t *testing.T) {
	ks := testingutils.Testing4SharesSet()

	c := &qbft.MsgContainer{
		Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
	}
	c.Msgs[1] = []*qbft.ProcessingMessage{testingutils.ToProcessingMessage(testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)))}

	byts, err := c.Encode()
	require.NoError(t, err)

	decoded := &qbft.MsgContainer{}
	require.NoError(t, decoded.Decode(byts))

	decodedByts, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, decodedByts)
}

func MessageWithSigners(signers []types.OperatorID, msg *qbft.Message) *qbft.ProcessingMessage {

	msgID := [56]byte{}
	copy(msgID[:], msg.Identifier)

	encodedMsg, err := msg.Encode()
	if err != nil {
		panic(err)
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    encodedMsg,
	}
	return testingutils.ToProcessingMessage(&types.SignedSSVMessage{
		OperatorIDs: signers,
		Signatures:  make([][]byte, len(signers)),
		SSVMessage:  ssvMsg,
	})
}

func TestMsgContainer_AddMsg(t *testing.T) {
	t.Run("same message, one with more signers", func(t *testing.T) {

		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		if err := c.AddMsg(MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Round: 1, Root: [32]byte{1, 2, 3, 4}})); err != nil {
			panic(err)
		}
		if err := c.AddMsg(MessageWithSigners([]types.OperatorID{1}, &qbft.Message{Round: 1, Root: [32]byte{1, 2, 3, 4}})); err != nil {
			panic(err)
		}
		if err := c.AddMsg(MessageWithSigners([]types.OperatorID{1, 2, 3, 4}, &qbft.Message{Round: 1, Root: [32]byte{1, 2, 3, 4}})); err != nil {
			panic(err)
		}

		cnt, msgs := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 4}, cnt)
		require.Len(t, msgs, 1)
	})
}

func TestMsgContainer_UniqueSignersSetForRoundAndValue(t *testing.T) {
	t.Run("multi common signers with different values", func(t *testing.T) {

		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{1, 2}, &qbft.Message{Root: [32]byte{1, 2, 3, 5}}),
			MessageWithSigners([]types.OperatorID{4}, &qbft.Message{Root: [32]byte{1, 2, 3, 6}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3}, cnt)

		cnt, _ = c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 6})
		require.EqualValues(t, []types.OperatorID{4}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{1, 2}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{4}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 4}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{1, 2, 5}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{4}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 4}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{1, 2, 5}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{4}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{3, 7, 8, 9, 10}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 5, 4, 3, 7, 8, 9, 10}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3}, cnt)
	})

	t.Run("multi common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{1}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3}, cnt)
	})

	t.Run("no common signers", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{6}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{4, 7}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(1, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{1, 2, 3, 6, 4, 7}, cnt)
	})

	t.Run("no round", func(t *testing.T) {
		c := &qbft.MsgContainer{
			Msgs: map[qbft.Round][]*qbft.ProcessingMessage{},
		}

		c.Msgs[1] = []*qbft.ProcessingMessage{
			MessageWithSigners([]types.OperatorID{1, 2, 3}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{6}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
			MessageWithSigners([]types.OperatorID{4, 7}, &qbft.Message{Root: [32]byte{1, 2, 3, 4}}),
		}
		cnt, _ := c.LongestUniqueSignersForRoundAndRoot(2, [32]byte{1, 2, 3, 4})
		require.EqualValues(t, []types.OperatorID{}, cnt)
	})
}
