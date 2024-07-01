package qbft_test

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestInstance_Marshaling(t *testing.T) {
	var TestingMessage = &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Root:       testingutils.TestingQBFTRootData,
	}
	keySet := testingutils.Testing4SharesSet()
	TestingRSASK := keySet.OperatorKeys[1]
	testingSignedMsg := testingutils.ToProcessingMessage(func() *types.SignedSSVMessage {
		return testingutils.SignQBFTMsg(TestingRSASK, 1, TestingMessage)
	}())
	testingCommitteeMember := testingutils.TestingCommitteeMember(keySet)

	i := &qbft.Instance{
		State: &qbft.State{
			CommitteeMember:                 testingCommitteeMember,
			ID:                              []byte{1, 2, 3, 4},
			Round:                           1,
			Height:                          1,
			LastPreparedRound:               1,
			LastPreparedValue:               []byte{1, 2, 3, 4},
			ProposalAcceptedForCurrentRound: testingSignedMsg,
			Decided:                         false,
			DecidedValue:                    []byte{1, 2, 3, 4},

			ProposeContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			PrepareContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			CommitContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			RoundChangeContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.ProcessingMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
		},
	}

	byts, err := i.Encode()
	require.NoError(t, err)

	decoded := &qbft.Instance{}
	require.NoError(t, decoded.Decode(byts))

	bytsDecoded, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, bytsDecoded)
}
