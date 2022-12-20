package frost

import (
	"testing"

	"github.com/bloxapp/ssv-spec/dkg"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

var testProtocolRound = Preparation

func testSignedMessage(round ProtocolRound) *dkg.SignedMessage {
	sk := testingutils.TestingKeygenKeySet().DKGOperators[1].SK
	msg := &dkg.Message{
		MsgType:    dkg.ProtocolMsgType,
		Identifier: testingutils.GetRandRequestID(),
	}
	switch round {
	case Preparation:
		msg.Data = Testing_PreparationMessageBytes(1, testingutils.KeygenMsgStore)
	case Round1:
		msg.Data = Testing_Round1MessageBytes(1, testingutils.KeygenMsgStore)
	case Round2:
		msg.Data = Testing_Round2MessageBytes(1, testingutils.KeygenMsgStore)
		// case Blame:
		// 	msg.Data = Testing_BlameMessageBytes(1, InvalidMessage, nil)
	}
	return testingutils.SignDKGMsg(sk, 1, msg)
}

func TestMsgContainer_SaveMsg(t *testing.T) {
	t.Run("new message", func(t *testing.T) {
		c := newMsgContainer()
		existingMessage, err := c.SaveMsg(testProtocolRound, testSignedMessage(testProtocolRound))
		require.NoError(t, err)
		require.Nil(t, existingMessage)
	})

	t.Run("message exist", func(t *testing.T) {
		c := newMsgContainer()
		testMsg := testSignedMessage(testProtocolRound)
		_, err := c.SaveMsg(testProtocolRound, testMsg)
		existingMessage, err := c.SaveMsg(testProtocolRound, testMsg)
		require.Error(t, err)
		require.NotNil(t, existingMessage)
	})
}
func TestMsgContainer_GetSignedMsg(t *testing.T)           {}
func TestMsgContainer_GetPreparationMsg(t *testing.T)      {}
func TestMsgContainer_GetRound1Msg(t *testing.T)           {}
func TestMsgContainer_GetRound2Msg(t *testing.T)           {}
func TestMsgContainer_GetBlameMsg(t *testing.T)            {}
func TestMsgContainer_GetMessage(t *testing.T)             {}
func TestMsgContainer_AllMessagesForRound(t *testing.T)    {}
func TestMsgContainer_AllMessagesReceivedFor(t *testing.T) {}
