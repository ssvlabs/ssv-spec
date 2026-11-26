package ssv_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestProposerPostConsensusRootsFollowDecidedValue pins that the proposer derives its post-consensus roots again
// for each duty's decided value. What it caches between packets is not part of the runner's JSON, and a vector
// runs a single duty on a runner that starts without it, so no vector can cover a second duty.
func TestProposerPostConsensusRootsFollowDecidedValue(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	runner := testingutils.ProposerRunner(ks)
	validator := testingutils.BaseValidator(ks)
	validator.DutyRunners[types.RoleProposer] = runner
	validator.Network = runner.GetNetwork()

	quorum := []types.OperatorID{1, 2, 3}
	var expectedRoots []string
	for _, slot := range []phase0.Slot{testingutils.TestingDutySlotGloas, testingutils.TestingDutySlotGloasNextEpoch} {
		duty := testingutils.TestingProposerDutyForSlot(slot)
		proposal := testingutils.TestingGloasProposalData(slot)
		consensusData := &types.ProposerConsensusData{
			Duty:    *duty,
			Version: gloas.DataVersionGloas,
			DataSSZ: testingutils.TestingGloasProposalDataBytes(slot),
		}

		var msgs []*types.SignedSSVMessage
		for _, opID := range quorum {
			msgs = append(msgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoMsgForSlot(ks.Shares[opID], opID, slot))))
		}
		msgs = append(msgs, testingutils.SSVDecidingMsgsForHeight(consensusData, testingutils.ProposerMsgID, qbft.Height(slot), ks)...)
		for _, opID := range quorum {
			msgs = append(msgs, testingutils.SignPartialSigSSVMessage(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgFor(ks.Shares[opID], opID, proposal))))
		}

		require.NoError(t, validator.StartDuty(duty))
		for _, msg := range msgs {
			require.NoError(t, validator.ProcessMessage(msg))
		}
		require.False(t, runner.HasRunningDuty())

		expectedRoots = append(expectedRoots,
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockOf(ks, proposal)),
			testingutils.TestingSignedEnvelopeRootOf(ks, proposal),
		)
	}

	// Each duty submitted its own block and published its own reveal.
	beaconNode := runner.GetBeaconNode().(*testingutils.TestingBeaconNode)
	testingutils.CompareBroadcastedBeaconMsgsInOrder(t, expectedRoots, beaconNode.BroadcastedRoots)
}
