package validation

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// ConsensusMessageSentBeforeSlotStarts tests a consensus message sent before the slot starts
func ConsensusMessageSentBeforeSlotStarts() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Round:      qbft.FirstRound,
				Identifier: testingutils.DefaultMsgID[:],
				Root:       [32]byte{1},
				Height:     qbft.FirstHeight,
			}),
		},
	}

	// Received time: slot start time - 2 * clock error tolerance
	receivedAt := time.Unix(testingutils.NewTestingBeaconNode().GetBeaconNetwork().EstimatedTimeAtSlot(phase0.Slot(qbft.FirstHeight)), 0)
	receivedAt = receivedAt.Add(time.Duration(-2 * validation.ClockErrorTolerance))

	return &MessageValidationTest{
		Name:          "consensus message sent before slot starts",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrEarlySlotMessage.Error(),
		ReceivedAt:    receivedAt,
	}
}

// PartialSignatureMessageSentBeforeSlotStarts tests a partial signature message sent before the slot starts
func PartialSignatureMessageSentBeforeSlotStarts() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVPartialSignatureMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodePartialSignatureMessage(&types.PartialSignatureMessages{
				Type: types.PostConsensusPartialSig,
				Slot: 1,
				Messages: []*types.PartialSignatureMessage{
					{
						PartialSignature: testingutils.MockPartialSignature[:],
						SigningRoot:      [32]byte{1},
						Signer:           1,
						ValidatorIndex:   1,
					},
				},
			}),
		},
	}

	// Received time: slot start time - 2 * clock error tolerance
	receivedAt := time.Unix(testingutils.NewTestingBeaconNode().GetBeaconNetwork().EstimatedTimeAtSlot(1), 0)
	receivedAt = receivedAt.Add(time.Duration(-2 * validation.ClockErrorTolerance))

	return &MessageValidationTest{
		Name:          "partial signature message sent before slot starts",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrEarlySlotMessage.Error(),
		ReceivedAt:    receivedAt,
	}
}
