package benchmark

import (
	"fmt"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
	"github.com/herumi/bls-eth-go-binary/bls"
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
	times := make([]int64, 0)
	var total int64

	for _, msg := range msgs {
		start := time.Now()
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
		end := time.Now()
		elapsed := end.Sub(start)
		total += elapsed.Microseconds()
		times = append(times, elapsed.Milliseconds())
	}

	fmt.Printf("Consensus round %v time: %v ms. Total: %v us\n", round, times, total)
}

func SinglePartialSigMessage() {
	ks := testingutils.Testing4SharesSet()
	validator := testingutils.BaseValidator(ks)
	duty := &testingutils.TestingAggregatorDuty
	role := duty.Type
	// cd := testingutils.TestAggregatorConsensusDataByts
	err := validator.StartDuty(duty)
	// height := qbft.Height(duty.Slot)
	// msgID := testingutils.AggregatorMsgID
	if err != nil {
		panic(err.Error())
	}

	msgs := make([]*types.SSVMessage, 0)

	// Pre-consensus
	msgs = append(msgs, PreConsensusF(ks, role, false)...)

	start := time.Now()
	err = validator.ProcessMessage(msgs[0])
	if err != nil {
		panic(err.Error())
	}
	end := time.Now()
	elapsed := end.Sub(start).Microseconds()
	fmt.Printf("Single partial message time: %v us.\n", elapsed)
}

