package benchmark

import (
	"math"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func PreConsensusF(ks *testingutils.TestKeySet, role types.BeaconRole, stopQuorum bool) []*types.SSVMessage {

	var genFunction func(opID types.OperatorID) *types.SignedPartialSignatureMessage

	switch role {
	case types.BNRoleAggregator:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PreConsensusSelectionProofMsg(ks.Shares[opID], ks.Shares[opID], opID, opID)
		}
	case types.BNRoleProposer:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PreConsensusRandaoMsg(ks.Shares[opID], opID)
		}
	case types.BNRoleSyncCommitteeContribution:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PreConsensusContributionProofMsg(ks.Shares[opID], ks.Shares[opID], opID, opID)
		}
	default:
		return []*types.SSVMessage{}
	}

	numMessages := 0
	msgs := make([]*types.SignedPartialSignatureMessage, 0)
	stopValue := len(ks.Committee())
	if stopQuorum {
		stopValue = int(ks.Threshold)
	}
	for _, op := range ks.Committee() {
		id := op.OperatorID
		msgs = append(msgs, genFunction(id))
		numMessages += 1
		if numMessages >= stopValue {
			break
		}
	}
	return PartialToSSVMessage(msgs, role)
}

func PostConsensusF(ks *testingutils.TestKeySet, role types.BeaconRole, stopQuorum bool, height qbft.Height, version spec.DataVersion) []*types.SSVMessage {

	var genFunction func(opID types.OperatorID) *types.SignedPartialSignatureMessage

	switch role {
	case types.BNRoleAttester:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PostConsensusAttestationMsg(ks.Shares[opID], opID, height)
		}
	case types.BNRoleAggregator:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PostConsensusAggregatorMsg(ks.Shares[opID], opID)
		}
	case types.BNRoleProposer:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PostConsensusProposerMsgV(ks.Shares[opID], opID, version)
		}
	case types.BNRoleSyncCommittee:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[opID], opID)
		}
	case types.BNRoleSyncCommitteeContribution:
		genFunction = func(opID types.OperatorID) *types.SignedPartialSignatureMessage {
			return testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[opID], opID, ks)
		}
	default:
		return []*types.SSVMessage{}
	}

	numMessages := 0
	msgs := make([]*types.SignedPartialSignatureMessage, 0)
	stopValue := len(ks.Committee())
	if stopQuorum {
		stopValue = int(ks.Threshold)
	}
	for _, op := range ks.Committee() {
		id := op.OperatorID
		msgs = append(msgs, genFunction(id))
		numMessages += 1
		if numMessages >= stopValue {
			break
		}
	}
	return PartialToSSVMessage(msgs, role)
}

func PartialToSSVMessage(msgs []*types.SignedPartialSignatureMessage, role types.BeaconRole) []*types.SSVMessage {

	var ssvMsgF func(msg *types.SignedPartialSignatureMessage) *types.SSVMessage
	switch role {
	case types.BNRoleAttester:
		ssvMsgF = func(msg *types.SignedPartialSignatureMessage) *types.SSVMessage {
			return testingutils.SSVMsgAttester(nil, msg)
		}
	case types.BNRoleAggregator:
		ssvMsgF = func(msg *types.SignedPartialSignatureMessage) *types.SSVMessage {
			return testingutils.SSVMsgAggregator(nil, msg)
		}
	case types.BNRoleProposer:
		ssvMsgF = func(msg *types.SignedPartialSignatureMessage) *types.SSVMessage {
			return testingutils.SSVMsgProposer(nil, msg)
		}
	case types.BNRoleSyncCommittee:
		ssvMsgF = func(msg *types.SignedPartialSignatureMessage) *types.SSVMessage {
			return testingutils.SSVMsgSyncCommittee(nil, msg)
		}
	case types.BNRoleSyncCommitteeContribution:
		ssvMsgF = func(msg *types.SignedPartialSignatureMessage) *types.SSVMessage {
			return testingutils.SSVMsgSyncCommitteeContribution(nil, msg)
		}
	default:
		return []*types.SSVMessage{}
	}

	ret := make([]*types.SSVMessage, 0)
	for _, msg := range msgs {
		ret = append(ret, ssvMsgF(msg))
	}
	return ret
}

func QbftToSSVMessage(msgs []*qbft.SignedMessage, role types.BeaconRole) []*types.SSVMessage {

	var ssvMsgF func(msg *qbft.SignedMessage) *types.SSVMessage
	switch role {
	case types.BNRoleAttester:
		ssvMsgF = func(msg *qbft.SignedMessage) *types.SSVMessage {
			return testingutils.SSVMsgAttester(msg, nil)
		}
	case types.BNRoleAggregator:
		ssvMsgF = func(msg *qbft.SignedMessage) *types.SSVMessage {
			return testingutils.SSVMsgAggregator(msg, nil)
		}
	case types.BNRoleProposer:
		ssvMsgF = func(msg *qbft.SignedMessage) *types.SSVMessage {
			return testingutils.SSVMsgProposer(msg, nil)
		}
	case types.BNRoleSyncCommittee:
		ssvMsgF = func(msg *qbft.SignedMessage) *types.SSVMessage {
			return testingutils.SSVMsgSyncCommittee(msg, nil)
		}
	case types.BNRoleSyncCommitteeContribution:
		ssvMsgF = func(msg *qbft.SignedMessage) *types.SSVMessage {
			return testingutils.SSVMsgSyncCommitteeContribution(msg, nil)
		}
	default:
		return []*types.SSVMessage{}
	}

	ret := make([]*types.SSVMessage, 0)
	for _, msg := range msgs {
		ret = append(ret, ssvMsgF(msg))
	}
	return ret
}

