package testingutils

import (
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

var SSVDecidingMsgsForCommitteeRunner = func(beaconVote *types.BeaconVote, ks *TestKeySet, height qbft.Height) []*types.SignedSSVMessage {
	id := CommitteeMsgID(ks)

	// consensus
	qbftMsgs := SSVDecidingMsgsForHeightAndBeaconVote(beaconVote, id[:], height, ks)
	return qbftMsgs
}

var SSVDecidingMsgsV = func(consensusData *types.ValidatorConsensusData, ks *TestKeySet, role types.RunnerRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	signedF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, _ := partialSigMsg.Encode()
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewTestingOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
		if err != nil {
			panic(err)
		}
		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{signer},
			Signatures:  [][]byte{sig},
			SSVMessage:  ssvMsg,
		}
	}

	// pre consensus msgs
	base := make([]*types.SignedSSVMessage, 0)
	if role == types.RoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedF(PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
		}
	}
	if role == types.RoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedF(PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.RoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, signedF(PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)
	return base
}

var ExpectedSSVDecidingMsgsV = func(consensusData *types.ValidatorConsensusData, ks *TestKeySet, role types.RunnerRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, _ := partialSigMsg.Encode()

		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewTestingOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
		if err != nil {
			panic(err)
		}
		return &types.SignedSSVMessage{
			OperatorIDs: []types.OperatorID{signer},
			Signatures:  [][]byte{sig},
			SSVMessage:  ssvMsg,
		}
	}

	// pre consensus msgs
	base := make([]*types.SignedSSVMessage, 0)
	if role == types.RoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
		}
	}
	if role == types.RoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.RoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// consensus and post consensus
	qbftMsgs := SSVExpectedDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)
	return base
}
