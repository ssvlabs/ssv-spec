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
	return &StructureSizeTest{
		Name:                  "expected prepare SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedPrepare(), 1),
		ExpectedEncodedLength: expectedSizePrepareSignedSSVMessage,
		IsMaxSize:             false,
	}
}

func ExpectedCommitSignedSSVMessage() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "expected commit SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedCommit(), 1),
		ExpectedEncodedLength: expectedSizeCommitSignedSSVMessage,
		IsMaxSize:             false,
	}
}

func ExpectedDecidedSignedSSVMessage() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "expected decided SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedCommit(), 3),
		ExpectedEncodedLength: expectedSizeDecidedSignedSSVMessage,
		IsMaxSize:             false,
	}
}

func ExpectedRoundChangeSignedSSVMessage() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "expected round change SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedRoundChange(3), 1),
		ExpectedEncodedLength: expectedSizeRoundChangeSignedSSVMessage,
		IsMaxSize:             false,
	}
}

func ExpectedProposalSignedSSVMessage() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "expected proposal SignedSSVMessage",
		Object:                expectedSignedSSVMessageWithFullDataFromObject(expectedProposal(3), 1),
		ExpectedEncodedLength: expectedSizeProposalSignedSSVMessage,
		IsMaxSize:             false,
	}
}

func ExpectedPartialSignatureSignedSSVMessage() *StructureSizeTest {
	return &StructureSizeTest{
		Name:                  "expected partial signature SignedSSVMessage",
		Object:                expectedSignedSSVMessageWithFullDataFromObject(expectedPartialSignatureMessages(1), 1),
		ExpectedEncodedLength: expectedSizePartialSignatureMessagesSignedSSVMessage,
		IsMaxSize:             false,
	}
}
