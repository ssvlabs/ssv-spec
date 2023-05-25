package testingutils

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

var SSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.BeaconRole) []*types.SSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
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
	for _, m := range PreConsensusQuorumV(ks, role, consensusData.Version) {
		base = append(base, ssvMsgF(nil, m))
	}

	// consensus and post consensus
	qbftMsgs := DecidingMsgsForHeight(consensusData, id[:], qbft.FirstHeight, ks)
	for _, msg := range qbftMsgs {
		base = append(base, ssvMsgF(msg, nil))
	}
	return base
}

var ExpectedSSVDecidingMsgsV = func(consensusData *types.ConsensusData, ks *TestKeySet, role types.BeaconRole) []*types.SSVMessage {
	id := types.NewMsgID(TestingSSVDomainType, TestingValidatorPubKey[:], role)

	ssvMsgF := func(qbftMsg *qbft.SignedMessage, partialSigMsg *types.SignedPartialSignatureMessage) *types.SSVMessage {
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
		for i := uint64(1); i <= ks.Threshold; i++ {
			base = append(base, ssvMsgF(nil, PreConsensusRandaoMsgV(ks.Shares[types.OperatorID(i)], types.OperatorID(i), consensusData.Version)))
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

	// consensus and post consensus
	qbftMsgs := ExpectedDecidingMsgsForHeight(consensusData, id[:], qbft.FirstHeight, ks)
	for _, msg := range qbftMsgs {
		base = append(base, ssvMsgF(msg, nil))
	}
	return base
}

var PreConsensusQuorumV = func(ks *TestKeySet, role types.BeaconRole, version spec.DataVersion) []*types.SignedPartialSignatureMessage {
	ret := make([]*types.SignedPartialSignatureMessage, 0)
	if role == types.BNRoleProposer {
		for i := uint64(1); i <= ks.Threshold; i++ {
			ret = append(ret, PreConsensusRandaoMsgV(ks.Shares[i], i, version))
		}
	} else if role == types.BNRoleAggregator {
		for i := uint64(1); i <= ks.Threshold; i++ {
			ret = append(ret, PreConsensusSelectionProofMsg(ks.Shares[i], ks.Shares[i], i, i))
		}
	} else if role == types.BNRoleSyncCommitteeContribution {
		for i := uint64(1); i <= ks.Threshold; i++ {
			ret = append(ret, PreConsensusContributionProofMsg(ks.Shares[i], ks.Shares[i], i, i))
		}
	}
	return ret
}
