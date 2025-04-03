package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgValidatorRegistration = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], types.RoleValidatorRegistration))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusValidatorRegistrationMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 0, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 2, false, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot, true)
}

var PreConsensusValidatorRegistrationWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, true, TestingDutySlot, false)
}

var PreConsensusValidatorRegistrationNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return validatorRegistrationMsg(msgSK, msgSK, msgID, msgID, 1, false, TestingDutySlot2, false)
}

var validatorRegistrationMsg = func(
	sk, beaconSK *bls.SecretKey,
	id, beaconID types.OperatorID,
	msgCnt int,
	wrongRoot bool,
	slot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	d, _ := beacon.DomainData(TestingDutyEpoch, types.DomainApplicationBuilder)

	signed, root, _ := signer.SignBeaconObject(TestingValidatorRegistrationBySlot(slot), d,
		beaconSK.GetPublicKey().Serialize(),
		types.DomainApplicationBuilder)
	if wrongRoot {
		signed, root, _ = signer.SignBeaconObject(TestingValidatorRegistrationWrong, d, beaconSK.GetPublicKey().Serialize(), types.DomainApplicationBuilder)
	}
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(TestingValidatorRegistration, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainApplicationBuilder)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.ValidatorRegistrationPartialSig,
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
