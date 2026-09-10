package ssv_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// TestProposerProcessPreConsensusExternalBid covers the Gloas (ePBS §4) external-build produce path: a
// p2p/builder-API bid win returns a bare block and NO envelope. ProcessPreConsensus must not dereference
// the nil envelope (it did before the fix, panicking) and must decide a value whose payload_root is zero —
// the value check pins payload_root zero iff the bid is not self-build. The randao quorum triggers the
// produce, so this drives three randao partial signatures to threshold and then inspects the value the
// runner started consensus on.
func TestProposerProcessPreConsensusExternalBid(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	version := gloas.DataVersionGloas
	duty := testingutils.TestingProposerDutyV(version)

	r := testingutils.ProposerRunner(ks)
	r.GetBeaconNode().(*testingutils.TestingBeaconNode).SetProducesExternalBid(true)
	r.GetBaseRunner().State = ssv.NewRunnerState(ks.Threshold, duty)

	for i := types.OperatorID(1); i <= types.OperatorID(ks.Threshold); i++ {
		require.NoError(t, r.ProcessPreConsensus(testingutils.PreConsensusRandaoMsgV(ks.Shares[i], i, version)))
	}

	// The randao quorum produced the external-bid block and started consensus on it.
	inst := r.GetBaseRunner().State.RunningInstance
	require.NotNil(t, inst)

	cd := &types.ProposerConsensusData{}
	require.NoError(t, cd.Decode(inst.StartValue))
	require.Equal(t, version, cd.Version)

	pd, err := gloas.DecodeGloasProposalData(cd.DataSSZ)
	require.NoError(t, err)
	require.Equal(t, phase0.Root{}, pd.PayloadRoot) // zero, because the bid is not self-build
	require.Equal(t, gloas.TestingExternalBuilderIndex, pd.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex)
	require.NotEqual(t, gloas.BuilderIndexSelfBuild, pd.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex)
}
