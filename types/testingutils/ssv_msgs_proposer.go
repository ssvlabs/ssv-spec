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

// PostConsensusProposerBlockOnlyMsgV keeps only the block entry, dropping the Gloas §6 envelope entry —
// the packet a run sees when the envelope partial-sigs miss quorum. The block entry is required, so the
// duty still finalizes (SIP #94 §4/§6). A no-op before Gloas, where the packet is block-only already.
var PostConsensusProposerBlockOnlyMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	msg.Messages = msg.Messages[:1] // block entry is built first; drop any trailing envelope entry
	return msg
}

// PostConsensusProposerEnvelopeFirstMsgV reverses the Gloas packet's entry order — the §6 envelope entry
// first, the block entry second. The receiver classifies entries by matched root, not position, so both
// must still submit; this guards against a future order-dependent rewrite (SIP #94 §4). A no-op before
// Gloas, where the packet is block-only.
var PostConsensusProposerEnvelopeFirstMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	if len(msg.Messages) == 2 {
		msg.Messages[0], msg.Messages[1] = msg.Messages[1], msg.Messages[0]
	}
	return msg
}

// PostConsensusProposerBadEnvelopeShareMsgV signs the block entry correctly but the §6 envelope entry with a
// different validator key (a structurally valid partial that fails beacon-sig verification), keeping both
// signing roots. The block reconstructs while the envelope fails — pinning that a bad envelope share does not
// strand the block (SIP #94 §4/§6). A no-op (block-only) before Gloas.
var PostConsensusProposerBadEnvelopeShareMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	if len(msg.Messages) < 2 {
		return msg // pre-Gloas: block-only, no envelope entry to corrupt
	}
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
// pinning that a bad block share does not strand the §6 envelope (SIP #94 §4/§6). A no-op (block-only) before Gloas.
var PostConsensusProposerBadBlockShareMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.PartialSignatureMessages {
	msg := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	if len(msg.Messages) < 2 {
		return msg // pre-Gloas: block-only, no envelope to keep alive
	}
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

	var blockRoot phase0.Root
	var err error
	if version == gloas.DataVersionGloas {
		// Gloas (ePBS §4): the block root is the bid-only fixture block's own hash tree root (a wrong root
		// comes from a wrong-slot block).
		slot := TestingDutySlotV(version)
		if wrongRoot {
			slot += 100
		}
		blockRoot, err = gloas.TestingBeaconBlock(slot).HashTreeRoot()
	} else if wrongRoot {
		blockRoot, err = TestingWrongBeaconBlockV(version).Root()
	} else {
		blockRoot, err = TestingBeaconBlockV(version).Root()
	}
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

	// Gloas self-build: the §6 blinded-envelope entry under DomainBeaconBuilder rides the same packet
	// (SIP #94 §4/§6). The wrong-root / wrong-sig variants exercise the block entry only.
	if version == gloas.DataVersionGloas && !wrongRoot && !wrongBeaconSig {
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
