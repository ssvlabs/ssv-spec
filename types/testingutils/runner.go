package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/ssv"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/gloas"
)

var TestingHighestDecidedSlot = phase0.Slot(0)

// committeeVoteValueCheckF routes a committee consensus value to the fork-appropriate value check:
// GloasBeaconVoteValueCheckF at Gloas slots (SIP #94 §2), BeaconVoteValueCheckF before. The fork must
// be decided by the duty's slot, not by the value's shape — otherwise a pre-Gloas BeaconVote proposed
// at a Gloas slot would pass the check and then fail to decode in ProcessConsensus, killing the duty
// after decision (and contradicting the GloasPreGloasVote value-check vector).
//
// The slot is resolved lazily because the runner, and the QBFT config holding this check, are built
// before the duty is known: production constructs the checker per duty, so the harness reads the
// running duty at call time.
func committeeVoteValueCheckF(
	signer types.BeaconSigner,
	slotF func() phase0.Slot,
	sharePublicKeys []types.ShareValidatorPK,
	expectedSource phase0.Epoch,
	expectedTarget phase0.Epoch,
) qbft.ProposedValueCheckF {
	return func(data []byte) error {
		slot := slotF()
		if VersionBySlot(slot) >= gloas.DataVersionGloas {
			return ssv.GloasBeaconVoteValueCheckF(signer, slot, sharePublicKeys, expectedSource, expectedTarget)(data)
		}
		return ssv.BeaconVoteValueCheckF(signer, slot, sharePublicKeys, expectedSource, expectedTarget)(data)
	}
}

// committeeDutySlotF resolves the running duty's slot for the committee value check.
//
// A missing runner means the construction site never back-patched the reference, which would silently
// validate every duty — Gloas ones included — against TestingDutySlot's pre-Gloas fork and surface as
// a bogus DecodeBeaconVoteErrorCode. That is a wiring bug, so fail loudly rather than fall back. Not
// having started a duty yet is legitimate (e.g. DontStartDuty tests) and keeps the previous
// TestingDutySlot behaviour.
func committeeDutySlotF(runner *ssv.Runner) func() phase0.Slot {
	return func() phase0.Slot {
		if runner == nil || *runner == nil {
			panic("committee value check has no runner: the value check was built without back-patching the runner reference")
		}
		state := (*runner).GetBaseRunner().State
		if state == nil || state.StartingDuty == nil {
			return TestingDutySlot
		}
		return state.StartingDuty.DutySlot()
	}
}

var CommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, keySet)
}

var CommitteeRunnerWithShareMap = func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
	return baseRunnerWithShareMap(types.RoleCommittee, shareMap)
}

var AggregatorCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleAggregatorCommittee, keySet)
}

var AggregatorCommitteeRunnerWithShareMap = func(shareMap map[phase0.ValidatorIndex]*types.Share) ssv.Runner {
	return baseRunnerWithShareMap(types.RoleAggregatorCommittee, shareMap)
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

var SyncCommitteeRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleCommittee, keySet)
}

var ValidatorRegistrationRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleValidatorRegistration, keySet)
}

var VoluntaryExitRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RoleVoluntaryExit, keySet)
}

var PTCAttesterRunner = func(keySet *TestKeySet) ssv.Runner {
	return baseRunner(types.RolePTCAttester, keySet)
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
	return ConstructBaseRunnerWithShareMapAndBeaconNode(role, shareMap, NewTestingBeaconNode())
}

