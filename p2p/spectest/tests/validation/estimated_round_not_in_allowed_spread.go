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

// RoundNotAllowedInTimeSpread tests a consensus message with a round value that is not allowed in the received time
func RoundNotAllowedInTimeSpread() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Round:      3, // Round above estimated time + 1
				Height:     qbft.FirstHeight,
				Root:       testingutils.TestingQBFTRootData,
				Identifier: testingutils.DefaultMsgID[:],
			}),
		},
	}

	// Duty's starting time
	receivedAt := time.Unix(testingutils.NewTestingBeaconNode().GetBeaconNetwork().EstimatedTimeAtSlot(phase0.Slot(qbft.FirstHeight)), 0)

	return &MessageValidationTest{
		Name:          "round not allowed in time spread",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrRoundNotAllowedInTimeSpread.Error(),
		ReceivedAt:    receivedAt,
	}
}
