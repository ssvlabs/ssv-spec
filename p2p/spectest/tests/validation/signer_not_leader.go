package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// SignerNotLeader tests a proposal message with a wrong proposer (not leader)
func SignerNotLeader() tests.SpecTest {

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{2}, // By TestingMessageValidator().Config, the default leader is 1. So, 2 will trigger "signer not leader"
		Signatures:  [][]byte{testingutils.MockRSASignature[:]},
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   testingutils.DefaultMsgID,
			Data: testingutils.EncodeQbftMessage(&qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Round:      qbft.FirstRound,
				Root:       testingutils.TestingQBFTRootData,
				Identifier: testingutils.DefaultMsgID[:],
			}),
		},
		FullData: testingutils.TestingQBFTFullData,
	}

	return &MessageValidationTest{
		Name:          "signer not leader",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrSignerNotLeader.Error(),
	}
}
