package maxmsgsize

import "github.com/ssvlabs/ssv-spec/types"

const (
	expectedSizePrepareSignedSSVMessage                  = 484
	expectedSizeCommitSignedSSVMessage                   = 484
	expectedSizeDecidedSignedSSVMessage                  = 1020
	expectedSizeRoundChangeSignedSSVMessage              = 1948
	expectedSizeProposalSignedSSVMessage                 = 7916
	expectedSizePartialSignatureMessagesSignedSSVMessage = 628
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
		"Test the expected size of a SignedSSVMessage containing a prepare QBFT message",
		expectedSignedSSVMessageFromObject(expectedPrepare(), 1),
		expectedSizePrepareSignedSSVMessage,
		false,
	)
}

func ExpectedCommitSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected commit SignedSSVMessage",
		"Test the expected size of a SignedSSVMessage containing a commit QBFT message",
		expectedSignedSSVMessageFromObject(expectedCommit(), 1),
		expectedSizeCommitSignedSSVMessage,
		false,
	)
}

func ExpectedDecidedSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected decided SignedSSVMessage",
		"Test the expected size of a SignedSSVMessage containing a decided QBFT message with 3 signers",
		expectedSignedSSVMessageFromObject(expectedCommit(), 3),
		expectedSizeDecidedSignedSSVMessage,
		false,
	)
}

func ExpectedRoundChangeSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected round change SignedSSVMessage",
		"Test the expected size of a SignedSSVMessage containing a round change QBFT message with 3 justifications",
		expectedSignedSSVMessageFromObject(expectedRoundChange(3), 1),
		expectedSizeRoundChangeSignedSSVMessage,
		false,
	)
}

func ExpectedProposalSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected proposal SignedSSVMessage",
		"Test the expected size of a SignedSSVMessage containing a proposal QBFT message with 3 justifications and full data",
		expectedSignedSSVMessageWithFullDataFromObject(expectedProposal(3), 1),
		expectedSizeProposalSignedSSVMessage,
		false,
	)
}

func ExpectedPartialSignatureSignedSSVMessage() *StructureSizeTest {
	return NewStructureSizeTest(
		"expected partial signature SignedSSVMessage",
		"Test the expected size of a SignedSSVMessage containing partial signature messages with full data",
		expectedSignedSSVMessageWithFullDataFromObject(expectedPartialSignatureMessages(1), 1),
		expectedSizePartialSignatureMessagesSignedSSVMessage,
		false,
	)
}
