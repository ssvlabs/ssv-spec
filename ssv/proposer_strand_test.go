package ssv_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestProposerBadEnvelopeShareDoesNotStrandBlock pins SIP #94 §4/§6 per-root isolation on the Gloas proposer
// post-consensus: a self-build packet's block and §6 envelope roots reconstruct independently, so an invalid
// envelope share — which crosses quorum-count then fails reconstruction — must not strand the required block.
// It's a Go unit test, not a spec vector, because the reconstruction failure is an untyped BLS error the
// MsgProcessingSpecTest ExpectedErrorCode framework can't match.
func TestProposerBadEnvelopeShareDoesNotStrandBlock(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas
	duty := testingutils.TestingProposerDutyV(version)
	cd := testingutils.TestProposerConsensusDataV(version)

	r := testingutils.ProposerRunner(ks)
	bn := r.GetBeaconNode().(*testingutils.TestingBeaconNode)

	// Bring the runner to a decided self-build value (block + §6 envelope roots) — the state a post-consensus
	// round runs against. producedEnvelope is not needed: the envelope reconstruction fails before the reveal
	// would be built.
	cdBytes, err := cd.Encode()
	require.NoError(t, err)
	base := r.GetBaseRunner()
	base.State = ssv.NewRunnerState(ks.Threshold, duty)
	inst := qbft.NewInstance(base.QBFTController.GetConfig(), base.QBFTController.CommitteeMember,
		base.QBFTController.Identifier, qbft.Height(duty.Slot), base.QBFTController.OperatorSigner)
	inst.State.Decided = true
	inst.State.DecidedValue = cdBytes
	base.State.RunningInstance = inst
	base.State.DecidedValue = cdBytes
	base.QBFTController.StoredInstances = append(base.QBFTController.StoredInstances, inst)
	base.QBFTController.Height = qbft.Height(duty.Slot)

	// Two valid full packets (block + envelope), below the quorum of three.
	require.NoError(t, r.ProcessPostConsensus(testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)))
	require.NoError(t, r.ProcessPostConsensus(testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)))
	require.Empty(t, bn.BroadcastedRoots, "nothing submitted before quorum")

	// op3 crosses both roots to quorum with an invalid envelope share, and lists the envelope entry first, so
	// a packet-order-iteration regression would process (and fail on) the envelope before submitting the
	// block. The block must still submit; the envelope fails reconstruction independently.
	badPkt := testingutils.PostConsensusProposerBadEnvelopeShareMsgV(ks.Shares[3], 3, version)
	badPkt.Messages[0], badPkt.Messages[1] = badPkt.Messages[1], badPkt.Messages[0] // envelope entry first
	err = r.ProcessPostConsensus(badPkt)
	require.Error(t, err, "envelope reconstruction fails on the bad share")
	require.Contains(t, err.Error(), "invalid signatures")

	// Exactly the block was submitted; the envelope reveal was not published.
	blockRoot := testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version))
	require.Len(t, bn.BroadcastedRoots, 1, "only the block is submitted; the envelope reveal is not")
	require.Equal(t, blockRoot, hex.EncodeToString(bn.BroadcastedRoots[0][:]))
}

// TestProposerFailedBlockSubmitStillPublishesReveal pins SIP #94 §6 "retry until they have": the §6 reveal is
// gated on the block submit being attempted, not on it succeeding, so when this operator's own block submit
// fails — the block still reaches beacon nodes from the other operators' submits, and a node ignores an
// envelope whose block it has not seen — the reveal must still publish. It runs the full self-build produce +
// decide flow (the reveal needs producedEnvelope) with the beacon node set to fail the block submit, and is a
// Go unit test because the failed submit surfaces as an untyped error.
func TestProposerFailedBlockSubmitStillPublishesReveal(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas
	duty := testingutils.TestingProposerDutyV(version)

	r := testingutils.ProposerRunner(ks)
	bn := r.GetBeaconNode().(*testingutils.TestingBeaconNode)
	base := r.GetBaseRunner()
	base.State = ssv.NewRunnerState(ks.Threshold, duty)

	// Self-build produce: the randao quorum yields the block and this operator's §6 envelope (producedEnvelope)
	// and starts consensus.
	for i := types.OperatorID(1); i <= types.OperatorID(ks.Threshold); i++ {
		require.NoError(t, r.ProcessPreConsensus(testingutils.PreConsensusRandaoMsgV(ks.Shares[i], i, version)))
	}
	inst := base.State.RunningInstance
	require.NotNil(t, inst)

	// Decide the produced value.
	cdBytes, err := testingutils.TestProposerConsensusDataV(version).Encode()
	require.NoError(t, err)
	inst.State.Decided = true
	inst.State.DecidedValue = cdBytes
	base.State.DecidedValue = cdBytes

	// This operator's own block submit fails.
	bn.SetFailGloasBlockSubmit(true)

	require.NoError(t, r.ProcessPostConsensus(testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, version)))
	require.NoError(t, r.ProcessPostConsensus(testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, version)))
	err = r.ProcessPostConsensus(testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, version))
	require.Error(t, err, "the local block submit fails")
	require.Contains(t, err.Error(), "could not submit")

	// The §6 reveal published despite the failed block submit; the failed block itself was not recorded.
	blockRoot := testingutils.GetSSZRootNoError(testingutils.TestingSignedBeaconBlockV(ks, version))
	envelopeRoot := testingutils.GetSSZRootNoError(testingutils.TestingBlindedExecutionPayloadEnvelope(testingutils.TestingDutySlotV(version)))
	var broadcast []string
	for _, rr := range bn.BroadcastedRoots {
		broadcast = append(broadcast, hex.EncodeToString(rr[:]))
	}
	require.Contains(t, broadcast, envelopeRoot, "the §6 reveal must publish even when the local block submit failed")
	require.NotContains(t, broadcast, blockRoot, "the failed block submit is not recorded")
}
