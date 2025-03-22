package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgCBSigning = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	msgID := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleCBSigning)

	if qbftMsg != nil {
		return &types.SSVMessage{
			MsgType: qbftMsg.SSVMessage.MsgType,
			MsgID:   msgID,
			Data:    qbftMsg.SSVMessage.Data,
		}
	}

	if partialSigMsg != nil {
		CBPreConsensusMsg := &types.CBPartialSignatures{
			RequestRoot: TestingCBSigningRequest.Root,
			PartialSig:  *partialSigMsg,
		}

		msgType := types.CommitBoostPartialSignatureMsgType
		data, err := CBPreConsensusMsg.Encode()
		if err != nil {
			panic(err)
		}
		return &types.SSVMessage{
			MsgType: msgType,
			MsgID:   msgID,
			Data:    data,
		}
	}

	panic("msg type undefined")
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusCBSigningMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return CBSigningMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusCBSigningNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return CBSigningMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var PreConsensusCBSigningTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return CBSigningMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusCBSigningTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return CBSigningMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusCBSigningWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return CBSigningMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, true)
}

var PreConsensusCBSigningWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return CBSigningMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var CBSigningMsg = func(
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

	signed, root, _ := signer.SignBeaconObject(TestingCBSigningRequest, d,
		beaconSK.GetPublicKey().Serialize(),
		types.DomainCommitBoost)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingCBSigningRequestWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainCommitBoost)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingCBSigningRequest, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainCommitBoost)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.CBSigningPartialSig,
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

	return &msgs
}
