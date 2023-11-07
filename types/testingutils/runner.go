package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
)

var TestingHighestDecidedSlot = phase0.Slot(0)

var AttesterRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.AttesterValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex, nil), keySet)
}

var AttesterRunner7Operators = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAttester, ssv.AttesterValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex, nil), keySet)
}

var ProposerRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleProposer, ssv.ProposerValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex, nil), keySet)
}

var ProposerBlindedBlockRunner = func(keySet *TestKeySet) ssv.Runner {
	ret := baseRunner(
		types.BNRoleProposer,
		ssv.ProposerValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex, nil),
		keySet,
	)
	ret.(*ssv.ProposerRunner).ProducesBlindedBlocks = true
	return ret
}

var AggregatorRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleAggregator, ssv.AggregatorValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var SyncCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleSyncCommittee, ssv.SyncCommitteeValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var SyncCommitteeContributionRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleSyncCommitteeContribution, ssv.SyncCommitteeContributionValueCheckF(NewTestingKeyManager(), types.BeaconTestNetwork, TestingValidatorPubKey[:], TestingValidatorIndex), keySet)
}

var ValidatorRegistrationRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleValidatorRegistration, nil, keySet)
}

var VoluntaryExitRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.BNRoleVoluntaryExit, nil, keySet)
}

var UnknownDutyTypeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(UnknownDutyType, UnknownDutyValueCheck(), keySet)
}

var baseRunner = func(role types.BeaconRole, valCheck qbft.ProposedValueCheckF, keySet *TestKeySet) ssv.Runner {
	share := TestingShare(keySet)
	identifier := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)
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
		identifier[:],
		share,
		config,
	)

	switch role {
	case types.BNRoleAttester:
		return ssv.NewAttesterRunnner(
			types.BeaconTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.BNRoleAggregator:
		return ssv.NewAggregatorRunner(
			types.BeaconTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.BNRoleProposer:
		return ssv.NewProposerRunner(
			types.BeaconTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.BNRoleSyncCommittee:
		return ssv.NewSyncCommitteeRunner(
			types.BeaconTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.BNRoleSyncCommitteeContribution:
		return ssv.NewSyncCommitteeAggregatorRunner(
			types.BeaconTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
			TestingHighestDecidedSlot,
		)
	case types.BNRoleValidatorRegistration:
		return ssv.NewValidatorRegistrationRunner(
			types.PraterNetwork,
			share,
			NewTestingBeaconNode(),
			net,
			km,
		)
	case types.BNRoleVoluntaryExit:
		return ssv.NewVoluntaryExitRunner(
			types.PraterNetwork,
			share,
			NewTestingBeaconNode(),
			net,
			km,
		)
	case UnknownDutyType:
		ret := ssv.NewAttesterRunnner(
			types.BeaconTestNetwork,
			share,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			valCheck,
			TestingHighestDecidedSlot,
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
	msgs := SSVDecidingMsgsForHeight(consensusInput, []byte{1, 2, 3, 4}, height, keySet)

	if err := v.DutyRunners[types.BNRoleAttester].StartNewDuty(&consensusInput.Duty); err != nil {
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

// //////////////////////////////// For SSV Tests ////////////////////////////////////////////////////////////////
var SSVDecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return SSVDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var SSVExpectedDecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return SSVExpectedDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var SSVDecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	msgs := make([]*qbft.SignedMessage, 0)

	// proposal
	s := SignQBFTMsg(keySet.Shares[1], 1, &qbft.Message{
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
		msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	// commit
	for i := uint64(1); i <= keySet.Threshold; i++ {
		msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	return msgs
}

var SSVExpectedDecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	msgs := make([]*qbft.SignedMessage, 0)

	// proposal
	s := SignQBFTMsg(keySet.Shares[1], 1, &qbft.Message{
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
		msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     height,
			Round:      qbft.FirstRound,
			Identifier: msgIdentifier,
			Root:       root,
		}))
	}
	// commit
	for i := uint64(1); i <= keySet.Threshold; i++ {
		msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
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
var DecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return DecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var ExpectedDecidingMsgsForHeight = func(consensusData *types.ConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return ExpectedDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var DecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	msgs := make([]*qbft.SignedMessage, 0)

	for h := qbft.FirstHeight; h <= height; h++ {
		// proposal
		s := SignQBFTMsg(keySet.Shares[1], 1, &qbft.Message{
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
			msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Root:       root,
			}))
		}
		// commit
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
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

var ExpectedDecidingMsgsForHeightWithRoot = func(root [32]byte, fullData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*qbft.SignedMessage {
	msgs := make([]*qbft.SignedMessage, 0)

	for h := qbft.FirstHeight; h <= height; h++ {
		// proposal
		s := SignQBFTMsg(keySet.Shares[1], 1, &qbft.Message{
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
			msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     h,
				Round:      qbft.FirstRound,
				Identifier: msgIdentifier,
				Root:       root,
			}))
		}
		// commit
		for i := uint64(1); i <= keySet.Threshold; i++ {
			msgs = append(msgs, SignQBFTMsg(keySet.Shares[types.OperatorID(i)], types.OperatorID(i), &qbft.Message{
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
