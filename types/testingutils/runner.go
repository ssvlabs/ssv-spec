package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingHighestDecidedSlot = phase0.Slot(0)

var CommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, ssv.BeaconVoteValueCheckF(NewTestingKeyManager(), TestingDutySlot, nil, TestingDutyEpoch), keySet)
}

var CommitteeRunnerWithKeySetMap = func(keySetMap map[phase0.ValidatorIndex]*TestKeySet) ssv.Runner {
	return baseRunnerWithKeySetMap(types.RoleCommittee, ssv.BeaconVoteValueCheckF(NewTestingKeyManager(), TestingDutySlot, nil, TestingDutyEpoch), keySetMap)
}

var AttesterRunner7Operators = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, ssv.BeaconVoteValueCheckF(NewTestingKeyManager(), TestingDutySlot, nil, TestingDutyEpoch), keySet)
}

var ProposerRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleProposer, ssv.ProposerValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, (types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex, nil), keySet)
}

var ProposerBlindedBlockRunner = func(keySet *TestKeySet) ssv.Runner {
	ret := baseRunner(
		types.RoleProposer,
		ssv.ProposerValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, (types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex, nil),
		keySet,
	)
	ret.(*ssv.ProposerRunner).ProducesBlindedBlocks = true
	return ret
}

var AggregatorRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleAggregator, ssv.AggregatorValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, (types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex), keySet)
}

var SyncCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, ssv.BeaconVoteValueCheckF(NewTestingKeyManager(), TestingDutySlot, nil, TestingDutyEpoch), keySet)
}

var SyncCommitteeContributionRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleSyncCommitteeContribution, ssv.SyncCommitteeContributionValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, (types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex), keySet)
}

var ValidatorRegistrationRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleValidatorRegistration, nil, keySet)
}

var VoluntaryExitRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleVoluntaryExit, nil, keySet)
}

var UnknownDutyTypeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(UnknownDutyType, UnknownDutyValueCheck(), keySet)
}

var baseRunnerWithKeySetMap = func(role types.RunnerRole, valCheck qbft.ProposedValueCheckF, keySetMap map[phase0.ValidatorIndex]*TestKeySet) ssv.Runner {
	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	for valIdx, keySet := range keySetMap {
		shareMap[valIdx] = TestingShare(keySet)
	}

	var keySetInstance *TestKeySet
	for _, keySet := range keySetMap {
		keySetInstance = keySet
		break
	}

	// Identifier
	ownerID := []byte{}
	if role == types.RoleCommittee {
		committee := make([]uint64, 0)
		for _, op := range keySetInstance.Committee() {
			committee = append(committee, op.Signer)
		}
		clusterID := types.GetClusterID(committee)
		copy(ownerID, clusterID[:])
	} else {
		ownerID = TestingValidatorPubKey[:]
	}
	identifier := types.NewMsgID(TestingSSVDomainType, ownerID[:], role)

	net := NewTestingNetwork(1, keySetInstance.OperatorKeys[1])

	km := NewTestingKeyManager()
	operator := TestingOperator(keySetInstance)
	opSigner := NewTestingOperatorSigner(keySetInstance, operator.OperatorID)

	config := TestingConfig(keySetInstance)
	config.ValueCheckF = valCheck
	config.ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}
	config.Network = net
	config.OperatorSigner = opSigner
	config.SignatureVerifier = NewTestingVerifier()

	contr := qbft.NewController(
		identifier[:],
		operator,
		config,
	)

	switch role {
	case types.RoleCommittee:
		return ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
		)
	case types.RoleAggregator:
		return ssv.NewAggregatorRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.RoleProposer:
		return ssv.NewProposerRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.RoleSyncCommitteeContribution:
		return ssv.NewSyncCommitteeAggregatorRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.RoleValidatorRegistration:
		return ssv.NewValidatorRegistrationRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case types.RoleVoluntaryExit:
		return ssv.NewVoluntaryExitRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case UnknownDutyType:
		ret := ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
		)
		ret.(*ssv.CommitteeRunner).BaseRunner.RunnerRoleType = UnknownDutyType
		return ret
	default:
		panic("unknown role type")
	}
}

