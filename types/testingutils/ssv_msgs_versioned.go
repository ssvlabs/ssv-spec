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
		justif = append(justif, PreConsensusRandaoMsg(ks.Shares[i+1], i+1))
	}

	cd := TestProposerConsensusDataV(version)
	cd.PreConsensusJustifications = justif
	return cd
}

var TestProposerBlindedWithJustificationsConsensusDataV = func(ks *TestKeySet, version spec.DataVersion) *types.ConsensusData {
	justif := make([]*types.SignedPartialSignatureMessage, 0)
	for i := uint64(0); i <= ks.Threshold; i++ {
		justif = append(justif, PreConsensusRandaoMsg(ks.Shares[i+1], i+1))
	}

	cd := TestProposerBlindedBlockConsensusDataV(version)
	cd.PreConsensusJustifications = justif
	return cd
}

var TestProposerBlindedBlockConsensusDataV = func(version spec.DataVersion) *types.ConsensusData {
	return &types.ConsensusData{
		Duty:    TestingProposerDuty,
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
		Slot:     TestingDutySlot,
		Messages: ret.Message.Messages,
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

	block := TestingBeaconBlockV(version)
	if wrongRoot {
		block = TestingWrongBeaconBlockV(version)
	}

	d, _ := beacon.DomainData(1, types.DomainProposer) // epoch doesn't matter here, hard coded
	sig, root, _ := signer.SignBeaconObject(block, d, sk.GetPublicKey().Serialize(), types.DomainProposer)
	if wrongBeaconSig {
		sig, root, _ = signer.SignBeaconObject(block, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainProposer)
	}
	blsSig := phase0.BLSSignature{}
	copy(blsSig[:], sig)

	msgs := types.PartialSignatureMessages{
		Type: types.PostConsensusPartialSig,
		Slot: TestingDutySlot,
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
