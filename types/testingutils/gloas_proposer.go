package testingutils

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// Gloas (ePBS) proposer fixtures for the §4 post-consensus-fold envelope path (SIP #94 §4/§6).

// WithGloasEnvelopeBroadcast appends the root of the published §6 reveal to a self-build proposer duty's
// expected beacon broadcast roots (SIP #94 §4/§6): the fixture envelope under the signature the cluster's
// shares reconstruct to.
var WithGloasEnvelopeBroadcast = func(roots []string, ks *TestKeySet, version spec.DataVersion) []string {
	return append(roots, TestingSignedEnvelopeRootOf(ks, TestingGloasProposalData(TestingDutySlotV(version))))
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
	return testingBlindedExecutionPayloadEnvelopeFor(gloas.TestingBeaconBlock(slot))
}

// testingBlindedExecutionPayloadEnvelopeFor derives the §6 blinded envelope of a self-build block carrying
// TestingGloasPayloadRoot.
func testingBlindedExecutionPayloadEnvelopeFor(block *gloas.BeaconBlock) *gloas.BlindedExecutionPayloadEnvelope {
	return TestingBlindedExecutionPayloadEnvelopeOf(&gloas.GloasProposalData{Block: block, PayloadRoot: TestingGloasPayloadRoot})
}

// TestingBlindedExecutionPayloadEnvelopeOf derives the §6 blinded envelope of a self-build proposal, field by
// field as SIP #94 §6 lists them.
var TestingBlindedExecutionPayloadEnvelopeOf = func(proposal *gloas.GloasProposalData) *gloas.BlindedExecutionPayloadEnvelope {
	root, err := proposal.Block.HashTreeRoot()
	if err != nil {
		panic(err.Error())
	}
	return &gloas.BlindedExecutionPayloadEnvelope{
		PayloadRoot:           proposal.PayloadRoot,
		ExecutionRequestsRoot: proposal.Block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot,
		BuilderIndex:          gloas.BuilderIndexSelfBuild,
		BeaconBlockRoot:       root,
		ParentBeaconBlockRoot: proposal.Block.ParentRoot,
	}
}

// TestingSignedBeaconBlockOf is a proposal's block under the signature the cluster's shares reconstruct to, for
// the epoch of the block's slot.
var TestingSignedBeaconBlockOf = func(ks *TestKeySet, proposal *gloas.GloasProposalData) *gloas.SignedBeaconBlock {
	return &gloas.SignedBeaconBlock{
		Message:   proposal.Block,
		Signature: signBeaconObjectAtSlot(proposal.Block, types.DomainProposer, ks, proposal.Block.Slot),
	}
}

// TestingSignedEnvelopeRootOf is the hex root of a self-build proposal's published §6 reveal: its blinded
// envelope under the signature the cluster's shares reconstruct to (see signedEnvelopeRoot).
var TestingSignedEnvelopeRootOf = func(ks *TestKeySet, proposal *gloas.GloasProposalData) string {
	envelope := TestingBlindedExecutionPayloadEnvelopeOf(proposal)
	root, err := signedEnvelopeRoot(envelope, signBeaconObjectAtSlot(envelope, types.DomainBeaconBuilder, ks, proposal.Block.Slot))
	if err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(root[:])
}

// TestingOtherSelfBuildBeaconBlock is a self-build block for the slot that differs from the fixture block
// (gloas.TestingBeaconBlock) in its graffiti only: what an operator whose own produce is not the decided block
// holds, so the envelope it produced is not the decided block's (SIP #94 §6).
var TestingOtherSelfBuildBeaconBlock = func(slot phase0.Slot) *gloas.BeaconBlock {
	block := gloas.TestingBeaconBlock(slot)
	block.Body.Graffiti = [32]byte{0x07, 0x1e, 0x12}
	return block
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
