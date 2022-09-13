package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

var AttesterRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.AttesterValueCheckF(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var AttesterRunner7Operators = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.AttesterValueCheckF(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var ProposerRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleProposer, ssv.ProposerValueCheckF(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var AggregatorRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAggregator, ssv.AggregatorValueCheckF(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var SyncCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleSyncCommittee, ssv.SyncCommitteeValueCheckF(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var SyncCommitteeContributionRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleSyncCommitteeContribution, ssv.SyncCommitteeContributionValueCheckF(NewTestingKeyManager(), types.NowTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var UnknownDutyTypeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(UnknownDutyType, UnknownDutyValueCheck(), keySet)
}

var baseRunner = func(role types.BeaconRole, valCheck qbft.ProposedValueCheckF, keySet *TestKeySet) ssv.Runner {
	share := TestingShare(keySet)
	identifier := types.NewBaseMsgID(TestingValidatorPubKey[:], role)

	proposerF := func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}

	net := NewTestingNetwork()
	km := NewTestingKeyManager()
	contr := qbft.NewController(
		identifier,
		share,
		types.PrimusTestnet,
		km,
		valCheck,
		NewTestingStorage(),
		net,
		proposerF,
	)

	switch role {
	case types.BNRoleAttester:
		return ssv.NewAttesterRunnner(
			types.NowTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
		)
	case types.BNRoleAggregator:
		return ssv.NewAggregatorRunner(
			types.NowTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
		)
	case types.BNRoleProposer:
		return ssv.NewProposerRunner(
			types.NowTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
		)
	case types.BNRoleSyncCommittee:
		return ssv.NewSyncCommitteeRunner(
			types.NowTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
		)
	case types.BNRoleSyncCommitteeContribution:
		return ssv.NewSyncCommitteeAggregatorRunner(
			types.NowTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
		)
	case UnknownDutyType:
		ret := ssv.NewAttesterRunnner(
			types.NowTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
		)
		ret.(*ssv.AttesterRunner).BaseRunner.BeaconRoleType = UnknownDutyType
		return ret
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
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msgs := DecidingMsgsForHeight(consensusDataByts, identifier, height, keySet)

	if err := v.DutyRunners[types.BNRoleAttester].StartNewDuty(consensusInput.Duty); err != nil {
		panic(err.Error())
	}
	for _, msg := range msgs {
		// TODO<olegshmuelov>: check if need to use SSVMsgAttester
		//ssvMsg := SSVMsgAttester(msg, nil)
		if err := v.ProcessMessage(msg); err != nil {
			panic(err.Error())
		}
	}

	return v.DutyRunners[types.BNRoleAttester]
}

var SSVDecidingMsgs = func(consensusData []byte, ks *TestKeySet, role types.BeaconRole) []*types.Message {
	id := types.NewBaseMsgID(TestingValidatorPubKey[:], role)

	ssvMsgF := func(qbftMsg *types.Message, partialSigMsg *ssv.SignedPartialSignatureMessage) *types.Message {
		//TODO<olegshmuelov>: rethink logic
		/*var byts []byte
		var msgType types.MsgType
		if partialSigMsg != nil {
			msgType = types.SSVPartialSignatureMsgType
			byts, _ = partialSigMsg.Encode()
		} else {
			msgType = types.SSVConsensusMsgType
			byts, _ = qbftMsg.Encode()
		}

		return &types.Message{
			MsgType: msgType,
			MsgID:   id,
			Data:    byts,
		}*/
		return &types.Message{}
	}

	// pre consensus msgs
	base := make([]*types.Message, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(nil, PreConsensusRandaoMsg(ks.Shares[types.OperatorID(i)], types.OperatorID(i))))
		}
	}
	if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(nil, PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(nil, PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	qbftMsgs := DecidingMsgsForHeight(consensusData, id, qbft.FirstHeight, ks)
	for _, msg := range qbftMsgs {
		base = append(base, ssvMsgF(msg, nil))
	}
	return base
}

var DecidingMsgsForHeight = func(consensusData []byte, msgID types.MessageID, height qbft.Height, keySet *TestKeySet) []*types.Message {
	//msgs := make([]*qbft.SignedMessage, 0)
	msgs := make([]*types.Message, 0)
	for h := qbft.FirstHeight; h <= height; h++ {
		signMsgEncoded, _ := SignQBFTMsg(keySet.Shares[1], types.OperatorID(1), &qbft.Message{
			Height: h,
			Round:  qbft.FirstRound,
			Input:  []byte{1, 2, 3, 4},
		}).Encode()

		msgID = types.PopulateMsgType(msgID, types.ConsensusProposeMsgType)
		msgs = append(msgs, &types.Message{
			ID:   msgID,
			Data: signMsgEncoded,
		})

		//msgs = append(msgs, SignQBFTMsg(keySet.Shares[1], 1, &qbft.Message{
		//	MsgType:    qbft.ProposalMsgType,
		//	Height:     h,
		//	Round:      qbft.FirstRound,
		//	Identifier: msgIdentifier,
		//	Data:       ProposalDataBytes(consensusData, nil, nil),
		//}))

		msgID = types.PopulateMsgType(msgID, types.ConsensusPrepareMsgType)
		// prepare
		for i := uint64(1); i <= keySet.Threshold; i++ {
			signMsgEncoded, _ := SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				Height: h,
				Round:  qbft.FirstRound,
				Input:  []byte{1, 2, 3, 4},
			}).Encode()
			msgs = append(msgs, &types.Message{
				ID:   msgID,
				Data: signMsgEncoded,
			})
			//msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			//	MsgType:    qbft.PrepareMsgType,
			//	Height:     h,
			//	Round:      qbft.FirstRound,
			//	Identifier: msgID,
			//	Data:       PrepareDataBytes(consensusData),
			//}))
		}

		msgID = types.PopulateMsgType(msgID, types.ConsensusCommitMsgType)
		// commit
		for i := uint64(1); i <= keySet.Threshold; i++ {
			signMsgEncoded, _ := SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				Height: h,
				Round:  qbft.FirstRound,
				Input:  []byte{1, 2, 3, 4},
			}).Encode()
			msgs = append(msgs, &types.Message{
				ID:   msgID,
				Data: signMsgEncoded,
			})
			//msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			//	MsgType:    qbft.CommitMsgType,
			//	Height:     h,
			//	Round:      qbft.FirstRound,
			//	Identifier: msgID,
			//	Data:       CommitDataBytes(consensusData),
			//}))
		}
	}
	return msgs
}
