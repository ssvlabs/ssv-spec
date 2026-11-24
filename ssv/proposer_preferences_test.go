package ssv_test

import (
	"errors"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

var errRequestAuthBroadcast = errors.New("forced request-auth broadcast failure")

// requestAuthFailingNetwork broadcasts every message except a RequestAuthPartialSig container, which fails.
type requestAuthFailingNetwork struct {
	*testingutils.TestingNetwork
}

func (n *requestAuthFailingNetwork) Broadcast(msgID types.MessageID, message *types.SignedSSVMessage) error {
	msgs := &types.PartialSignatureMessages{}
	if err := msgs.Decode(message.SSVMessage.Data); err == nil && msgs.Type == types.RequestAuthPartialSig {
		return errRequestAuthBroadcast
	}
	return n.TestingNetwork.Broadcast(msgID, message)
}

// TestProposerPreferencesRequestAuthBroadcastFailureReturned pins that a failure to broadcast the
// builder-request-auth container is returned from the duty start without undoing the independent preference
// round (SIP #94 §5): the preference partial is still broadcast, and the slot's sub-runner keeps running with
// both rounds frozen. No vector can cover it, as the testing network never fails.
func TestProposerPreferencesRequestAuthBroadcastFailureReturned(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	network := &requestAuthFailingNetwork{TestingNetwork: testingutils.NewTestingNetwork(1, ks.OperatorKeys[1])}

	runner, err := ssv.NewProposerPreferencesRunner(
		types.BeaconTestNetwork,
		map[phase0.ValidatorIndex]*types.Share{
			testingutils.TestingValidatorIndex: testingutils.TestingShare(ks, testingutils.TestingValidatorIndex),
		},
		testingutils.NewTestingBeaconNode(),
		network,
		testingutils.NewTestingKeyManager(),
		testingutils.TestingOperatorSigner(ks),
		types.DefaultGasLimit,
		testingutils.TestingBuilderEntries,
	)
	require.NoError(t, err)

	duty := testingutils.TestingProposerPreferencesDuty()
	require.ErrorIs(t, runner.StartNewDuty(duty, ks.Threshold), errRequestAuthBroadcast)

	// The preference round went out, and the duty runs with both rounds frozen.
	require.Len(t, network.BroadcastedMsgs, 1)
	broadcast := &types.PartialSignatureMessages{}
	require.NoError(t, broadcast.Decode(network.BroadcastedMsgs[0].SSVMessage.Data))
	require.Equal(t, types.ProposerPreferencesPartialSig, broadcast.Type)
	require.True(t, runner.HasRunningDuty())

	sub := runner.(*ssv.ProposerPreferencesRunner).BySlot[duty.Slot]
	require.NotNil(t, sub.ProposerPreferences)
	require.Len(t, sub.BuilderRequestAuths, len(testingutils.TestingBuilderEntries))
}
