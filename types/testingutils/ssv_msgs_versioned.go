package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/types"
)

var TestProposerConsensusDataV = func(version spec.DataVersion) *types.ConsensusData {
	duty := TestingProposerDutyV(version)
	return &types.ConsensusData{
		Duty:    *duty,
		Version: version,
		DataSSZ: TestingBeaconBlockBytesV(version),
	}
}

var TestProposerConsensusDataBytsV = func(version spec.DataVersion) []byte {
	cd := TestProposerConsensusDataV(version)
	byts, _ := cd.Encode()
	return byts
}

var TestProposerWithJustificationsConsensusDataV = func(ks *TestKeySet, version spec.DataVersion) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusRandaoMsgV(ks.Shares[i+1], i+1, version))
	}

	cd := TestProposerConsensusDataV(version)
	cd.PreConsensusJustifications = justif
	return cd
}

var TestProposerBlindedWithJustificationsConsensusDataV = func(ks *TestKeySet, version spec.DataVersion) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusRandaoMsgV(ks.Shares[i+1], i+1, version))
	}

	cd := TestProposerBlindedBlockConsensusDataV(version)
	cd.PreConsensusJustifications = justif
	return cd
}

var TestProposerBlindedBlockConsensusDataV = func(version spec.DataVersion) *types.ConsensusData {
	return &types.ConsensusData{
		Duty:    *TestingProposerDutyV(version),
		Version: version,
		DataSSZ: TestingBlindedBeaconBlockBytesV(version),
	}
}

var TestProposerBlindedBlockConsensusDataBytsV = func(version spec.DataVersion) []byte {
	cd := TestProposerBlindedBlockConsensusDataV(version)
	byts, _ := cd.Encode()
	return byts
}

var PostConsensusProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsgV(sk, id, false, false, version)
}

var PostConsensusProposerTooManyRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	ret := postConsensusBeaconBlockMsgV(sk, id, false, false, version)
	ret.Message.Messages = append(ret.Message.Messages, ret.Message.Messages[0])

	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingProposerDutyV(version).Slot,
		Messages: ret.Message.Messages,
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusProposerTooFewRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	msg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     TestingProposerDutyV(version).Slot,
		Messages: []*types.PartialSignatureMessage{},
	}

	sig, _ := NewTestingKeyManager().SignRoot(msg, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: sig,
		Signer:    id,
	}
}

var PostConsensusWrongProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsgV(sk, id, true, false, version)
}

var PostConsensusWrongSigProposerMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return postConsensusBeaconBlockMsgV(sk, id, false, true, version)
}

var PostConsensusSigProposerWrongBeaconSignerMsgV = func(sk *bls.SecretKey, id, beaconSigner types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	ret := postConsensusBeaconBlockMsgV(sk, beaconSigner, false, true, version)
	ret.Signer = id
	return ret
}

var postConsensusBeaconBlockMsgV = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	wrongBeaconSig bool,
	version spec.DataVersion,
) *types.SignedPartialSignatureMessage {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()

	var root phase0.Root
	var err error
	if wrongRoot {
		blk := TestingWrongBeaconBlockV(version)
		root, err = blk.Root()
	} else {
		blk := TestingBeaconBlockV(version)
		root, err = blk.Root()
	}
	if err != nil {
		panic(err)
	}
	hashRoot := types.SSZ32Bytes(root)

	d, _ := beacon.DomainData(1, types.DomainProposer) // epoch doesn't matter here, hard coded
	sig, root, _ := signer.SignBeaconObject(hashRoot, d, sk.GetPublicKey().Serialize(), types.DomainProposer)
	if wrongBeaconSig {
		sig, root, _ = signer.SignBeaconObject(hashRoot, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainProposer)
	}
	blsSig := phase0.BLSSignature{}
	copy(blsSig[:], sig)

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingProposerDutyV(version).Slot,
		Messages: []*types.PartialSignatureMessage{
			{
				PartialSignature: blsSig[:],
				SigningRoot:      root,
				Signer:           id,
			},
		},
	}
	msgSig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: msgSig,
		Signer:    id,
	}
}

var PreConsensusRandaoMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 1, false, version)
}

// PreConsensusRandaoNextEpochMsgV testing for a second duty start
var PreConsensusRandaoNextEpochMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version)+1, 1, false, version)
}

var PreConsensusRandaoDifferentEpochMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version)+1, 1, false, version)
}

var PreConsensusRandaoTooManyRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 2, false, version)
}

var PreConsensusRandaoTooFewRootsMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 0, false, version)
}

var PreConsensusRandaoNoMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 0, false, version)
}

var PreConsensusRandaoWrongBeaconSigMsgV = func(sk *bls.SecretKey, id types.OperatorID, version spec.DataVersion) *types.SignedPartialSignatureMessage {
	return randaoMsgV(sk, id, false, TestingDutyEpochV(version), 1, true, version)
}

var PreConsensusRandaoDifferentSignerMsgV = func(
	msgSigner, randaoSigner *bls.SecretKey,
	msgSignerID,
	randaoSignerID types.OperatorID,
	version spec.DataVersion,
) *types.SignedPartialSignatureMessage {
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
			},
		},
	}
	sig, _ := signer.SignRoot(msg, types.PartialSignatureType, msgSigner.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msg,
		Signature: sig,
		Signer:    msgSignerID,
	}
}

var randaoMsgV = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	wrongRoot bool,
	epoch phase0.Epoch,
	msgCnt int,
	wrongBeaconSig bool,
	version spec.DataVersion,
) *types.SignedPartialSignatureMessage {
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
		}
		if wrongRoot {
			msg.SigningRoot = [32]byte{}
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	sig, _ := signer.SignRoot(msgs, types.PartialSignatureType, sk.GetPublicKey().Serialize())
	return &types.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: sig,
		Signer:    id,
	}
}
