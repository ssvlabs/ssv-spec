package testingutils

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

var AttesterRunner = func(keySet *TestKeySet) *ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), types.NowTestNetwork), keySet)
}

var AttesterRunner7Operators = func(keySet *TestKeySet) *ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.BeaconAttestationValueCheck(NewTestingKeyManager(), types.NowTestNetwork), keySet)
}

var ProposerRunner = func(keySet *TestKeySet) *ssv.Runner {
	return baseRunner(types.BNRoleProposer, ssv.BeaconBlockValueCheck(NewTestingKeyManager(), types.NowTestNetwork), keySet)
}

var AggregatorRunner = func(keySet *TestKeySet) *ssv.Runner {
	return baseRunner(types.BNRoleAggregator, ssv.AggregatorValueCheck(NewTestingKeyManager(), types.NowTestNetwork), keySet)
}

var SyncCommitteeRunner = func(keySet *TestKeySet) *ssv.Runner {
	return baseRunner(types.BNRoleSyncCommittee, ssv.SyncCommitteeValueCheck(NewTestingKeyManager(), types.NowTestNetwork), keySet)
}

var SyncCommitteeContributionRunner = func(keySet *TestKeySet) *ssv.Runner {
	return baseRunner(types.BNRoleSyncCommitteeContribution, ssv.SyncCommitteeContributionValueCheck(NewTestingKeyManager(), types.NowTestNetwork), keySet)
}

var baseRunner = func(role types.BeaconRole, valCheck qbft.ProposedValueCheckF, keySet *TestKeySet) *ssv.Runner {
	share := TestingShare(keySet)
	identifier := types.NewBaseMsgID(TestingValidatorPubKey[:], role)

	proposerF := func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}

	return ssv.NewDutyRunner(
		role,
		types.NowTestNetwork,
		share,
		NewTestingQBFTController(identifier, share, valCheck, proposerF),
		NewTestingStorage(),
		valCheck,
	)
}

var DecidedRunner = func(keySet *TestKeySet) *ssv.Runner {
	return decideRunner(TestAttesterConsensusDataByts, qbft.FirstHeight, keySet)
}

var DecidedRunnerWithHeight = func(height qbft.Height, keySet *TestKeySet) *ssv.Runner {
	return decideRunner(TestAttesterConsensusDataByts, height, keySet)
}

var DecidedRunnerUnknownDutyType = func(keySet *TestKeySet) *ssv.Runner {
	return decideRunner(TestConsensusUnkownDutyTypeDataByts, qbft.FirstHeight, keySet)
}

var decideRunner = func(consensusData []byte, height qbft.Height, keySet *TestKeySet) *ssv.Runner {
	v := BaseValidator(keySet)
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	msgs := DecidingMsgsForHeight(consensusData, identifier, height, keySet)

	if err := v.DutyRunners[types.BNRoleAttester].Decide(TestAttesterConsensusData); err != nil {
		panic(err.Error())
	}
	for _, msg := range msgs {
		//ssvMsg := SSVMsgAttester(msg, nil)
		if err := v.ProcessMessageSIP(msg); err != nil {
			panic(err.Error())
		}
	}

	return v.DutyRunners[types.BNRoleAttester]
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
