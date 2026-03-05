package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

var SSVDecidingMsgsForCommitteeRunner = func(beaconVote *types.BeaconVote, ks *TestKeySet, height qbft.Height) []*types.SignedSSVMessage {
	id := CommitteeMsgID(ks)

	// consensus
	qbftMsgs := SSVDecidingMsgsForHeightAndBeaconVote(beaconVote, id[:], height, ks)
	return qbftMsgs
}

var SSVDecidingMsgsForAggregatorCommitteeRunnerForKS = func(duty *types.AggregatorCommitteeDuty, ks *TestKeySet, version spec.DataVersion) []*types.SignedSSVMessage {
	ksMap := make(map[phase0.ValidatorIndex]*TestKeySet)
	ksMap[duty.ValidatorDuties[0].ValidatorIndex] = ks
	slot := TestingDutySlotV(version)
	return SSVDecidingMsgsForAggregatorCommitteeRunner(duty, ksMap, qbft.Height(slot), version)
}

var SSVDecidingMsgsForAggregatorCommitteeRunner = func(duty *types.AggregatorCommitteeDuty, keySetMap map[phase0.ValidatorIndex]*TestKeySet, height qbft.Height, version spec.DataVersion) []*types.SignedSSVMessage {
	// ks is taken as any element from keySetMap
	ks := keySetMap[duty.ValidatorDuties[0].ValidatorIndex]

	id := AggregatorCommitteeMsgIDForKeySet(ks)

	signedF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, err := partialSigMsg.Encode()
		if err != nil {
			panic(err)
		}
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
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
	for i := uint64(1); i <= ks.Threshold; i++ {
		base = append(base, signedF(PreConsensusAggregatorCommitteeMsgForDuty(duty, keySetMap, i)))
	}

	// consensus
	accd := TestAggregatorCommitteeConsensusDataForDuty(duty, version, nil)
	fullData, err := accd.Encode()
	if err != nil {
		panic(err)
	}
	r, err := qbft.HashDataRoot(fullData)
	if err != nil {
		panic(err)
	}
	qbftMsgs := SSVDecidingMsgsForHeightWithRoot(r, fullData, id[:], height, ks)
	base = append(base, qbftMsgs...)

	return base
}

var QBFTMsgsForAggregatorCommitteeRunner = func(duty *types.AggregatorCommitteeDuty, ks *TestKeySet, height qbft.Height, version spec.DataVersion) []*types.SignedSSVMessage {
	id := AggregatorCommitteeMsgIDForKeySet(ks)
	// consensus
	accd := TestAggregatorCommitteeConsensusDataForDuty(duty, version, nil)
	fullData, err := accd.Encode()
	if err != nil {
		panic(err)
	}
	r, err := qbft.HashDataRoot(fullData)
	if err != nil {
		panic(err)
	}
	return SSVDecidingMsgsForHeightWithRoot(r, fullData, id[:], height, ks)
}

var SSVDecidingMsgsV = func(consensusData *types.ProposerConsensusData, ks *TestKeySet, role types.RunnerRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	signedF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, _ := partialSigMsg.Encode()
		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
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

	// consensus and post consensus
	qbftMsgs := SSVDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)
	return base
}

var ExpectedSSVDecidingMsgsV = func(consensusData *types.ProposerConsensusData, ks *TestKeySet, role types.RunnerRole) []*types.SignedSSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(partialSigMsg *types.PartialSignatureMessages) *types.SignedSSVMessage {
		byts, _ := partialSigMsg.Encode()

		ssvMsg := &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   id,
			Data:    byts,
		}
		signer := partialSigMsg.Messages[0].Signer
		sig, err := NewOperatorSigner(ks, signer).SignSSVMessage(ssvMsg)
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

	// consensus and post consensus
	qbftMsgs := SSVExpectedDecidingMsgsForHeight(consensusData, id[:], qbft.Height(consensusData.Duty.Slot), ks)
	base = append(base, qbftMsgs...)
	return base
}
