package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	ExpectedSizePrepareSignedSSVMessage                  = 484
	ExpectedSizeCommitSignedSSVMessage                   = 484
	ExpectedSizeDecidedSignedSSVMessage                  = 1020
	ExpectedSizeRoundChangeSignedSSVMessage              = 1948
	ExpectedSizeProposalSignedSSVMessage                 = 7916
	ExpectedSizePartialSignatureMessagesSignedSSVMessage = 628
)

func expectedFullData() []byte {
	bv := maxBeaconVote()
	bvBytes, err := bv.Encode()
	if err != nil {
		panic(err)
	}
	return bvBytes
}

func expectedSignedSSVMessageFromObject(obj types.Encoder, numSigners int) *types.SignedSSVMessage {

	objBytes, err := obj.Encode()
	if err != nil {
		panic(err)
	}

	signatures := [][]byte{}
	signers := []types.OperatorID{}
	for i := 0; i < numSigners; i++ {
		sig := [256]byte{1}
		signatures = append(signatures, sig[:])

		signers = append(signers, 1)
	}

	return &types.SignedSSVMessage{
		Signatures:  signatures,
		OperatorIDs: signers[:],
		SSVMessage: &types.SSVMessage{
			MsgType: types.SSVConsensusMsgType,
			MsgID:   [56]byte{1},
			Data:    objBytes,
		},
	}
}

func expectedSignedSSVMessageWithFullDataFromObject(obj types.Encoder, numSigners int) *types.SignedSSVMessage {
	msg := expectedSignedSSVMessageFromObject(obj, numSigners)
	msg.FullData = expectedFullData()
	return msg
}

func ExpectedPrepareSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected prepare SignedSSVMessage",
		testdoc.StructureSizeTestExpectedPrepareSignedSSVMessageDoc,
		expectedSignedSSVMessageFromObject(expectedPrepare(), 1),
		ExpectedSizePrepareSignedSSVMessage,
		false,
	)
}

func ExpectedCommitSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected commit SignedSSVMessage",
		testdoc.StructureSizeTestExpectedCommitSignedSSVMessageDoc,
		expectedSignedSSVMessageFromObject(expectedCommit(), 1),
		ExpectedSizeCommitSignedSSVMessage,
		false,
	)
}

func ExpectedDecidedSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected decided SignedSSVMessage",
		testdoc.StructureSizeTestExpectedDecidedSignedSSVMessageDoc,
		expectedSignedSSVMessageFromObject(expectedCommit(), 3),
		ExpectedSizeDecidedSignedSSVMessage,
		false,
	)
}

func ExpectedRoundChangeSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected round change SignedSSVMessage",
		testdoc.StructureSizeTestExpectedRoundChangeSignedSSVMessageDoc,
		expectedSignedSSVMessageFromObject(expectedRoundChange(3), 1),
		ExpectedSizeRoundChangeSignedSSVMessage,
		false,
	)
}

func ExpectedProposalSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected proposal SignedSSVMessage",
		testdoc.StructureSizeTestExpectedProposalSignedSSVMessageDoc,
		expectedSignedSSVMessageWithFullDataFromObject(expectedProposal(3), 1),
		ExpectedSizeProposalSignedSSVMessage,
		false,
	)
}

func ExpectedPartialSignatureSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected partial signature SignedSSVMessage",
		testdoc.StructureSizeTestExpectedPartialSignatureSignedSSVMessageDoc,
		expectedSignedSSVMessageWithFullDataFromObject(expectedPartialSignatureMessages(1), 1),
		ExpectedSizePartialSignatureMessagesSignedSSVMessage,
		false,
	)
}
