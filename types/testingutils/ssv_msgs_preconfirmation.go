package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgPreconf = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RolePreconfirmation))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusPreconfMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return PreconfMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusPreconfNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return PreconfMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var PreConsensusPreconfTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return PreconfMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusPreconfTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return PreconfMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusPreconfWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return PreconfMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, true)
}

var PreConsensusPreconfWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return PreconfMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var PreconfMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainCommitBoost)

	signed, root, _ := signer.SignBeaconObject(TestingPreconfRequest, d,
		beaconSK.GetPublicKey().Serialize(),
		types.DomainCommitBoost)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingPreconfRequestWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainCommitBoost)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingPreconfRequest, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainCommitBoost)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.PreconfPartialSig,
		Slot:     slot,
		Messages: []*types.PartialSignatureMessage{},
	}

	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           beaconID,
			ValidatorIndex:   TestingValidatorIndex,
		}
		msgs.Messages = append(msgs.Messages, msg)
	}

	msg := &types.PartialSignatureMessage{
		PartialSignature: signed[:],
		SigningRoot:      root,
		Signer:           id,
		ValidatorIndex:   TestingValidatorIndex,
	}
	if wrongRoot {
		msg.SigningRoot = [32]byte{}
	}

	return &msgs
}
