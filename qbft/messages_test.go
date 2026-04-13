package qbft_test

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestRoundChangeJustificationProcessingMessagesReturnsDecodeError(t *testing.T) {
	keySet := testingutils.Testing4SharesSet()
	justification := testingutils.TestingProposalMessage(keySet.OperatorKeys[1], types.OperatorID(1))
	justification.SSVMessage.Data = []byte{1, 2, 3}

	justifications, err := qbft.MarshalJustifications([]*types.SignedSSVMessage{justification})
	require.NoError(t, err)

	msg := &qbft.Message{
		RoundChangeJustification: justifications,
	}

	processingMessages, err := msg.RoundChangeJustificationProcessingMessages()
	require.Nil(t, processingMessages)
	require.ErrorContains(t, err, "decode justification message")
}

func TestRoundChangeJustificationProcessingMessages(t *testing.T) {
	keySet := testingutils.Testing4SharesSet()
	justification := testingutils.TestingRoundChangeMessageWithRound(keySet.OperatorKeys[1], types.OperatorID(1), 2)

	justifications, err := qbft.MarshalJustifications([]*types.SignedSSVMessage{justification})
	require.NoError(t, err)

	msg := &qbft.Message{
		RoundChangeJustification: justifications,
	}

	processingMessages, err := msg.RoundChangeJustificationProcessingMessages()
	require.NoError(t, err)
	require.Len(t, processingMessages, 1)
	require.Equal(t, qbft.RoundChangeMsgType, processingMessages[0].QBFTMessage.MsgType)
	require.Equal(t, qbft.FirstHeight, processingMessages[0].QBFTMessage.Height)
	require.Equal(t, qbft.Round(2), processingMessages[0].QBFTMessage.Round)
}

func TestPrepareJustificationProcessingMessages(t *testing.T) {
	keySet := testingutils.Testing4SharesSet()
	justification := testingutils.TestingPrepareMessage(keySet.OperatorKeys[1], types.OperatorID(1))

	justifications, err := qbft.MarshalJustifications([]*types.SignedSSVMessage{justification})
	require.NoError(t, err)

	msg := &qbft.Message{
		PrepareJustification: justifications,
	}

	processingMessages, err := msg.PrepareJustificationProcessingMessages()
	require.NoError(t, err)
	require.Len(t, processingMessages, 1)
	require.Equal(t, qbft.PrepareMsgType, processingMessages[0].QBFTMessage.MsgType)
	require.Equal(t, qbft.FirstHeight, processingMessages[0].QBFTMessage.Height)
	require.Equal(t, qbft.FirstRound, processingMessages[0].QBFTMessage.Round)
}

func TestPrepareJustificationProcessingMessagesReturnsDecodeError(t *testing.T) {
	keySet := testingutils.Testing4SharesSet()
	justification := testingutils.TestingPrepareMessage(keySet.OperatorKeys[1], types.OperatorID(1))
	justification.SSVMessage.Data = []byte{1, 2, 3}

	justifications, err := qbft.MarshalJustifications([]*types.SignedSSVMessage{justification})
	require.NoError(t, err)

	msg := &qbft.Message{
		PrepareJustification: justifications,
	}

	processingMessages, err := msg.PrepareJustificationProcessingMessages()
	require.Nil(t, processingMessages)
	require.ErrorContains(t, err, "decode justification message")
}

