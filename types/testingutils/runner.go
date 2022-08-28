package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

var AttesterRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:]), keySet)
}

var AttesterRunner7Operators = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:]), keySet)
}

var ProposerRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleProposer, ssv.BeaconBlockValueCheck(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:]), keySet)
}

var AggregatorRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAggregator, ssv.AggregatorValueCheck(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:]), keySet)
}

var SyncCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleSyncCommittee, ssv.SyncCommitteeValueCheck(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:]), keySet)
}

var SyncCommitteeContributionRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleSyncCommitteeContribution, ssv.SyncCommitteeContributionValueCheck(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:]), keySet)
}

var baseRunner = func(role types.BeaconRole, valCheck qbft.ProposedValueCheckF, keySet *TestKeySet) ssv.Runner {
	share := TestingShare(keySet)
	identifier := types.NewMsgID(TestingValidatorPubKey[:], role)

	proposerF := func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}

	switch role {
	case types.BNRoleAttester:
		return ssv.NewAttesterRunnner(
			types.NowTestNetwork,
			share,
			NewTestingQBFTController(identifier[:], share, valCheck, proposerF),
			NewTestingBeaconNode(),
			NewTestingNetwork(),
			NewTestingKeyManager(),
			valCheck,
		)
	case types.BNRoleAggregator:
		return ssv.NewAggregatorRunner(
			types.NowTestNetwork,
			share,
			NewTestingQBFTController(identifier[:], share, valCheck, proposerF),
			NewTestingBeaconNode(),
			NewTestingNetwork(),
			NewTestingKeyManager(),
			valCheck,
		)
	case types.BNRoleProposer:
		return ssv.NewProposerRunner(
			types.NowTestNetwork,
			share,
			NewTestingQBFTController(identifier[:], share, valCheck, proposerF),
			NewTestingBeaconNode(),
			NewTestingNetwork(),
			NewTestingKeyManager(),
			valCheck,
		)
	case types.BNRoleSyncCommittee:
		return ssv.NewSyncCommitteeRunner(
			types.NowTestNetwork,
			share,
			NewTestingQBFTController(identifier[:], share, valCheck, proposerF),
			NewTestingBeaconNode(),
			NewTestingNetwork(),
			NewTestingKeyManager(),
			valCheck,
		)
	case types.BNRoleSyncCommitteeContribution:
		return ssv.NewSyncCommitteeAggregatorRunner(
			types.NowTestNetwork,
			share,
			NewTestingQBFTController(identifier[:], share, valCheck, proposerF),
			NewTestingBeaconNode(),
			NewTestingNetwork(),
			NewTestingKeyManager(),
			valCheck,
		)
	default:
		panic("unknown role type")
	}
}

var DecidedRunner = func(keySet *TestKeySet) ssv.Runner {
	return decideRunner(TestAttesterConsensusData, qbft.FirstHeight, keySet)
}

var DecidedRunnerWithHeight = func(height qbft.Height, keySet *TestKeySet) ssv.Runner {
	return decideRunner(TestAttesterConsensusData, height, keySet)
}

var DecidedRunnerUnknownDutyType = func(keySet *TestKeySet) ssv.Runner {
	return decideRunner(TestConsensusUnkownDutyTypeData, qbft.FirstHeight, keySet)
}

var decideRunner = func(consensusInput *types.ConsensusData, height qbft.Height, keySet *TestKeySet) ssv.Runner {
	v := BaseValidator(keySet)
	consensusDataByts, _ := consensusInput.Encode()
	msgs := DecidingMsgsForHeight(consensusDataByts, []byte{1, 2, 3, 4}, height, keySet)

	if err := v.DutyRunners[types.BNRoleAttester].StartNewDuty(consensusInput.Duty); err != nil {
		panic(err.Error())
	}
	for _, msg := range msgs {
		ssvMsg := SSVMsgAttester(msg, nil)
		if err := v.ProcessMessage(ssvMsg); err != nil {
			panic(err.Error())
		}
	}

	return v.DutyRunners[types.BNRoleAttester]
}

var SSVDecidingMsgs = func(consensusData []byte, ks *TestKeySet, role types.BeaconRole) []*types.SSVMessage {
	id := types.NewMsgID(TestingValidatorPubKey[:], role)

	ssvMsgF := func(qbftMsg *qbft.SignedMessage, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.SSVMessage {
		var byts []byte
		var msgType types.MsgType
		if partialSigMsg != nil {
			msgType = types.SSVPartialSignatureMsgType
			byts, _ = partialSigMsg.Encode()
		} else {
			msgType = types.SSVConsensusMsgType
			byts, _ = qbftMsg.Encode()
		}

		return &types.SSVMessage{
			MsgType: msgType,
			MsgID:   id,
			Data:    byts,
		}
	}

	// pre consensus msgs
	base := make([]*types.SSVMessage, 0)
	if role == types.BNRoleProposer {
		base = []*types.SSVMessage{
			ssvMsgF(nil, PreConsensusRandaoMsg(ks.Shares[1], 1)),
			ssvMsgF(nil, PreConsensusRandaoMsg(ks.Shares[2], 2)),
			ssvMsgF(nil, PreConsensusRandaoMsg(ks.Shares[3], 3)),
		}
	}
	if role == types.BNRoleAggregator {
		base = []*types.SSVMessage{
			ssvMsgF(nil, PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			ssvMsgF(nil, PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			ssvMsgF(nil, PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		base = []*types.SSVMessage{
			ssvMsgF(nil, PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
			ssvMsgF(nil, PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2)),
			ssvMsgF(nil, PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3)),
		}
	}

	qbftMsgs := DecidingMsgsForHeight(consensusData, id[:], qbft.FirstHeight, ks)
	for _, msg := range qbftMsgs {
		base = append(base, ssvMsgF(msg, nil))
	}
	return base
}

var DecidingMsgsForHeight = func(consensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	msgs := make([]*qbft.SignedMessage, 0)
	for h := qbft.Height(qbft.FirstHeight); h <= height; h++ {
		msgs = append(msgs, SignQBFTMsg(keySet.Shares[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     h,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Data:       ProposalDataBytes(consensusData, nil, nil),
		}))

		// prepare
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Data:       PrepareDataBytes(consensusData),
			}))
		}
		// commit
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Data:       CommitDataBytes(consensusData),
			}))
		}
	}
	return msgs
}
