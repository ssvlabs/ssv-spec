package frost

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/stretchr/testify/require"
)

func TestMsgContainer_SaveMsg(t *testing.T) {
	t.Run("new message", func(t *testing.T) {
		c := newMsgContainer()

		existingMessage, err := c.SaveMsg(testProtocolRound, testSignedMessage(testProtocolRound, 1))
		require.NoError(t, err)
		require.Nil(t, existingMessage)
	})

	t.Run("message exist", func(t *testing.T) {
		testMsg := testSignedMessage(testProtocolRound, 1)
		c := newMsgContainer()
		_, _ = c.SaveMsg(testProtocolRound, testMsg)

		existingMessage, err := c.SaveMsg(testProtocolRound, testMsg)
		require.Error(t, err)
		require.NotNil(t, existingMessage)
	})
}

func TestMsgContainer_GetSignedMsg(t *testing.T) {
	t.Run("signed message not found", func(t *testing.T) {
		c := newMsgContainer()

		returnedMsg, err := c.GetSignedMsg(testProtocolRound, 1)
		require.Error(t, err)
		require.Nil(t, returnedMsg)
	})

	t.Run("signed message exist", func(t *testing.T) {
		c := newMsgContainer()
		testMsg := testSignedMessage(testProtocolRound, 1)
		_, _ = c.SaveMsg(testProtocolRound, testMsg)

		returnedMsg, err := c.GetSignedMsg(testProtocolRound, 1)
		require.NoError(t, err)
		require.NotNil(t, returnedMsg)
	})
}

func TestMsgContainer_GetPreparationMsg(t *testing.T) {
	t.Run("preparation message exists", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Preparation, testSignedMessage(Preparation, 1))

		returnedMsg, err := c.GetPreparationMsg(1)
		require.NoError(t, err)
		require.NotNil(t, returnedMsg)
	})

	t.Run("preparation message doesn't exists", func(t *testing.T) {
		c := newMsgContainer()

		returnedMsg, err := c.GetPreparationMsg(1)
		require.Error(t, err)
		require.Nil(t, returnedMsg)
	})

	t.Run("preparation message is nil", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Preparation, testSignedMessage(Round1, 1))

		returnedMsg, err := c.GetPreparationMsg(1)
		require.ErrorIs(t, err, ErrMsgNil{round: Preparation, operatorID: 1})
		require.Nil(t, returnedMsg)
	})
}

func TestMsgContainer_GetRound1Msg(t *testing.T) {
	t.Run("round1 message exists", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Round1, testSignedMessage(Round1, 1))

		returnedMsg, err := c.GetRound1Msg(1)
		require.NoError(t, err)
		require.NotNil(t, returnedMsg)
	})

	t.Run("round1 message doesn't exists", func(t *testing.T) {
		c := newMsgContainer()

		returnedMsg, err := c.GetRound1Msg(1)
		require.Error(t, err)
		require.Nil(t, returnedMsg)
	})

	t.Run("round1 message is nil", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Round1, testSignedMessage(Preparation, 1))

		returnedMsg, err := c.GetRound1Msg(1)
		require.ErrorIs(t, err, ErrMsgNil{round: Round1, operatorID: 1})
		require.Nil(t, returnedMsg)
	})
}

func TestMsgContainer_GetRound2Msg(t *testing.T) {
	t.Run("round2 message exists", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Round2, testSignedMessage(Round2, 1))

		returnedMsg, err := c.GetRound2Msg(1)
		require.NoError(t, err)
		require.NotNil(t, returnedMsg)
	})

	t.Run("round2 message doesn't exists", func(t *testing.T) {
		c := newMsgContainer()

		returnedMsg, err := c.GetRound2Msg(1)
		require.Error(t, err)
		require.Nil(t, returnedMsg)
	})

	t.Run("round2 message is nil", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Round2, testSignedMessage(Preparation, 1))

		returnedMsg, err := c.GetRound2Msg(1)
		require.ErrorIs(t, err, ErrMsgNil{round: Round2, operatorID: 1})
		require.Nil(t, returnedMsg)
	})
}

func TestMsgContainer_GetBlameMsg(t *testing.T) {
	t.Run("round2 message exists", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Blame, testSignedMessage(Blame, 1))

		returnedMsg, err := c.GetBlameMsg(1)
		require.NoError(t, err)
		require.NotNil(t, returnedMsg)
	})

	t.Run("round2 message doesn't exists", func(t *testing.T) {
		c := newMsgContainer()

		returnedMsg, err := c.GetBlameMsg(1)
		require.Error(t, err)
		require.Nil(t, returnedMsg)
	})

	t.Run("round2 message is nil", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Blame, testSignedMessage(Preparation, 1))

		returnedMsg, err := c.GetBlameMsg(1)
		require.ErrorIs(t, err, ErrMsgNil{round: Blame, operatorID: 1})
		require.Nil(t, returnedMsg)
	})
}

func TestMsgContainer_AllMessagesForRound(t *testing.T) {
	t.Run("default case", func(t *testing.T) {
		c := newMsgContainer()
		expected := map[uint32]*dkg.SignedMessage{
			1: testSignedMessage(Preparation, 1),
			2: testSignedMessage(Preparation, 2),
		}
		_, _ = c.SaveMsg(Preparation, expected[1])
		_, _ = c.SaveMsg(Preparation, expected[2])

		actual := c.AllMessagesForRound(Preparation)
		require.EqualValues(t, expected, actual)
	})
}

func TestMsgContainer_AllMessagesReceivedFor(t *testing.T) {
	t.Run("true case", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Preparation, testSignedMessage(Preparation, 1))
		_, _ = c.SaveMsg(Preparation, testSignedMessage(Preparation, 2))

		ok := c.AllMessagesReceivedFor(Preparation, []uint32{1, 2})
		require.EqualValues(t, true, ok)
	})

	t.Run("false case", func(t *testing.T) {
		c := newMsgContainer()
		_, _ = c.SaveMsg(Preparation, testSignedMessage(Preparation, 1))

		ok := c.AllMessagesReceivedFor(Preparation, []uint32{1, 2})
		require.EqualValues(t, false, ok)
	})
}
