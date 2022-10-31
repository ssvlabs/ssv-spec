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
	net := NewTestingNetwork()
	km := NewTestingKeyManager()

	config := TestingConfig(keySet)
	config.ValueCheckF = valCheck
	config.ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}
	config.Network = net
	config.Signer = km

	contr := qbft.NewController(
		identifier,
		share,
		types.PrimusTestnet,
		config,
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
	consensusDataByts, _ := consensusInput.MarshalSSZ()
	identifier := types.NewBaseMsgID([]byte{1, 2, 3, 4}, types.BNRoleAttester)
	// TODO<olegshmuelov>: pass InputData
	msgs := DecidingMsgsForHeight(&qbft.Data{
		Root:   [32]byte{},
		Source: consensusDataByts,
	}, identifier, height, keySet)

	if err := v.DutyRunners[types.BNRoleAttester].StartNewDuty(consensusInput.Duty); err != nil {
		panic(err.Error())
	}
	for _, msg := range msgs {
		if err := v.ProcessMessage(msg); err != nil {
			panic(err.Error())
		}
	}

	return v.DutyRunners[types.BNRoleAttester]
}

var SSVDecidingMsgs = func(consensusData []byte, ks *TestKeySet, role types.BeaconRole) []*types.Message {
	id := types.NewBaseMsgID(TestingValidatorPubKey[:], role)

	ssvMsgF := func(msgType types.MsgType, partialSigMsg *ssv.SignedPartialSignature) *types.Message {
		partialSigMsgEncoded, _ := partialSigMsg.Encode()
		return &types.Message{
			ID:   types.PopulateMsgType(id, msgType),
			Data: partialSigMsgEncoded,
		}
	}

	// pre consensus msgs
	base := make([]*types.Message, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(types.PartialRandaoSignatureMsgType, PreConsensusRandaoMsg(ks.Shares[types.OperatorID(i)], types.OperatorID(i))))
		}
	}
	if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(types.PartialSelectionProofSignatureMsgType, PreConsensusSelectionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}
	if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(types.PartialContributionProofSignatureMsgType, PreConsensusContributionProofMsg(ks.Shares[types.OperatorID(i)], ks.Shares[types.OperatorID(i)], types.OperatorID(i), types.OperatorID(i))))
		}
	}

	// TODO<olegshmuelov>: pass InputData
	qbftMsgs := DecidingMsgsForHeight(&qbft.Data{
		Root:   [32]byte{},
		Source: consensusData,
	}, id, qbft.FirstHeight, ks)
	base = append(base, qbftMsgs...)
	return base
}

var DecidingMsgsForHeight = func(consensusData *qbft.Data, msgID types.MessageID, height qbft.Height, keySet *TestKeySet) []*types.Message {
	msgs := make([]*types.Message, 0)
	for h := qbft.FirstHeight; h <= height; h++ {
		signMsgEncoded, _ := SignQBFTMsg(keySet.Shares[1], types.OperatorID(1), &qbft.Message{
			Height: h,
			Round:  qbft.FirstRound,
			Input:  consensusData,
		}).Encode()

		msgs = append(msgs, &types.Message{
			ID:   types.PopulateMsgType(msgID, types.ConsensusProposeMsgType),
			Data: signMsgEncoded,
		})

		// prepare
		msgID = types.PopulateMsgType(msgID, types.ConsensusPrepareMsgType)
		for i := uint64(1); i <= keySet.Threshold; i++ {
			signMsgEncoded, _ := SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				Height: h,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: consensusData.Root},
			}).Encode()
			msgs = append(msgs, &types.Message{
				ID:   msgID,
				Data: signMsgEncoded,
			})
		}

		// commit
		msgID = types.PopulateMsgType(msgID, types.ConsensusCommitMsgType)
		for i := uint64(1); i <= keySet.Threshold; i++ {
			signMsgEncoded, _ := SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				Height: h,
				Round:  qbft.FirstRound,
				Input:  &qbft.Data{Root: consensusData.Root},
			}).Encode()
			msgs = append(msgs, &types.Message{
				ID:   msgID,
				Data: signMsgEncoded,
			})
		}
	}
	return msgs
}