var ConstructBaseRunnerWithShareMapAndBeaconNode = func(role types.RunnerRole, shareMap map[phase0.ValidatorIndex]*types.Share, beacon ssv.BeaconNode) (ssv.Runner, error) {

	var identifier types.MessageID
	var net *TestingNetwork
	var opSigner *types.OperatorSigner
	var valCheck qbft.ProposedValueCheckF
	var contr *qbft.Controller
	// Assigned once the runner exists; the committee value check reads the running duty through it.
	var valCheckRunner ssv.Runner

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
		if role == types.RoleCommittee || role == types.RoleAggregatorCommittee {
			committee := make([]uint64, 0)
			for _, op := range keySetInstance.Committee() {
				committee = append(committee, op.Signer)
			}
			committeeID := types.GetCommitteeID(committee)
			identifier = types.NewCommitteeMsgID(TestingSSVDomainType, committeeID, role)
		} else {
			identifier = types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), role)
		}

		// Network
		net = NewTestingNetwork(1, keySetInstance.OperatorKeys[1])

		// Create CommitteeMember
		committeeMember := TestingCommitteeMember(keySetInstance)

		// Create OperatorSigner
		opSigner = NewOperatorSigner(keySetInstance, committeeMember.OperatorID)

		// Create ValueCheck
		switch role {
		case types.RoleCommittee:
			valCheck = committeeVoteValueCheckF(km, committeeDutySlotF(&valCheckRunner),
				sharePubKeys, TestBeaconVote.Source.Epoch, TestBeaconVote.Target.Epoch)
		case types.RoleProposer:
			valCheck = ssv.ProposerValueCheckF(km, types.BeaconTestNetwork,
				(types.ValidatorPK)(shareInstance.ValidatorPubKey), shareInstance.ValidatorIndex, shareInstance.SharePubKey, VersionByEpoch)
		case types.RoleAggregatorCommittee:
			valCheck = ssv.AggregatorCommitteeValueCheckF(km, types.BeaconTestNetwork)
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
			beacon,
			net,
			km,
			opSigner,
			valCheck,
		)
	case types.RoleProposer:
		runner, err = ssv.NewProposerRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			beacon,
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
			beacon,
			net,
			km,
			opSigner,
			types.DefaultGasLimit,
		)
	case types.RoleVoluntaryExit:
		runner, err = ssv.NewVoluntaryExitRunner(
			types.BeaconTestNetwork,
			shareMap,
			beacon,
			net,
			km,
			opSigner,
		)
	case types.RolePTCAttester:
		runner, err = ssv.NewPTCAttesterRunner(
			types.BeaconTestNetwork,
			shareMap,
			beacon,
			net,
			km,
			opSigner,
		)
	case types.RoleAggregatorCommittee:
		runner, err = ssv.NewAggregatorCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			beacon,
			net,
			km,
			opSigner,
			valCheck,
		)
	case UnknownDutyType:
		runner, err = ssv.NewCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			beacon,
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
	valCheckRunner = runner
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
	// Assigned once the runner exists; the committee value check reads the running duty through it.
	var valCheckRunner ssv.Runner

	// Identifier
	var identifier types.MessageID
	if role == types.RoleCommittee || role == types.RoleAggregatorCommittee {
		committee := make([]uint64, 0)
		for _, op := range keySet.Committee() {
			committee = append(committee, op.Signer)
		}
		clusterID := types.GetCommitteeID(committee)
		identifier = types.NewCommitteeMsgID(TestingSSVDomainType, clusterID, role)
	} else {
		identifier = types.NewValidatorMsgID(TestingSSVDomainType, types.ValidatorPK(TestingValidatorPubKey), role)
	}

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
		valCheck = committeeVoteValueCheckF(km, committeeDutySlotF(&valCheckRunner),
			[]types.ShareValidatorPK{share.SharePubKey}, TestBeaconVote.Source.Epoch, TestBeaconVote.Target.Epoch)
	case types.RoleProposer:
		valCheck = ssv.ProposerValueCheckF(km, types.BeaconTestNetwork,
			(types.ValidatorPK)(TestingValidatorPubKey), TestingValidatorIndex, share.SharePubKey, VersionByEpoch)
	case types.RoleAggregatorCommittee:
		valCheck = ssv.AggregatorCommitteeValueCheckF(km, types.BeaconTestNetwork)
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
	case types.RolePTCAttester:
		runner, err = ssv.NewPTCAttesterRunner(
			types.BeaconTestNetwork,
			shareMap,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
		)
	case types.RoleAggregatorCommittee:
		runner, err = ssv.NewAggregatorCommitteeRunner(
			types.BeaconTestNetwork,
			shareMap,
			contr,
			NewTestingBeaconNode(),
			net,
			km,
			opSigner,
			valCheck,
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
	valCheckRunner = runner
	return runner, err
}

// //////////////////////////////// For SSV Tests ////////////////////////////////////////////////////////////////
var SSVDecidingMsgsForHeight = func(consensusData *types.ProposerConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return SSVDecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var SSVExpectedDecidingMsgsForHeight = func(consensusData *types.ProposerConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
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
var DecidingMsgsForHeight = func(consensusData *types.ProposerConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
	byts, _ := consensusData.Encode()
	r, _ := qbft.HashDataRoot(byts)
	fullData, _ := consensusData.MarshalSSZ()

	return DecidingMsgsForHeightWithRoot(r, fullData, msgIdentifier, height, keySet)
}

var ExpectedDecidingMsgsForHeight = func(consensusData *types.ProposerConsensusData, msgIdentifier []byte, height qbft.Height, keySet *TestKeySet) []*types.SignedSSVMessage {
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
