package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgProposer = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleProposer))
}

// ==================================================
// PostConsensus
// ==================================================

var PostConsensusProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusBeaconBlockMsgV(sk, id, false, false, version)
}

// PostConsensusProposerBlockOnlyMsgV keeps only the block entry, dropping the §6 envelope entry — the packet
// a run sees when the envelope partial-sigs miss quorum. The block entry is required, so the duty still
// finalizes (SIP #94 §4/§6).
var PostConsensusProposerBlockOnlyMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	msg.Messages = msg.Messages[:1] // block entry is built first; drop any trailing envelope entry
	return msg
}

// PostConsensusProposerEnvelopeFirstMsgV reverses the packet's entry order — the §6 envelope entry first, the
// block entry second. The receiver classifies entries by matched root, not position, so both must still
// submit; this guards against a future order-dependent rewrite (SIP #94 §4).
var PostConsensusProposerEnvelopeFirstMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	msg.Messages[0], msg.Messages[1] = msg.Messages[1], msg.Messages[0]
	return msg
}

// PostConsensusProposerBadEnvelopeShareMsgV signs the block entry correctly but the §6 envelope entry with a
// different validator key (a structurally valid partial that fails beacon-sig verification), keeping both
// signing roots. The block reconstructs while the envelope fails — pinning that a bad envelope share does not
// strand the block (SIP #94 §4/§6).
var PostConsensusProposerBadEnvelopeShareMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	// Re-sign the envelope entry with a different validator key: same signing root, invalid signature.
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	envelope := TestingBlindedExecutionPayloadEnvelope(TestingDutySlotV(version))
	dBuilder, _ := beacon.DomainData(1, types.DomainBeaconBuilder)
	badSig, envSigningRoot, _ := signer.SignBeaconObject(envelope, dBuilder, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainBeaconBuilder)
	badBls := phase0.BLSSignature{}
	copy(badBls[:], badSig)
	msg.Messages[1].PartialSignature = badBls[:]
	msg.Messages[1].SigningRoot = envSigningRoot
	return msg
}

// PostConsensusProposerBadBlockShareMsgV signs the §6 envelope entry correctly but the block entry with a
// different validator key (a structurally valid partial that fails beacon-sig verification), keeping both
// signing roots. The envelope reconstructs while the block fails — the mirror of the bad-envelope case,
// pinning that a bad block share does not strand the §6 envelope (SIP #94 §4/§6).
var PostConsensusProposerBadBlockShareMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	// Re-sign the block entry (Messages[0]) with a different validator key: same signing root, invalid signature.
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	blockRoot, err := gloas.TestingBeaconBlock(TestingDutySlotV(version)).HashTreeRoot()
	if err != nil {
		panic(err)
	}
	dProposer, _ := beacon.DomainData(1, types.DomainProposer)
	badSig, blockSigningRoot, _ := signer.SignBeaconObject(types.SSZ32Bytes(blockRoot), dProposer, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainProposer)
	badBls := phase0.BLSSignature{}
	copy(badBls[:], badSig)
	msg.Messages[0].PartialSignature = badBls[:]
	msg.Messages[0].SigningRoot = blockSigningRoot
	return msg
}

// PostConsensusProposerBadSharesMsgV signs both entries with a different validator key, keeping both signing
// roots: an operator whose block share and §6 envelope share are both bad. Each root must drop the share and
// still reconstruct from the others (SIP #94 §4/§6).
var PostConsensusProposerBadSharesMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := PostConsensusProposerBadBlockShareMsgV(sk, id, version)
	msg.Messages[1] = PostConsensusProposerBadEnvelopeShareMsgV(sk, id, version).Messages[1]
	return msg
}

// PostConsensusProposerEnvelopeOnlyMsgV keeps only the §6 envelope entry, dropping the block entry. The block
// root must be present in a post-consensus packet, so receivers reject it (SIP #94 §4).
var PostConsensusProposerEnvelopeOnlyMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	msg.Messages = msg.Messages[1:] // the block entry is built first
	return msg
}

