package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgVoluntaryExit = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleVoluntaryExit))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusVoluntaryExitMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var PreConsensusVoluntaryExitTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusVoluntaryExitWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, true)
}

var PreConsensusVoluntaryExitWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return VoluntaryExitMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var VoluntaryExitMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainVoluntaryExit)

	signed, root, _ := signer.SignBeaconObject(TestingVoluntaryExitBySlot(slot), d,
		beaconSK.GetPublicKey().Serialize(),
		types.DomainVoluntaryExit)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingVoluntaryExitWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainVoluntaryExit)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingVoluntaryExit, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainVoluntaryExit)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.VoluntaryExitPartialSig,
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
