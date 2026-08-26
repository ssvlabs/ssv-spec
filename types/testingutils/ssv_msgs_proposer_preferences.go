package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec/types"
)

// ==================================================
// SSVMessage
// ==================================================

var SSVMsgProposerPreferences = func(qbftMsg *types.SignedSSVMessage, partialSigMsg *types.PartialSignatureMessages) *types.SSVMessage {
	return ssvMsg(qbftMsg, partialSigMsg, types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), types.RoleProposerPreferences))
}

// ==================================================
// PreConsensus
// ==================================================

var PreConsensusProposerPreferencesMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 1, false, TestingDutySlotGloas, false)
}

var PreConsensusProposerPreferencesTooFewRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 0, false, TestingDutySlotGloas, false)
}

var PreConsensusProposerPreferencesTooManyRootsMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 2, false, TestingDutySlotGloas, false)
}

var PreConsensusProposerPreferencesWrongBeaconSigMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 1, false, TestingDutySlotGloas, true)
}

var PreConsensusProposerPreferencesNextEpochMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 1, false, TestingDutySlotGloasNextEpoch, false)
}

// PreConsensusProposerPreferencesSecondSlotMsg targets the second concurrently-active lookahead
// proposal slot (TestingProposerPreferencesSecondDuty).
var PreConsensusProposerPreferencesSecondSlotMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 1, false, TestingDutySlotGloas+1, false)
}

// PreConsensusProposerPreferencesWrongRootMsg carries a preference derived from a diverging
// dependent root: a valid message from a peer whose beacon node observed a different duty
// assignment, so it must fail the operator's expected-root check rather than mix into a quorum.
var PreConsensusProposerPreferencesWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID) *types.PartialSignatureMessages {
	return proposerPreferencesMsg(msgSK, msgID, 1, true, TestingDutySlotGloas, false)
}

var proposerPreferencesMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	msgCnt int,
	divergingDependentRoot bool,
	proposalSlot phase0.Slot,
	wrongBeaconSig bool,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(proposalSlot)
	d, _ := beacon.DomainData(epoch, types.DomainProposerPreferences)

	preferences := TestingProposerPreferences(proposalSlot)
	if divergingDependentRoot {
		preferences.DependentRoot = TestingBlockRoot
	}
	signed, root, _ := signer.SignBeaconObject(preferences, d, sk.GetPublicKey().Serialize(), types.DomainProposerPreferences)
	if wrongBeaconSig {
		signed, root, _ = signer.SignBeaconObject(preferences, d, Testing7SharesSet().ValidatorPK.Serialize(), types.DomainProposerPreferences)
	}

	msgs := types.PartialSignatureMessages{
		Type:     types.ProposerPreferencesPartialSig,
		Slot:     proposalSlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for i := 0; i < msgCnt; i++ {
		msg := &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		}
		msgs.Messages = append(msgs.Messages, msg)
	}
	return &msgs
}