// PostConsensusProposerMsgFor is an operator's post-consensus packet for a decided proposal (SIP #94 §4): the
// block entry and, when the bid is self-build, the §6 envelope entry, each under its own domain for the epoch
// of the block's slot. Unlike PostConsensusProposerMsgV it is not tied to the fixture block and slot.
var PostConsensusProposerMsgFor = func(sk *bls.SecretKey, id types.OperatorID, proposal *gloas.GloasProposalData) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(proposal.Block.Slot)
	pk := sk.GetPublicKey().Serialize()

	entry := func(obj types.HashRoot, domainType phase0.DomainType) *types.PartialSignatureMessage {
		d, _ := beacon.DomainData(epoch, domainType)
		sig, signingRoot, err := signer.SignBeaconObject(obj, d, pk, domainType)
		if err != nil {
			panic(err.Error())
		}
		return &types.PartialSignatureMessage{
			PartialSignature: sig,
			SigningRoot:      signingRoot,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
	}

	entries := []*types.PartialSignatureMessage{entry(proposal.Block, types.DomainProposer)}
	if proposal.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex == gloas.BuilderIndexSelfBuild {
		entries = append(entries, entry(TestingBlindedExecutionPayloadEnvelopeOf(proposal), types.DomainBeaconBuilder))
	}
	return &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     proposal.Block.Slot,
		Messages: entries,
	}
}

var PostConsensusProposerTooManyRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	ret := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	ret.Messages = append(ret.Messages, ret.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingProposerDutyV(version).Slot,
		Messages: ret.Messages,
	}
	return msg
}

var PostConsensusProposerTooFewRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingProposerDutyV(version).Slot,
		Messages: []*types.PartialSignatureMessage{},
	}
	return msg
}

var PostConsensusWrongProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusBeaconBlockMsgV(sk, id, true, false, version)
}

var PostConsensusWrongValidatorIndexProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, true, false, version)
	for _, m := range msg.Messages {
		m.ValidatorIndex = TestingWrongValidatorIndex
	}
	return msg
}

var PostConsensusWrongSigProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return postConsensusBeaconBlockMsgV(sk, id, false, true, version)
}

var postConsensusBeaconBlockMsgV = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
	version spec.DataVersion,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()

	// The block root is the bid-only fixture block's own hash tree root (a wrong root comes from a wrong-slot
	// block).
	slot := TestingDutySlotV(version)
	if wrongRoot {
		slot += 100
	}
	blockRoot, err := gloas.TestingBeaconBlock(slot).HashTreeRoot()
	if err != nil {
		panic(err)
	}

	pk := sk.GetPublicKey().Serialize()
	if wrongBeaconSig {
		pk = Testing7SharesSet().ValidatorPK.Serialize()
	}

	// block entry, under DomainProposer
	dProposer, _ := beacon.DomainData(1, types.DomainProposer) // epoch doesn't matter here, hard coded
	blockSig, blockSigningRoot, _ := signer.SignBeaconObject(types.SSZ32Bytes(blockRoot), dProposer, pk, types.DomainProposer)
	blockBls := phase0.BLSSignature{}
	copy(blockBls[:], blockSig)

	entries := []*types.PartialSignatureMessage{{
		PartialSignature: blockBls[:],
		SigningRoot:      blockSigningRoot,
		Signer:           id,
		ValidatorIndex:   TestingValidatorIndex,
	}}

	// Self-build: the §6 blinded-envelope entry under DomainBeaconBuilder rides the same packet (SIP #94 §4/§6).
	// The wrong-root / wrong-sig variants exercise the block entry only.
	if !wrongRoot && !wrongBeaconSig {
		envelope := TestingBlindedExecutionPayloadEnvelope(TestingDutySlotV(version))
		dBuilder, _ := beacon.DomainData(1, types.DomainBeaconBuilder)
		envSig, envSigningRoot, _ := signer.SignBeaconObject(envelope, dBuilder, pk, types.DomainBeaconBuilder)
		envBls := phase0.BLSSignature{}
		copy(envBls[:], envSig)
		entries = append(entries, &types.PartialSignatureMessage{
			PartialSignature: envBls[:],
			SigningRoot:      envSigningRoot,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		})
	}

	return &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingProposerDutyV(version).Slot,
		Messages: entries,
	}
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusRandaoMsg = func(sk *bls.SecretKey, id types.OperatorID) *types.PartialSignatureMessages {
	return randaoMsg(sk, id, false, TestingDutyEpoch, 1, false)
}

var PreConsensusRandaoMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 1, false, version)
}

// PreConsensusRandaoMsgForSlot is an operator's randao partial for a slot's proposer duty: the reveal of the
// slot's epoch, under the domain for that epoch.
var PreConsensusRandaoMsgForSlot = func(sk *bls.SecretKey, id types.OperatorID, slot phase0.Slot) *types.PartialSignatureMessages {
	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(slot)
	d, _ := NewTestingBeaconNode().DomainData(epoch, types.DomainRandao)
	signed, root, err := NewTestingKeyManager().SignBeaconObject(types.SSZUint64(epoch), d, sk.GetPublicKey().Serialize(), types.DomainRandao)
	if err != nil {
		panic(err.Error())
	}
	return &types.PartialSignatureMessages{
		Type: types.RandaoPartialSig,
		Slot: slot,
		Messages: []*types.PartialSignatureMessage{{
			PartialSignature: signed,
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}},
	}
}

var PreConsensusRandaoNextEpochMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version)+1, 1, false, version)
}

var PreConsensusRandaoDifferentEpochMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version)+1, 1, false, version)
}

var PreConsensusRandaoTooManyRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 2, false, version)
}

var PreConsensusRandaoTooFewRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 0, false, version)
}

var PreConsensusRandaoNoMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 0, false, version)
}

var PreConsensusRandaoWrongBeaconSigMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 1, true, version)
}

var randaoMsgV = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	epoch phase0.Epoch,
	msgCnt int,
	wrongBeaconSig bool,
	version spec.DataVersion,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(epoch, types.DomainRandao)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(epoch), d, sk.GetPublicKey().Serialize(), types.DomainRandao)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(types.SSZUint64(TestingDutyEpochV(version)), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainRandao)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     TestingProposerDutyV(version).Slot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
		if wrongRoot {
			msg.SigningRoot = [32]byte{}
		}
		msgs.Messages = append(msgs.Messages, msg)
	}
	return &msgs
}

var randaoMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	epoch phase0.Epoch,
	msgCnt int,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(epoch, types.DomainRandao)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(epoch), d, sk.GetPublicKey().Serialize(), types.DomainRandao)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(types.SSZUint64(TestingDutyEpoch), d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainRandao)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.RandaoPartialSig,
		Slot:     TestingDutySlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
		if wrongRoot {
			msg.SigningRoot = [32]byte{}
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	return &msgs
}

var PreConsensusRandaoDifferentSignerMsgV = func(
	msgSigner, randaoSigner *bls.SecretKey,
	msgSignerID,
	randaoSignerID types.OperatorID,
	version spec.DataVersion,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	epoch := TestingDutyEpochV(version)
	d, _ := beacon.DomainData(epoch, types.DomainRandao)
	signed, root, _ := signer.SignBeaconObject(types.SSZUint64(epoch), d, randaoSigner.GetPublicKey().Serialize(), types.DomainRandao)

	msg := types.PartialSignatureMessages{
		Type: types.RandaoPartialSig,
		Slot: TestingProposerDutyV(version).Slot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: signed[:],
				SigningRoot:      root,
				Signer:           randaoSignerID,
				ValidatorIndex:   TestingValidatorIndex,
			},
		},
	}
	return &msg
}
