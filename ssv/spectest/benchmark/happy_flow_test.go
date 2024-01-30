package benchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

func FullDutyXRounds(numRounds int) {
	ks := testingutils.Testing4SharesSet()
	validator := testingutils.BaseValidator(ks)
	duty := &testingutils.TestingAggregatorDuty
	role := duty.Type
	cd := testingutils.TestAggregatorConsensusDataByts
	err := validator.StartDuty(duty)
	height := qbft.Height(duty.Slot)
	msgID := testingutils.AggregatorMsgID
	if err != nil {
		panic(err.Error())
	}

	msgs := make([]*types.SSVMessage, 0)

	// Pre-consensus
	msgs = append(msgs, PreConsensusF(ks, role, false)...)

	// Consensus
	for i := 1; i <= numRounds; i++ {
		usePreparedValue := (i > 1)
		if i == numRounds {
			msgs = append(msgs, ConsensusForRound(ks, role, height, qbft.Round(i), msgID, cd, usePreparedValue, int(ks.ShareCount), int(ks.ShareCount))...)
		} else {
			msgs = append(msgs, ConsensusForRound(ks, role, height, qbft.Round(i), msgID, cd, usePreparedValue, int(ks.ShareCount), int(ks.ShareCount)-1)...)
		}
	}

	// Post-consensus
	msgs = append(msgs, PostConsensusF(ks, role, true, height, spec.DataVersionDeneb)...)

	// Process
	start := time.Now()

	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	end := time.Now()
	elapsed := end.Sub(start).Microseconds()
	fmt.Printf("Total duty (%v rounds): %v us.\n", numRounds, elapsed)
}

func ConsensusXRound(round int) {
	ks := testingutils.Testing4SharesSet()
	validator := testingutils.BaseValidator(ks)
	duty := &testingutils.TestingAggregatorDuty
	role := duty.Type
	cd := testingutils.TestAggregatorConsensusDataByts
	err := validator.StartDuty(duty)
	height := qbft.Height(duty.Slot)
	msgID := testingutils.AggregatorMsgID
	if err != nil {
		panic(err.Error())
	}

	msgs := make([]*types.SSVMessage, 0)

	// Pre-consensus
	msgs = append(msgs, PreConsensusF(ks, role, false)...)

	// Consensus
	for i := 1; i <= round-1; i++ {
		usePreparedValue := (i > 1)
		msgs = append(msgs, ConsensusForRound(ks, role, height, qbft.Round(i), msgID, cd, usePreparedValue, int(ks.ShareCount), int(ks.ShareCount)-1)...)
	}

	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	msgs = make([]*types.SSVMessage, 0)

	usePreparedValue := (round > 1)
	msgs = append(msgs, ConsensusForRound(ks, role, height, qbft.Round(round), msgID, cd, usePreparedValue, int(ks.ShareCount), int(ks.ShareCount))...)

	// Process
	start := time.Now()

	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	end := time.Now()
	elapsed := end.Sub(start).Microseconds()
	fmt.Printf("Consensus round %v time: %v us.\n", round, elapsed)
}

func TestFullDutyRound(t *testing.T) {
	FullDutyXRounds(1)
	FullDutyXRounds(2)
	FullDutyXRounds(3)
	FullDutyXRounds(4)
	FullDutyXRounds(5)
	FullDutyXRounds(6)
}

func TestConsensusRounds(t *testing.T) {
	ConsensusXRound(1)
	ConsensusXRound(2)
	ConsensusXRound(3)
}

func TestPreConsensus(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	validator := testingutils.BaseValidator(ks)
	duty := &testingutils.TestingAggregatorDuty
	role := duty.Type
	err := validator.StartDuty(duty)
	if err != nil {
		panic(err.Error())
	}

	msgs := make([]*types.SSVMessage, 0)

	// Pre-consensus
	msgs = append(msgs, PreConsensusF(ks, role, false)...)

	// Process
	start := time.Now()

	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	end := time.Now()
	elapsed := end.Sub(start).Microseconds()
	fmt.Printf("Pre-consensus full committee: %v us.\n", elapsed)
}

func TestPostConsensus(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	validator := testingutils.BaseValidator(ks)
	duty := &testingutils.TestingAggregatorDuty
	role := duty.Type
	cd := testingutils.TestAggregatorConsensusDataByts
	err := validator.StartDuty(duty)
	height := qbft.Height(duty.Slot)
	msgID := testingutils.AggregatorMsgID
	if err != nil {
		panic(err.Error())
	}

	msgs := make([]*types.SSVMessage, 0)

	// Pre-consensus
	msgs = append(msgs, PreConsensusF(ks, role, false)...)

	// Consensus
	msgs = append(msgs, ConsensusForRound(ks, role, height, 1, msgID, cd, false, int(ks.ShareCount), int(ks.ShareCount))...)

	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	msgs = make([]*types.SSVMessage, 0)
	// Post-consensus
	msgs = append(msgs, PostConsensusF(ks, role, true, height, spec.DataVersionDeneb)...)

	// Process
	start := time.Now()

	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	end := time.Now()
	elapsed := end.Sub(start).Microseconds()
	fmt.Printf("Post-consensus: %v us.\n", elapsed)
}
