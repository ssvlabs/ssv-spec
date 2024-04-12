package qbft_test

import (
	"testing"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestState_Decoding(t *testing.T) {

	proposalMsg := &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     1,
		Round:      2,
		Identifier: []byte{1, 2, 3, 4},
		Root:       testingutils.TestingSyncCommitteeBlockRoot,
	}
	proposalMsgBytes, err := proposalMsg.Encode()
	if err != nil {
		panic(err)
	}
	signedProposalMsg := &types.SignedSSVMessage{
		OperatorID: []types.OperatorID{1},
		Signature:  [][]byte{{1, 2, 3, 4}},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   [56]byte{1, 2, 3, 4},
			Data:    proposalMsgBytes,
		},
	}

	state := &qbft.State{
		Share: &types.Share{
			OperatorID:      1,
			ValidatorPubKey: []byte{1, 2, 3, 4},
			Committee: []*types.Operator{
				{
					OperatorID:  1,
					SharePubKey: []byte{1, 2, 3, 4},
				},
			},
			DomainType: testingutils.TestingSSVDomainType,
		},
		ID:                              []byte{1, 2, 3, 4},
		Round:                           1,
		Height:                          2,
		LastPreparedRound:               3,
		LastPreparedValue:               []byte{1, 2, 3, 4},
		ProposalAcceptedForCurrentRound: signedProposalMsg,
	}

	byts, err := state.Encode()
	require.NoError(t, err)

	decodedState := &qbft.State{}
	require.NoError(t, decodedState.Decode(byts))

	require.EqualValues(t, 1, decodedState.Share.OperatorID)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Share.ValidatorPubKey)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Share.Committee[0].SharePubKey)
	require.EqualValues(t, 1, decodedState.Share.Committee[0].OperatorID)
	require.EqualValues(t, testingutils.TestingSSVDomainType, decodedState.Share.DomainType)

	require.EqualValues(t, 3, decodedState.LastPreparedRound)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.LastPreparedValue)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ID)
	require.EqualValues(t, 2, decodedState.Height)
	require.EqualValues(t, 1, decodedState.Round)

	decodedProposalMsg := &qbft.Message{}
	if err := decodedProposalMsg.Decode(decodedState.ProposalAcceptedForCurrentRound.SSVMessage.Data); err != nil {
		panic(err)
	}

	require.EqualValues(t, [][]byte{{1, 2, 3, 4}}, decodedState.ProposalAcceptedForCurrentRound.Signature)
	require.EqualValues(t, []types.OperatorID{1}, decodedState.ProposalAcceptedForCurrentRound.GetOperatorID())
	require.EqualValues(t, qbft.CommitMsgType, decodedProposalMsg.MsgType)
	require.EqualValues(t, 1, decodedProposalMsg.Height)
	require.EqualValues(t, 2, decodedProposalMsg.Round)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedProposalMsg.Identifier)
	require.EqualValues(t, testingutils.TestingSyncCommitteeBlockRoot, decodedProposalMsg.Root)
}
