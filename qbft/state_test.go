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
	signedProposalMsg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{{1, 2, 3, 4}},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   [56]byte{1, 2, 3, 4},
			Data:    proposalMsgBytes,
		},
	}

	ks := testingutils.Testing4SharesSet()

	state := &qbft.State{
		SharedValidator:                 testingutils.TestingSharedValidator(ks, testingutils.TestingValidatorIndex),
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

	require.EqualValues(t, 1, decodedState.SharedValidator.OwnValidatorShare.OperatorID)
	//require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Operator.ValidatorPubKey)
	//require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.Operator.Committee[0].SharePubKey)
	require.EqualValues(t, 1, decodedState.SharedValidator.Committee[0].OperatorID)
	//require.EqualValues(t, testingutils.TestingSSVDomainType, decodedState.Operator.DomainType)

	require.EqualValues(t, 3, decodedState.LastPreparedRound)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.LastPreparedValue)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedState.ID)
	require.EqualValues(t, 2, decodedState.Height)
	require.EqualValues(t, 1, decodedState.Round)

	decodedProposalMsg := &qbft.Message{}
	if err := decodedProposalMsg.Decode(decodedState.ProposalAcceptedForCurrentRound.SSVMessage.Data); err != nil {
		panic(err)
	}

	require.EqualValues(t, [][]byte{{1, 2, 3, 4}}, decodedState.ProposalAcceptedForCurrentRound.Signatures)
	require.EqualValues(t, []types.OperatorID{1}, decodedState.ProposalAcceptedForCurrentRound.GetOperatorIDs())
	require.EqualValues(t, qbft.CommitMsgType, decodedProposalMsg.MsgType)
	require.EqualValues(t, 1, decodedProposalMsg.Height)
	require.EqualValues(t, 2, decodedProposalMsg.Round)
	require.EqualValues(t, []byte{1, 2, 3, 4}, decodedProposalMsg.Identifier)
	require.EqualValues(t, testingutils.TestingSyncCommitteeBlockRoot, decodedProposalMsg.Root)
}
