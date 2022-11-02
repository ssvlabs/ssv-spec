package qbft

import (
	"github.com/bloxapp/ssv-spec/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestState_Decoding(t *testing.T) {
	inputData := &Data{
		Root:   [32]byte{1, 2, 3, 4},
		Source: []byte{1, 2, 3, 4},
	}
	state := &State{
		Share: &types.Share{
			OperatorID:      1,
			ValidatorPubKey: []byte{1, 2, 3, 4},
			Committee: []*types.Operator{
				{
					OperatorID: 1,
					PubKey:     []byte{1, 2, 3, 4},
				},
			},
			DomainType: types.PrimusTestnet,
		},
		//ID:                []byte{1, 2, 3, 4},
		ID:                types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester),
		Round:             1,
		Height:            2,
		LastPreparedRound: 3,
		LastPreparedValue: inputData,
		ProposalAcceptedForCurrentRound: &SignedMessage{
			Message: &Message{
				Height: 1,
				Round:  2,
				Input:  inputData,
			},
			Signature: []byte{1, 2, 3, 4},
			Signers:   []types.OperatorID{1},
		},
	}

	byts, err := state.Encode()
	require.NoError(t, err)

	decodedState := &State{}
	require.NoError(t, decodedState.Decode(byts))

	require.EqualValues(t, 1, decodedState.Share.OperatorID)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Share.ValidatorPubKey)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Share.Committee[0].PubKey)
	require.EqualValues(t, 1, decodedState.Share.Committee[0].OperatorID)
	require.EqualValues(t, types.PrimusTestnet, decodedState.Share.DomainType)

	require.EqualValues(t, 3, decodedState.LastPreparedRound)
	require.EqualValues(t, [32]byte{1, 2, 3, 4}, decodedState.LastPreparedValue.Root)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.LastPreparedValue.Source)
	require.True(t, decodedState.ID.Compare(types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)))
	require.EqualValues(t, 2, decodedState.Height)
	require.EqualValues(t, 1, decodedState.Round)

	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.Signature)
	require.EqualValues(t, []types.OperatorID{1}, decodedState.ProposalAcceptedForCurrentRound.Signers)
	require.EqualValues(t, 1, decodedState.ProposalAcceptedForCurrentRound.Message.Height)
	require.EqualValues(t, 2, decodedState.ProposalAcceptedForCurrentRound.Message.Round)
	require.EqualValues(t, [32]byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.Message.Input.Root)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ProposalAcceptedForCurrentRound.Message.Input.Source)
}