func ConsensusMessageXRound(round int, msgType qbft.MessageType) {
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

	// Process pre-consensus and old-consensus
	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	msgs = make([]*types.SSVMessage, 0)

	usePreparedValue := (round > 1)
	msgs = append(msgs, ConsensusForRound(ks, role, height, qbft.Round(round), msgID, cd, usePreparedValue, int(ks.ShareCount), int(ks.ShareCount))...)

	for _, msg := range msgs {
		signedMessage := &qbft.SignedMessage{}
		err := signedMessage.Decode(msg.Data)
		if err != nil {
			panic(err.Error())
		}
		mType := signedMessage.Message.MsgType
		if mType == msgType {
			start := time.Now()
			err := validator.ProcessMessage(msg)
			if err != nil {
				panic(err.Error())
			}
			end := time.Now()
			elapsed := end.Sub(start).Microseconds()
			fmt.Printf("Consensus round %v msg type %v time: %v us.\n", round, msgType, elapsed)
			return
		} else {
			err := validator.ProcessMessage(msg)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func YConsensusMessageXRound(round int, messageNumber int, msgType qbft.MessageType) {
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

	// Process pre-consensus and old-consensus
	for _, msg := range msgs {
		err := validator.ProcessMessage(msg)
		if err != nil {
			panic(err.Error())
		}
	}

	msgs = make([]*types.SSVMessage, 0)

	usePreparedValue := (round > 1)
	msgs = append(msgs, ConsensusForRound(ks, role, height, qbft.Round(round), msgID, cd, usePreparedValue, int(ks.ShareCount), int(ks.ShareCount))...)

	messageOccurrenceNumber := 0
	for _, msg := range msgs {
		signedMessage := &qbft.SignedMessage{}
		err := signedMessage.Decode(msg.Data)
		if err != nil {
			panic(err.Error())
		}
		mType := signedMessage.Message.MsgType

		foundMessage := false
		if mType == msgType {
			messageOccurrenceNumber += 1
			if messageOccurrenceNumber == messageNumber {
				foundMessage = true
			}
		}

		if foundMessage {
			start := time.Now()
			err := validator.ProcessMessage(msg)
			if err != nil {
				panic(err.Error())
			}
			end := time.Now()
			elapsed := end.Sub(start).Microseconds()
			fmt.Printf("Consensus round %v msg type %v occurrence %v time: %v us.\n", round, msgType, messageNumber, elapsed)
			return
		} else {
			err := validator.ProcessMessage(msg)
			if err != nil {
				panic(err.Error())
			}
		}
	}
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

func TestConsensusMessage(t *testing.T) {
	ConsensusMessageXRound(1, qbft.ProposalMsgType)
	ConsensusMessageXRound(1, qbft.PrepareMsgType)
	ConsensusMessageXRound(1, qbft.CommitMsgType)
	ConsensusMessageXRound(2, qbft.RoundChangeMsgType)
	ConsensusMessageXRound(2, qbft.ProposalMsgType)
	ConsensusMessageXRound(2, qbft.PrepareMsgType)
	ConsensusMessageXRound(2, qbft.CommitMsgType)
	ConsensusMessageXRound(3, qbft.RoundChangeMsgType)
	ConsensusMessageXRound(3, qbft.ProposalMsgType)
	ConsensusMessageXRound(3, qbft.PrepareMsgType)
	ConsensusMessageXRound(3, qbft.CommitMsgType)
}

func TestRoundChangeMessage(t *testing.T) {
	YConsensusMessageXRound(2, 1, qbft.RoundChangeMsgType)
	YConsensusMessageXRound(2, 2, qbft.RoundChangeMsgType)
	YConsensusMessageXRound(2, 3, qbft.RoundChangeMsgType)
	YConsensusMessageXRound(3, 1, qbft.RoundChangeMsgType)
	YConsensusMessageXRound(3, 2, qbft.RoundChangeMsgType)
	YConsensusMessageXRound(3, 3, qbft.RoundChangeMsgType)
}

func TestSinglePartialMessage(t *testing.T) {
	SinglePartialSigMessage()
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

func TestCryptographyPrimitives(t *testing.T) {

	// Init
	ks := testingutils.Testing4SharesSet()
	signer := testingutils.NewTestingKeyManager()
	beacon := testingutils.NewTestingBeaconNode()
	epoch := testingutils.TestingDutyEpoch
	d, _ := beacon.DomainData(phase0.Epoch(epoch), types.DomainRandao)

	// Result data
	signingTimes := make([]float64, 0)
	verifyingTimes := make([]float64, 0)
	PartialSigContainer := ssv.NewPartialSigContainer(ks.Threshold)
	partialSignatureMessages := make([]*types.PartialSignatureMessage, 0)

	var signingRoot [32]byte
	opID := uint64(1)
	// Create partial signature messages, storing the signing and verification times
	for opID <= ks.Threshold {

		sk := ks.Shares[opID]

		// Sign
		start := time.Now()
		signature, root, _ := signer.SignBeaconObject(types.SSZUint64(epoch), d, sk.GetPublicKey().Serialize(), types.DomainRandao)
		end := time.Now()
		signingTimes = append(signingTimes, float64(end.Sub(start).Microseconds()))
		signingRoot = root

		// Create and store message
		msg := &types.PartialSignatureMessage{
			PartialSignature: signature[:],
			SigningRoot:      root,
			Signer:           opID,
		}
		partialSignatureMessages = append(partialSignatureMessages, msg)
		PartialSigContainer.AddSignature(msg)

		// Verify message
		pk := &bls.PublicKey{}
		committee := ks.Committee()
		for idx, op := range committee {
			if op.OperatorID == opID {
				start := time.Now()
				if err := pk.Deserialize(committee[idx].GetPublicKey()); err != nil {
					panic("Test failed to deserialize public key")
				}
				sig := &bls.Sign{}
				if err := sig.Deserialize(signature); err != nil {
					panic("Test failed to deserialize signature")
				}
				if !sig.VerifyByte(pk, root[:]) {
					panic("Test failed to verify signature")
				}
				end := time.Now()
				verifyingTimes = append(verifyingTimes, float64(end.Sub(start).Microseconds()))
			}
		}

		opID += 1
	}

	// Print signing time
	meanSign, stddevSign := GetMeanAndStddev(signingTimes)
	fmt.Printf("BLS Signing time: %.2f ± %.2f us.\n", meanSign, stddevSign)

	// Print verification time
	meanVerify, stddevVerify := GetMeanAndStddev(verifyingTimes)
	fmt.Printf("BLS Verification time: %.2f ± %.2f us.\n", meanVerify, stddevVerify)

	// Reconstruct time
	validatorPubKey := ks.ValidatorPK.Serialize()
	start := time.Now()
	_, err := PartialSigContainer.ReconstructSignature(signingRoot, validatorPubKey)
	if err != nil {
		panic(err)
	}
	end := time.Now()
	elapsed := end.Sub(start).Microseconds()
	fmt.Printf("Reconstruct + Verify signature: %v us.\n", elapsed)
	fmt.Printf("Reconstruct (from the above diff): %.2f us.\n", float64(elapsed)-meanVerify)

	// Aggregation time
	var aggregated types.Signature
	start = time.Now()
	for _, msg := range partialSignatureMessages {
		if aggregated == nil {
			aggregated = msg.PartialSignature
		} else {
			aggregated, err = aggregated.Aggregate(msg.PartialSignature)
			if err != nil {
				panic(err)
			}
		}
	}
	end = time.Now()
	elapsed = end.Sub(start).Microseconds()
	fmt.Printf("Aggregate 3 signatures: %v us.\n", elapsed)
}
