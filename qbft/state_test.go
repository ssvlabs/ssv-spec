package qbft_test

import (
	"testing"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestState_Decoding(t *testing.T) {

	ks := testingutils.Testing4SharesSet()
	proposalMessage := testingutils.SignQBFTMsg(ks.NetworkKeys[1], 1, &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     1,
		Round:      2,
		Identifier: []byte{1, 2, 3, 4},
		Root:       testingutils.TestingSyncCommitteeBlockRoot,
	})

	state := &qbft.State{
		Share: &types.Share{
			OperatorID:      1,
			ValidatorPubKey: []byte{1, 2, 3, 4},
			Committee: []*types.Operator{
				{
					OperatorID:   1,
					BeaconPubKey: []byte{1, 2, 3, 4},
				},
			},
			DomainType: testingutils.TestingSSVDomainType,
		},
		ID:                              []byte{1, 2, 3, 4},
		Round:                           1,
		Height:                          2,
		LastPreparedRound:               3,
		LastPreparedValue:               []byte{1, 2, 3, 4},
		ProposalAcceptedForCurrentRound: proposalMessage,
	}

	byts, err := state.Encode()
	require.NoError(t, err)

	decodedState := &qbft.State{}
	require.NoError(t, decodedState.Decode(byts))

	require.EqualValues(t, 1, decodedState.Share.OperatorID)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Share.ValidatorPubKey)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Share.Committee[0].BeaconPubKey)
	require.EqualValues(t, 1, decodedState.Share.Committee[0].OperatorID)
	require.EqualValues(t, testingutils.TestingSSVDomainType, decodedState.Share.DomainType)

	require.EqualValues(t, 3, decodedState.LastPreparedRound)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.LastPreparedValue)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ID)
	require.EqualValues(t, 2, decodedState.Height)
	require.EqualValues(t, 1, decodedState.Round)

	decodedProposal := &qbft.Message{}
	require.NoError(t, decodedProposal.Decode(decodedState.ProposalAcceptedForCurrentRound.SSVMessage.Data))

	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.Signature)
	require.EqualValues(t, []types.OperatorID{1}, decodedState.ProposalAcceptedForCurrentRound.OperatorID)
	require.EqualValues(t, qbft.CommitMsgType, decodedProposal.MsgType)
	require.EqualValues(t, 1, decodedProposal.Height)
	require.EqualValues(t, 2, decodedProposal.Round)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedProposal.Identifier)
	require.EqualValues(t, testingutils.TestingSyncCommitteeBlockRoot, decodedProposal.Root)
}
