package qbft_test

import (
	"testing"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
	"github.com/stretchr/testify/require"
)

func TestState_Decoding(t *testing.T) {
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
		ID:                []byte{1, 2, 3, 4},
		Round:             1,
		Height:            2,
		LastPreparedRound: 3,
		LastPreparedValue: []byte{1, 2, 3, 4},
		ProposalAcceptedForCurrentRound: &qbft.SignedMessage{
			Message: qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     1,
				Round:      2,
				Identifier: []byte{1, 2, 3, 4},
				Root:       testingutils.TestingSyncCommitteeBlockRoot,
			},
			Signature: []byte{1, 2, 3, 4},
			Signers:   []types.OperatorID{1},
		},
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

	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.Signature)
	require.EqualValues(t, []types.OperatorID{1}, decodedState.ProposalAcceptedForCurrentRound.Signers)
	require.EqualValues(t, qbft.CommitMsgType, decodedState.ProposalAcceptedForCurrentRound.Message.MsgType)
	require.EqualValues(t, 1, decodedState.ProposalAcceptedForCurrentRound.Message.Height)
	require.EqualValues(t, 2, decodedState.ProposalAcceptedForCurrentRound.Message.Round)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.Message.Identifier)
	require.EqualValues(t, testingutils.TestingSyncCommitteeBlockRoot, decodedState.ProposalAcceptedForCurrentRound.Message.Root)
}
