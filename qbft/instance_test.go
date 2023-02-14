package qbft

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstance_Marshaling(t *testing.T) {
	var TestingMessage = &Message{
		MsgType:    ProposalMsgType,
		Height:     FirstHeight,
		Round:      FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Root:       testingutils.TestingQBFTRootData,
	}
	TestingSK := func() *bls.SecretKey {
		types.InitBLS()
		ret := &bls.SecretKey{}
		ret.SetByCSPRNG()
		return ret
	}()
	testingSignedMsg := func() *SignedMessage {
		return testingutils.SignQBFTMsg(TestingSK, 1, TestingMessage)
	}()
	testingValidatorPK := spec.BLSPubKey{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	testingShare := &types.Share{
		OperatorID:      1,
		ValidatorPubKey: testingValidatorPK[:],
		SharePubKey:     TestingSK.GetPublicKey().Serialize(),
		DomainType:      types.PrimusTestnet,
		Quorum:          3,
		PartialQuorum:   2,
		Committee: []*types.Operator{
			{
				OperatorID: 1,
				PubKey:     TestingSK.GetPublicKey().Serialize(),
			},
		},
	}
	i := &Instance{
		State: &State{
			Share:                           testingShare,
			ID:                              []byte{1, 2, 3, 4},
			Round:                           1,
			Height:                          1,
			LastPreparedRound:               1,
			LastPreparedValue:               []byte{1, 2, 3, 4},
			ProposalAcceptedForCurrentRound: testingSignedMsg,
			Decided:                         false,
			DecidedValue:                    []byte{1, 2, 3, 4},

			ProposeContainer: &MsgContainer{
				Msgs: map[Round][]*SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			PrepareContainer: &MsgContainer{
				Msgs: map[Round][]*SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			CommitContainer: &MsgContainer{
				Msgs: map[Round][]*SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			RoundChangeContainer: &MsgContainer{
				Msgs: map[Round][]*SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
		},
	}

	byts, err := i.Encode()
	require.NoError(t, err)

	decoded := &Instance{}
	require.NoError(t, decoded.Decode(byts))

	bytsDecoded, err := decoded.Encode()
	require.NoError(t, err)
	require.EqualValues(t, byts, bytsDecoded)
}
