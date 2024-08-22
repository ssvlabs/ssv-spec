package validation

import (
	"github.com/ssvlabs/ssv-spec/p2p/spectest/tests"
	"github.com/ssvlabs/ssv-spec/p2p/validation"
	"github.com/ssvlabs/ssv-spec/qbft"
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/testingutils"
)

// WrongSignatureConsensusMessage tests a consensus message with wrong signature
func WrongSignatureConsensusMessage() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVConsensusMsgType,
		MsgID:   msgID,
		Data: testingutils.EncodeQbftMessage(&qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Round:      qbft.FirstRound,
			Identifier: msgID[:],
			Root:       [32]byte{1},
			Height:     qbft.FirstHeight,
		}),
	}

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.SignSSVMessage(ks.OperatorKeys[2], ssvMsg)}, // Wrong key
		SSVMessage:  ssvMsg,
	}

	return &MessageValidationTest{
		Name:          "wrong signature consensus message",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrSignatureVerification.Error(),
	}
}

// WrongSignaturePartialSignatureMessage tests a partial sig message with wrong signature
func WrongSignaturePartialSignatureMessage() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()
	msgID := types.NewMsgID(testingutils.TestingSSVDomainType, testingutils.TestingValidatorPK[:], types.RoleProposer)

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   msgID,
		Data: testingutils.EncodePartialSignatureMessage(&types.PartialSignatureMessages{
			Type: types.PostConsensusPartialSig,
			Slot: 0,
			Messages: []*types.PartialSignatureMessage{
				{
					PartialSignature: testingutils.MockPartialSignature[:],
					SigningRoot:      [32]byte{1},
					Signer:           1,
					ValidatorIndex:   1,
				},
			},
		}),
	}

	msg := &types.SignedSSVMessage{
		OperatorIDs: []types.OperatorID{1},
		Signatures:  [][]byte{testingutils.SignSSVMessage(ks.OperatorKeys[2], ssvMsg)}, // Wrong key
		SSVMessage:  ssvMsg,
	}

	return &MessageValidationTest{
		Name:          "wrong signature partial signature",
		Messages:      [][]byte{testingutils.EncodeMessage(msg)},
		ExpectedError: validation.ErrSignatureVerification.Error(),
	}
}
