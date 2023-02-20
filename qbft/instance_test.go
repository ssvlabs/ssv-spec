package qbft_test

import (
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
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
	TestingSK := func() *bls.SecretKey {
		types.InitBLS()
		ret := &bls.SecretKey{}
		ret.SetByCSPRNG()
		return ret
	}()
	testingSignedMsg := func() *qbft.SignedMessage {
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
	i := &qbft.Instance{
		State: &qbft.State{
			Share:                           testingShare,
			ID:                              []byte{1, 2, 3, 4},
			Round:                           1,
			Height:                          1,
			LastPreparedRound:               1,
			LastPreparedValue:               []byte{1, 2, 3, 4},
			ProposalAcceptedForCurrentRound: testingSignedMsg,
			Decided:                         false,
			DecidedValue:                    []byte{1, 2, 3, 4},

			ProposeContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			PrepareContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			CommitContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.SignedMessage{
					1: {
						testingSignedMsg,
					},
				},
			},
			RoundChangeContainer: &qbft.MsgContainer{
				Msgs: map[qbft.Round][]*qbft.SignedMessage{
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