func Proposal(ks *testingutils.TestKeySet, opID types.OperatorID, height qbft.Height, round qbft.Round, msgID []byte, fullData []byte, root [32]byte, rcJustification [][]byte, prepareJustification [][]byte) *qbft.SignedMessage {
	s := testingutils.SignQBFTMsg(ks.Shares[opID], opID, &qbft.Message{
		MsgType:                  qbft.ProposalMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               msgID,
		Root:                     root,
		RoundChangeJustification: rcJustification,
		PrepareJustification:     prepareJustification,
	})
	s.FullData = fullData
	return s
}

func Prepare(ks *testingutils.TestKeySet, opID types.OperatorID, height qbft.Height, round qbft.Round, msgID []byte, root [32]byte) *qbft.SignedMessage {
	s := testingutils.SignQBFTMsg(ks.Shares[opID], opID, &qbft.Message{
		MsgType:    qbft.PrepareMsgType,
		Height:     height,
		Round:      round,
		Identifier: msgID,
		Root:       root,
	})
	return s
}

func Commit(ks *testingutils.TestKeySet, opID types.OperatorID, height qbft.Height, round qbft.Round, msgID []byte, root [32]byte) *qbft.SignedMessage {
	s := testingutils.SignQBFTMsg(ks.Shares[opID], opID, &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     height,
		Round:      round,
		Identifier: msgID,
		Root:       root,
	})
	return s
}

func RoundChange(ks *testingutils.TestKeySet, opID types.OperatorID, height qbft.Height, round qbft.Round, msgID []byte, fullData []byte, root [32]byte, dataRound qbft.Round, rcJustification [][]byte) *qbft.SignedMessage {
	s := testingutils.SignQBFTMsg(ks.Shares[opID], opID, &qbft.Message{
		MsgType:                  qbft.RoundChangeMsgType,
		Height:                   height,
		Round:                    round,
		Identifier:               msgID,
		Root:                     root,
		DataRound:                dataRound,
		RoundChangeJustification: rcJustification,
	})
	s.FullData = fullData
	return s
}

func ConsensusForRound(ks *testingutils.TestKeySet, role types.BeaconRole,
	height qbft.Height, round qbft.Round,
	msgID []byte,
	fullData []byte,
	preparedValue bool,
	maxPrepare int,
	maxCommit int) []*types.SSVMessage {

	loopMessages := func(creator func(opID types.OperatorID) *qbft.SignedMessage, maxMessages int) []*qbft.SignedMessage {
		numMsgs := 0
		ret := make([]*qbft.SignedMessage, 0)
		for _, op := range ks.Committee() {
			opID := op.OperatorID
			ret = append(ret, creator(opID))
			numMsgs += 1
			if numMsgs >= maxMessages {
				break
			}
		}
		return ret
	}

	// full data hash root
	root, err := qbft.HashDataRoot(fullData)
	if err != nil {
		panic(err.Error())
	}

	// return variable
	allMsgs := make([]*qbft.SignedMessage, 0)

	// Justification
	preparedMsgsEncoded := make([][]byte, 0)

	// Prepare messages for justification
	if round > 1 && preparedValue {
		preparedMsgs := loopMessages(func(opID types.OperatorID) *qbft.SignedMessage {
			return Prepare(ks, opID, height, round-1, msgID, root)
		}, int(ks.Threshold))
		preparedMsgsEncoded, err = qbft.MarshalJustifications(preparedMsgs)
		if err != nil {
			panic(err.Error())
		}
	}

	// Round-Change messages for justification
	if round > 1 {
		allMsgs = append(allMsgs, loopMessages(func(opID types.OperatorID) *qbft.SignedMessage {
			return RoundChange(ks, opID, height, round, msgID, fullData, root, round-1, preparedMsgsEncoded)
		}, int(ks.Threshold))...)
	}

	rcMsgsEncoded, err := qbft.MarshalJustifications(allMsgs)
	if err != nil {
		panic(err.Error())
	}

	// Proposal
	allMsgs = append(allMsgs, Proposal(ks, 1, height, round, msgID, fullData, root, rcMsgsEncoded, preparedMsgsEncoded))

	// Prepare
	allMsgs = append(allMsgs, loopMessages(func(opID types.OperatorID) *qbft.SignedMessage {
		return Prepare(ks, opID, height, round, msgID, root)
	}, maxPrepare)...)

	// Commit
	allMsgs = append(allMsgs, loopMessages(func(opID types.OperatorID) *qbft.SignedMessage {
		return Commit(ks, opID, height, round, msgID, root)
	}, maxCommit)...)

	return QbftToSSVMessage(allMsgs, role)
}

// Returns mean and population standard deviation
func GetMeanAndStddev(values []float64) (float64, float64) {
	var mean, stddev float64

	for _, value := range values {
		mean += value
	}
	mean = mean / float64(len(values))

	for _, value := range values {
		stddev += math.Pow(mean-float64(value), 2)
	}
	stddev = math.Pow(stddev/float64(len(values)), 1.0/2)

	return mean, stddev
}
