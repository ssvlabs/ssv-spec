package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// Gloas (ePBS) proposer fixtures for the §4 post-consensus-fold envelope path (SIP #94 §4/§6).

// WithGloasEnvelopeBroadcast appends the §6 envelope publish root to a self-build proposer duty's expected
// beacon broadcast roots (SIP #94 §4/§6).
var WithGloasEnvelopeBroadcast = func(roots []string, version spec.DataVersion) []string {
	return append(roots, GetSSZRootNoError(TestingBlindedExecutionPayloadEnvelope(TestingDutySlotV(version))))
}

// TestingGloasPayloadRoot is the self-build payload_root the fixture proposer carries in its decided value
// (hash_tree_root(envelope.payload) node-side). Non-zero, so it satisfies the §4 self-build value check.
var TestingGloasPayloadRoot = phase0.Root{0x50, 0x51, 0x52}

// TestingGloasProposalData is the decided value's DataSSZ payload: the fixture block plus
// TestingGloasPayloadRoot (SIP #94 §4).
var TestingGloasProposalData = func(slot phase0.Slot) *gloas.GloasProposalData {
	return &gloas.GloasProposalData{
		Block:       gloas.TestingBeaconBlock(slot),
		PayloadRoot: TestingGloasPayloadRoot,
	}
}

var TestingGloasProposalDataBytes = func(slot phase0.Slot) []byte {
	byts, err := TestingGloasProposalData(slot).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	return byts
}

// TestingBlindedExecutionPayloadEnvelope is the §6 blinded envelope the cluster threshold-signs, derived
// from the decided value exactly as the runner derives it (SIP #94 §6): payload_root from the decided
// value, execution_requests_root from the bid, the block's root and parent root.
var TestingBlindedExecutionPayloadEnvelope = func(slot phase0.Slot) *gloas.BlindedExecutionPayloadEnvelope {
	block := gloas.TestingBeaconBlock(slot)
	root, err := block.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}
	return &gloas.BlindedExecutionPayloadEnvelope{
		PayloadRoot:           TestingGloasPayloadRoot,
		ExecutionRequestsRoot: block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot,
		BuilderIndex:          gloas.BuilderIndexSelfBuild,
		BeaconBlockRoot:       root,
		ParentBeaconBlockRoot: block.ParentRoot,
	}
}

// TestProposerExternalBidConsensusDataV is the decided value of an external bid win (SIP #94 §4): the bare
// external-bid block with a zero payload_root, so the post-consensus packet carries the block entry only.
var TestProposerExternalBidConsensusDataV = func(version spec.DataVersion) *types.ProposerConsensusData {
	byts, err := (&gloas.GloasProposalData{Block: gloas.TestingBeaconBlockExternalBuild(TestingDutySlotV(version))}).MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}
	return &types.ProposerConsensusData{
		Duty:    *TestingProposerDutyV(version),
		Version: version,
		DataSSZ: byts,
	}
}
