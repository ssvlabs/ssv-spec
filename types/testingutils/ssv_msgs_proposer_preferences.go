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

// ==================================================
// Builder Request Auth (SIP #94 §5)
// ==================================================

// PreConsensusBuilderRequestAuthMsg is one operator's RequestAuthPartialSig container: a partial signature
// per configured entry data, all riding one message (the multi-root pre-consensus shape).
var PreConsensusBuilderRequestAuthMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID, dataList [][]byte) *types.PartialSignatureMessages {
	return builderRequestAuthMsg(msgSK, msgID, dataList, TestingDutySlotGloas)
}

// PreConsensusBuilderRequestAuthWrongRootMsg carries a container whose count matches but one partial is
// for a divergent auth data — a valid message from a peer whose configured entries disagree, so its
// signing root fails the operator's expected-root check rather than mixing into a quorum.
var PreConsensusBuilderRequestAuthWrongRootMsg = func(msgSK *bls.SecretKey, msgID types.OperatorID, matchingData []byte) *types.PartialSignatureMessages {
	return builderRequestAuthMsg(msgSK, msgID, [][]byte{matchingData, []byte("divergent-builder-auth-data")}, TestingDutySlotGloas)
}

var builderRequestAuthMsg = func(
	sk *bls.SecretKey,
	id types.OperatorID,
	dataList [][]byte,
	proposalSlot phase0.Slot,
) *types.PartialSignatureMessages {
	signer := NewTestingKeyManager()
	beacon := NewTestingBeaconNode()
	epoch := types.BeaconTestNetwork.EstimatedEpochAtSlot(proposalSlot)
	d, _ := beacon.DomainData(epoch, types.DomainBuilderRequestAuth)

	msgs := &types.PartialSignatureMessages{
		Type:     types.RequestAuthPartialSig,
		Slot:     proposalSlot,
		Messages: []*types.PartialSignatureMessage{},
	}
	for _, data := range dataList {
		auth := TestingBuilderRequestAuth(data, proposalSlot)
		signed, root, _ := signer.SignBeaconObject(auth, d, sk.GetPublicKey().Serialize(), types.DomainBuilderRequestAuth)
		msgs.Messages = append(msgs.Messages, &types.PartialSignatureMessage{
			PartialSignature: signed[:],
			SigningRoot:      root,
			Signer:           id,
			ValidatorIndex:   TestingValidatorIndex,
		})
	}
	return msgs
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