var baseRunner = func(role types.RunnerRole, valCheck qbft.ProposedValueCheckF, keySet *TestKeySet) ssv.Runner {
	share := TestingShare(keySet)

	// Identifier
	ownerID := []byte{}
	if role == types.RoleCommittee {
		committee := make([]uint64, 0)
		for _, op := range keySet.Committee() {
			committee = append(committee, op.Signer)
		}
		clusterID := types.GetClusterID(committee)
		copy(ownerID, clusterID[:])
	} else {
		ownerID = TestingValidatorPubKey[:]
	}
	identifier := types.NewMsgID(TestingSSVDomainType, ownerID[:], role)

	net := NewTestingNetwork(1, keySet.OperatorKeys[1])
	km := NewTestingKeyManager()
	operator := TestingOperator(keySet)
	opSigner := NewTestingOperatorSigner(keySet, operator.OperatorID)

	config := TestingConfig(keySet)
	config.ValueCheckF = valCheck
	config.ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}
	config.Network = net
	config.OperatorSigner = opSigner
	config.SignatureVerifier = NewTestingVerifier()

	contr := qbft.NewController(
		identifier[:],
		operator,
		config,
	)

	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	shareMap[share.ValidatorIndex] = share

	switch role {
	case types.RoleCommittee:
		return ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
		)
	case types.RoleAggregator:
		return ssv.NewAggregatorRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.RoleProposer:
		return ssv.NewProposerRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.RoleSyncCommitteeContribution:
		return ssv.NewSyncCommitteeAggregatorRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.RoleValidatorRegistration:
		return ssv.NewValidatorRegistrationRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case types.RoleVoluntaryExit:
		return ssv.NewVoluntaryExitRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case UnknownDutyType:
		ret := ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
		)
		ret.(*ssv.CommitteeRunner).BaseRunner.RunnerRoleType = UnknownDutyType
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
	msgs := SSVDecidingMsgsForHeight(consensusInput, AttesterMsgID, height, keySet)

	if err := v.DutyRunners[types.RoleCommittee].StartNewDuty(&consensusInput.Duty); err != nil {
		panic(err.Error())
	}
	for _, msg := range msgs {
		if err := v.ProcessMessage(msg); err != nil {
			panic(err.Error())
		}
	}

	return v.DutyRunners[types.RoleCommittee]
}

// //////////////////////////////// For SSV Tests ////////////////////////////////////////////////////////////////
var SSVDecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return SSVDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var SSVDecidingMsgsForHeightAndBeaconVote = func(beaconVote *types.BeaconVote, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	fullData, err := beaconVote.Encode()
	if err != nil {
		panic(err)
	}
	r, err := qbft.HashDataRoot(fullData)
	if err != nil {
		panic(err)
	}

	return SSVDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var SSVExpectedDecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return SSVExpectedDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var SSVDecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	msgs := make([]*types.SignedSSVMessage, 0)

	// proposal
	s := SignQBFTMsg(keySet.OperatorKeys[1], 1, &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     height,
		Round:      qbft.FirstRound,
		Identifier: msgIdentifier,
		Root:       root,
	})
	s.FullData = fullData
	msgs = append(msgs, s)

	// prepare
	for i := uint64(1); i <= keySet.Threshold; i++ {
		msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	// commit
	for i := uint64(1); i <= keySet.Threshold; i++ {
		msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	return msgs
}

var SSVExpectedDecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	msgs := make([]*types.SignedSSVMessage, 0)

	// proposal
	s := SignQBFTMsg(keySet.OperatorKeys[1], 1, &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     height,
		Round:      qbft.FirstRound,
		Identifier: msgIdentifier,
		Root:       root,
	})
	s.FullData = fullData
	msgs = append(msgs, s)

	// prepare
	for i := uint64(1); i <= keySet.Threshold; i++ {
		msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	// commit
	for i := uint64(1); i <= keySet.Threshold; i++ {
		msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	return msgs
}

// //////////////////// For QBFT TESTS /////////////////////////////////////////////////////////////////////////
var DecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return DecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var ExpectedDecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return ExpectedDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var DecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	msgs := make([]*types.SignedSSVMessage, 0)

	for h := qbft.FirstHeight; h <= height; h++ {
		// proposal
		s := SignQBFTMsg(keySet.OperatorKeys[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     h,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		})
		s.FullData = fullData
		msgs = append(msgs, s)

		// prepare
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Root:       root,
			}))
		}
		// commit
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Root:       root,
			}))
		}
	}
	return msgs
}

var ExpectedDecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	msgs := make([]*types.SignedSSVMessage, 0)

	for h := qbft.FirstHeight; h <= height; h++ {
		// proposal
		s := SignQBFTMsg(keySet.OperatorKeys[1], 1, &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     h,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		})
		s.FullData = fullData
		msgs = append(msgs, s)

		// prepare
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Root:       root,
			}))
		}
		// commit
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.OperatorKeys[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Root:       root,
			}))
		}
	}
	return msgs
}
