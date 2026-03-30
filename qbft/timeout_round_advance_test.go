package qbft_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

type failingNetwork struct {
	err      error
	attempts int
}

func (n *failingNetwork) Broadcast(_ types.MessageID, _ *types.SignedSSVMessage) error {
	n.attempts++
	return n.err
}

func newTestInstanceWithFailingNetwork(t *testing.T) (*qbft.Instance, *testingutils.TestQBFTTimer, *failingNetwork, *testingutils.TestKeySet) {
	t.Helper()

	ks := testingutils.Testing4SharesSet()
	config := testingutils.TestingConfig(ks)
	timer, ok := config.Timer.(*testingutils.TestQBFTTimer)
	require.True(t, ok)

	network := &failingNetwork{err: errors.New("network down")}
	config.Network = network

	instance := qbft.NewInstance(
		config,
		testingutils.TestingCommitteeMember(ks),
		testingutils.TestingIdentifier,
		qbft.FirstHeight,
		testingutils.TestingOperatorSigner(ks),
	)
	instance.StartValue = testingutils.TestingQBFTFullData

	return instance, timer, network, ks
}

func TestUponRoundTimeoutDoesNotAdvanceRoundWhenBroadcastFails(t *testing.T) {
	instance, timer, network, ks := newTestInstanceWithFailingNetwork(t)
	originalProposal := testingutils.ToProcessingMessage(
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	)
	instance.State.ProposalAcceptedForCurrentRound = originalProposal

	err := instance.UponRoundTimeout()
	require.ErrorContains(t, err, "failed to broadcast round change message")

	require.Equal(t, qbft.FirstRound, instance.State.Round)
	require.Same(t, originalProposal, instance.State.ProposalAcceptedForCurrentRound)
	require.Equal(t, 0, timer.State.Timeouts)
	require.Equal(t, qbft.NoRound, timer.State.Round)
	require.Equal(t, 1, network.attempts)
}

func TestProcessMsgDoesNotAdvanceRoundOnPartialQuorumWhenBroadcastFails(t *testing.T) {
	instance, timer, network, ks := newTestInstanceWithFailingNetwork(t)
	originalProposal := testingutils.ToProcessingMessage(
		testingutils.TestingProposalMessage(ks.OperatorKeys[1], types.OperatorID(1)),
	)
	instance.State.ProposalAcceptedForCurrentRound = originalProposal

	msg1, err := qbft.NewProcessingMessage(
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[2], types.OperatorID(2), qbft.FirstRound+1),
	)
	require.NoError(t, err)
	_, _, _, err = instance.ProcessMsg(msg1)
	require.NoError(t, err)

	msg2, err := qbft.NewProcessingMessage(
		testingutils.TestingRoundChangeMessageWithRound(ks.OperatorKeys[3], types.OperatorID(3), qbft.FirstRound+1),
	)
	require.NoError(t, err)
	_, _, _, err = instance.ProcessMsg(msg2)
	require.ErrorContains(t, err, "failed to broadcast round change message")

	require.Equal(t, qbft.FirstRound, instance.State.Round)
	require.Same(t, originalProposal, instance.State.ProposalAcceptedForCurrentRound)
	require.Equal(t, 0, timer.State.Timeouts)
	require.Equal(t, qbft.NoRound, timer.State.Round)
	require.Equal(t, 1, network.attempts)
}
