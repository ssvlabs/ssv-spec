package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

var SSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.BeaconRole) []*types.SSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
		byts, _ := partialSigMsg.Encode()
		return &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
	}

	// pre consensus msgs
	base := make([]*types.SSVMessage, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
		}
	}
	if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	for _, msg := range qbftMsgs {
		base = append(base, msg.SSVMessage)
	}
	return base
}

var ExpectedSSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.BeaconRole) []*types.SSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
		byts, _ := partialSigMsg.Encode()

		return &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
	}

	// pre consensus msgs
	base := make([]*types.SSVMessage, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
		}
	}
	if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVExpectedDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	for _, msg := range qbftMsgs {
		base = append(base, msg.SSVMessage)
	}
	return base
}
