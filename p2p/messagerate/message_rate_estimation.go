package messagerate

import (
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/types"
)

// Ethereum parameters
const (
	EthereumValidators                       = 1000000.0 // May be taken from the network
	SyncCommitteeSize                        = 512.0
	EstimatedAttestationCommitteeSize        = EthereumValidators / 2048.0
	AggregatorProbability                    = 16.0 / EstimatedAttestationCommitteeSize
	ProposalProbability                      = 1.0 / EthereumValidators
	SyncCommitteeProbability                 = SyncCommitteeSize / EthereumValidators
	SyncCommitteeAggProb                     = SyncCommitteeProbability * 16.0 / (SyncCommitteeSize / 4.0)
	SlotsPerEpoch                            = 32.0
	MaxAttestationDutiesPerEpochForCommittee = SlotsPerEpoch
)

// Expected number of messages per duty step

func consensusMessages(n int) int {
	return 1 + n + n + 2 // 1 Proposal + n Prepares + n Commits + 2 Decideds (average)
}

func partialSignatureMessages(n int) int {
	return n
}

func dutyWithPreConsensus(n int) int {
	// Pre-Consensus + Consensus + Post-Consensus
	return partialSignatureMessages(n) + consensusMessages(n) + partialSignatureMessages(n)
}

func dutyWithoutPreConsensus(n int) int {
	// Consensus + Post-Consensus
	return consensusMessages(n) + partialSignatureMessages(n)
}

// Expected number of committee duties per epoch due to attestations
func expectedNumberOfCommitteeDutiesPerEpochDueToAttestation(numValidators int) float64 {
	k := float64(numValidators)
	n := SlotsPerEpoch

	// Probability that all validators are not assigned to slot i
	probabilityAllNotOnSlotI := math.Pow((n-1)/n, k)
	// Probability that at least one validator is assigned to slot i
	probabilityAtLeastOneOnSlotI := 1 - probabilityAllNotOnSlotI
	// Expected value for duty existence ({0,1}) on slot i
	expectedDutyExistenceOnSlotI := 0*probabilityAllNotOnSlotI + 1*probabilityAtLeastOneOnSlotI
	// Expected number of duties per epoch
	expectedNumberOfDutiesPerEpoch := n * expectedDutyExistenceOnSlotI

	return expectedNumberOfDutiesPerEpoch
}

// Expected committee duties per epoch that are due to only sync committee beacon duties
func expectedSingleSCCommitteeDutiesPerEpoch(numValidators int) float64 {
	// Probability that a validator is not in sync committee
	chanceOfNotBeingInSyncCommittee := 1.0 - SyncCommitteeProbability
	// Probability that all validators are not in sync committee
	chanceThatAllValidatorsAreNotInSyncCommittee := math.Pow(chanceOfNotBeingInSyncCommittee, float64(numValidators))
	// Probability that at least one validator is in sync committee
	chanceOfAtLeastOneValidatorBeingInSyncCommittee := 1.0 - chanceThatAllValidatorsAreNotInSyncCommittee

	// Expected number of slots with no attestation duty
	expectedSlotsWithNoDuty := 32.0 - expectedNumberOfCommitteeDutiesPerEpochDueToAttestation(numValidators)

	// Expected number of committee duties per epoch created due to only sync committee duties
	return chanceOfAtLeastOneValidatorBeingInSyncCommittee * expectedSlotsWithNoDuty
}

// Committee represents a Committee entity with a certain number of operators and validators
type Committee struct {
	Operators  []types.OperatorID
	Validators []phase0.ValidatorIndex
}

// Estimates the message rate for a topic given the topic's committees
func EstimateMessageRateForTopic(committees []*Committee) float64 {
	if len(committees) == 0 {
		return 0
	}

	totalMsgRate := 0.0

	for _, committee := range committees {
		committeeSize := len(committee.Operators)
		numValidators := len(committee.Validators)

		totalMsgRate += expectedNumberOfCommitteeDutiesPerEpochDueToAttestation(numValidators) * float64(dutyWithoutPreConsensus(committeeSize))
		totalMsgRate += expectedSingleSCCommitteeDutiesPerEpoch(numValidators) * float64(dutyWithoutPreConsensus(committeeSize))
		totalMsgRate += float64(numValidators) * AggregatorProbability * float64(dutyWithPreConsensus(committeeSize))
		totalMsgRate += float64(numValidators) * SlotsPerEpoch * ProposalProbability * float64(dutyWithPreConsensus(committeeSize))
		totalMsgRate += float64(numValidators) * SlotsPerEpoch * SyncCommitteeAggProb * float64(dutyWithPreConsensus(committeeSize))
	}

	// Convert rate to seconds
	totalEpochSeconds := float64(SlotsPerEpoch * 12)
	totalMsgRate = totalMsgRate / totalEpochSeconds

	return totalMsgRate
}
