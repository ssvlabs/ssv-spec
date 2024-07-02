package qbft_test

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
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
	signedProposalMsg := testingutils.ToProcessingMessage(&types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{{1, 2, 3, 4}},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   [56]byte{1, 2, 3, 4},
			Data:    proposalMsgBytes,
		},
	})

	state := &qbft.State{
		CommitteeMember: &types.CommitteeMember{
			OperatorID: 1,
			Committee: []*types.Operator{
				{
					OperatorID:        1,
					SSVOperatorPubKey: []byte{1, 2, 3, 4},
				},
			},
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

	require.EqualValues(t, 1, decodedState.CommitteeMember.OperatorID)
	//require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.CommitteeMember.ValidatorPubKey)
	//require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.CommitteeMember.Committee[0].SharePubKey)
	require.EqualValues(t, 1, decodedState.CommitteeMember.Committee[0].OperatorID)
	//require.EqualValues(t, testingutils.TestingSSVDomainType, decodedState.CommitteeMember.DomainType)

	require.EqualValues(t, 3, decodedState.LastPreparedRound)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.LastPreparedValue)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ID)
	require.EqualValues(t, 2, decodedState.Height)
	require.EqualValues(t, 1, decodedState.Round)

	require.EqualValues(t, [][]byte{{1, 2, 3, 4}}, decodedState.ProposalAcceptedForCurrentRound.SignedMessage.Signatures)
	require.EqualValues(t, []types.OperatorID{1}, decodedState.ProposalAcceptedForCurrentRound.SignedMessage.OperatorIDs)
	require.EqualValues(t, qbft.CommitMsgType, decodedState.ProposalAcceptedForCurrentRound.QBFTMessage.MsgType)
	require.EqualValues(t, 1, decodedState.ProposalAcceptedForCurrentRound.QBFTMessage.Height)
	require.EqualValues(t, 2, decodedState.ProposalAcceptedForCurrentRound.QBFTMessage.Round)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.QBFTMessage.Identifier)
	require.EqualValues(t, testingutils.TestingSyncCommitteeBlockRoot, decodedState.ProposalAcceptedForCurrentRound.QBFTMessage.Root)
}
