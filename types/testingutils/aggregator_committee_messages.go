package testingutils

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
)

// Return the input for the committee for an aggregator committee duty for a given slot
func AggregatorCommitteeInputForSlot(aggVals []int, scVals []int, ksMap map[phase0.ValidatorIndex]*TestKeySet, slot phase0.Slot, version spec.DataVersion) (*types.AggregatorCommitteeDuty, []*types.SignedSSVMessage) {
	// Duty
	duty := TestingAggregatorCommitteeDutyForSlot(slot, aggVals, scVals)

	// Pre-Consensus, QBFT and Post-Consensus
	msgs := AggregatorCommitteeInputForDuty(duty, ksMap, version)

	return duty, msgs
}

func AggregatorCommitteeInputForDuty(duty *types.AggregatorCommitteeDuty, ksMap map[phase0.ValidatorIndex]*TestKeySet, version spec.DataVersion) []*types.SignedSSVMessage {

	ks := ksMap[duty.ValidatorDuties[0].ValidatorIndex]
	slot := duty.Slot
	height := qbft.Height(slot)

	msgID := AggregatorCommitteeMsgIDForKeySet(ks)

	consensusDataBytes := TestAggregatorCommitteeConsensusDataBytesForDuty(duty, version)
	cdRoot := sha256.Sum256(consensusDataBytes)

	msgs := []*types.SignedSSVMessage{
		// Pre-consensus messages
		SignPartialSigSSVMessage(ks, SSVMsgAggregatorCommittee(ks, nil, PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1))),
		SignPartialSigSSVMessage(ks, SSVMsgAggregatorCommittee(ks, nil, PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 2))),
		SignPartialSigSSVMessage(ks, SSVMsgAggregatorCommittee(ks, nil, PreConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 3))),

		// Consensus messages
		TestingProposalMessageWithIdentifierAndFullData(ks.OperatorKeys[1], types.OperatorID(1), msgID[:], consensusDataBytes, height),
		TestingPrepareMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID[:], cdRoot),
		TestingPrepareMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID[:], cdRoot),
		TestingPrepareMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID[:], cdRoot),
		TestingCommitMessageWithParams(ks.OperatorKeys[1], 1, 1, height, msgID[:], cdRoot),
		TestingCommitMessageWithParams(ks.OperatorKeys[2], 2, 1, height, msgID[:], cdRoot),
		TestingCommitMessageWithParams(ks.OperatorKeys[3], 3, 1, height, msgID[:], cdRoot),

		// Post-consensus messages
		SignPartialSigSSVMessage(ks, SSVMsgAggregatorCommittee(ks, nil, PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 1, version))),
		SignPartialSigSSVMessage(ks, SSVMsgAggregatorCommittee(ks, nil, PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 2, version))),
		SignPartialSigSSVMessage(ks, SSVMsgAggregatorCommittee(ks, nil, PostConsensusAggregatorCommitteeMsgForDuty(duty, ksMap, 3, version))),
	}

	return msgs
}
