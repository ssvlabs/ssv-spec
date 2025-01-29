package testingutils

import (
	"bytes"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
)

var TestingHighestDecidedSlot = phase0.Slot(0)

var CommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, keySet)
}

var CommitteeRunnerWithShareMap = func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
	return baseRunnerWithShareMap(types.RoleCommittee, shareMap)
}

var AttesterRunner7Operators = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, keySet)
}

var ProposerRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleProposer, keySet)
}

var ProposerBlindedBlockRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleProposer, keySet)
}

var AggregatorRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleAggregator, keySet)
}

var SyncCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, keySet)
}

var SyncCommitteeContributionRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleSyncCommitteeContribution, keySet)
}

var ValidatorRegistrationRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleValidatorRegistration, keySet)
}

var VoluntaryExitRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleVoluntaryExit, keySet)
}

var UnknownDutyTypeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(UnknownDutyType, keySet)
}

var baseRunnerWithShareMap = func(role types.RunnerRole, shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
	runner, err := ConstructBaseRunnerWithShareMap(role, shareMap)
	if err != nil {
		panic(err)
	}
	return runner
}

var ConstructBaseRunnerWithShareMap = func(role types.RunnerRole, shareMap map[phase0.ValidatorIndex]*types.Share) (ssv.Runner, error) {

	var identifier types.MessageID
	var net *TestingNetwork
	var opSigner *types.OperatorSigner
	var valCheck qbft.ProposedValueCheckF
	var contr *qbft.Controller

	km := NewTestingKeyManager()

	if len(shareMap) > 0 {
		// Get sample instance for share and key set
		var keySetInstance *TestKeySet
		var shareInstance *types.Share
		for _, share := range shareMap {
			keySetInstance = KeySetForShare(share)
			shareInstance = TestingShare(keySetInstance, share.ValidatorIndex)
			break
		}

		// Get list of shares' pub keys
		sharePubKeys := make([]types.ShareValidatorPK, 0)
		for _, share := range shareMap {
			sharePubKeys = append(sharePubKeys, share.SharePubKey)
		}

		// Identifier
		var ownerID []byte
		if role == types.RoleCommittee {
			committee := make([]uint64, 0)
			for _, op := range keySetInstance.Committee() {
				committee = append(committee, op.Signer)
			}
			committeeID := types.GetCommitteeID(committee)
			ownerID = bytes.Clone(committeeID[:])
		} else {
			ownerID = TestingValidatorPubKey[:]
		}
		identifier = types.NewMsgID(TestingSSVDomainType, ownerID, role)

		// Network
		net = NewTestingNetwork(1, keySetInstance.OperatorKeys[1])

		// Create CommitteeMember
		committeeMember := TestingCommitteeMember(keySetInstance)

		// Create OperatorSigner
		opSigner = NewOperatorSigner(keySetInstance, committeeMember.OperatorID)

		// Create ValueCheck
		switch role {
		case types.RoleCommittee:
			valCheck = ssv.BeaconVoteValueCheckF(km, TestingDutySlot,
				sharePubKeys, TestingDutyEpoch)
		case types.RoleProposer:
			valCheck = ssv.ProposerValueCheckF(km, types.BeaconTestNetwork,
				(types.ValidatorPK)(shareInstance.ValidatorPubKey), shareInstance.ValidatorIndex, shareInstance.SharePubKey)
		case types.RoleAggregator:
			valCheck = ssv.AggregatorValueCheckF(km, types.BeaconTestNetwork,
				(types.ValidatorPK)(shareInstance.ValidatorPubKey), shareInstance.ValidatorIndex)
		case types.RoleSyncCommitteeContribution:
			valCheck = ssv.SyncCommitteeContributionValueCheckF(km, types.BeaconTestNetwork,
				(types.ValidatorPK)(shareInstance.ValidatorPubKey), shareInstance.ValidatorIndex)
		default:
			valCheck = nil
		}

		// Create qbft.Controller
		config := TestingConfig(keySetInstance)
		config.ValueCheckF = valCheck
		config.ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
			return 1
		}
		config.Network = net

		contr = qbft.NewController(identifier[:], committeeMember, config, opSigner)
	}

	var runner ssv.Runner
	var err error
	switch role {
	case types.RoleCommittee:
		runner, err = ssv.NewCommitteeRunner(
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
		runner, err = ssv.NewAggregatorRunner(
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
		runner, err = ssv.NewProposerRunner(
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
		runner, err = ssv.NewSyncCommitteeAggregatorRunner(
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
		runner, err = ssv.NewValidatorRegistrationRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			types.DefaultGasLimit,
		)
	case types.RoleVoluntaryExit:
		runner, err = ssv.NewVoluntaryExitRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case UnknownDutyType:
		runner, err = ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
		)
		if runner != nil {
			runner.(*ssv.CommitteeRunner).BaseRunner.RunnerRoleType = UnknownDutyType
		}
	default:
		return nil, errors.New("unknown role type")
	}
	return runner, err
}

var baseRunner = func(role types.RunnerRole, keySet *TestKeySet) ssv.Runner {
	runner, err := ConstructBaseRunner(role, keySet)
	if err != nil {
		panic(err)
	}
	return runner
}

var ConstructBaseRunner = func(role types.RunnerRole, keySet *TestKeySet) (ssv.Runner, error) {
	share := TestingShare(keySet, TestingValidatorIndex)
	km := NewTestingKeyManager()

	// Identifier
	var ownerID []byte
	if role == types.RoleCommittee {
		committee := make([]uint64, 0)
		for _, op := range keySet.Committee() {
			committee = append(committee, op.Signer)
		}
		clusterID := types.GetCommitteeID(committee)
		ownerID = clusterID[:]
	} else {
		ownerID = TestingValidatorPubKey[:]
	}
	identifier := types.NewMsgID(TestingSSVDomainType, ownerID[:], role)

	// Network
	net := NewTestingNetwork(1, keySet.OperatorKeys[1])

	// Create CommitteeMember
	committeeMember := TestingCommitteeMember(keySet)

	// Create OperatorSigner
	opSigner := NewOperatorSigner(keySet, committeeMember.OperatorID)

	// Create ValueCheck
	var valCheck qbft.ProposedValueCheckF
	switch role {
	case types.RoleCommittee:
		valCheck = ssv.BeaconVoteValueCheckF(km, TestingDutySlot,
			[]types.ShareValidatorPK{share.SharePubKey}, TestingDutyEpoch)
	case types.RoleProposer:
		valCheck = ssv.ProposerValueCheckF(km, types.BeaconTestNetwork,
			(types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex, share.SharePubKey)
	case types.RoleAggregator:
		valCheck = ssv.AggregatorValueCheckF(km, types.BeaconTestNetwork,
			(types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex)
	case types.RoleSyncCommitteeContribution:
		valCheck = ssv.SyncCommitteeContributionValueCheckF(km, types.BeaconTestNetwork,
			(types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex)
	default:
		valCheck = nil
	}

	// Create qbft.Controller
	config := TestingConfig(keySet)
	config.ValueCheckF = valCheck
	config.ProposerF = func(state *qbft.State, round qbft.Round) types.OperatorID {
		return 1
	}
	config.Network = net

	contr := qbft.NewController(identifier[:], committeeMember, config, opSigner)

	// Build share map
	shareMap := make(map[phase0.ValidatorIndex]*types.Share)
	shareMap[share.ValidatorIndex] = share

	var runner ssv.Runner
	var err error
	switch role {
	case types.RoleCommittee:
		runner, err = ssv.NewCommitteeRunner(
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
		runner, err = ssv.NewAggregatorRunner(
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
		runner, err = ssv.NewProposerRunner(
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
		runner, err = ssv.NewSyncCommitteeAggregatorRunner(
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
		runner, err = ssv.NewValidatorRegistrationRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			types.DefaultGasLimit,
		)
	case types.RoleVoluntaryExit:
		runner, err = ssv.NewVoluntaryExitRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case UnknownDutyType:
		runner, err = ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
		)
		if runner != nil {
			runner.(*ssv.CommitteeRunner).BaseRunner.RunnerRoleType = UnknownDutyType
		}
	default:
		return nil, errors.New("unknown role type")
	}
	return runner, err
}

// //////////////////////////////// For SSV Tests ////////////////////////////////////////////////////////////////
var SSVDecidingMsgsForHeight = func(consensusData *types.ValidatorConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
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

var SSVExpectedDecidingMsgsForHeight = func(consensusData *types.ValidatorConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
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
var DecidingMsgsForHeight = func(consensusData *types.ValidatorConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return DecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var ExpectedDecidingMsgsForHeight = func(consensusData *types.ValidatorConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
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
