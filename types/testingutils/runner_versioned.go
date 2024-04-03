package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

var SSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.BeaconRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	// pre consensus msgs
	base := make([]*types.SignedSSVMessage, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedSSVMsgWithMessageID(i, ks.NetworkKeys[i], nil, PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version), id))
		}
	}
	if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedSSVMsgWithMessageID(i, ks.NetworkKeys[i], nil, PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], types.OperatorID(i)), id))
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedSSVMsgWithMessageID(i, ks.NetworkKeys[i], nil, PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], types.OperatorID(i)), id))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)

	return base
}

var ExpectedSSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.BeaconRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	// pre consensus msgs
	base := make([]*types.SignedSSVMessage, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {

			base = append(base, signedSSVMsgWithMessageID(i, ks.NetworkKeys[i], nil, PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version), id))
		}
	}
	if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedSSVMsgWithMessageID(i, ks.NetworkKeys[i], nil, PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], types.OperatorID(i)), id))
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedSSVMsgWithMessageID(i, ks.NetworkKeys[i], nil, PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], types.OperatorID(i)), id))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVExpectedDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)

	return base
}
