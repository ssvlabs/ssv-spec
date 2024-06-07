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

func ExpectedPrepareSignedSSVMessage() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "expected prepare SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedPrepare(), 1),
		ExpectedEncodedLength: expectedSizePrepareSignedSSVMessage,
		IsMaxSizeForType:      false,
	}
}

func ExpectedCommitSignedSSVMessage() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "expected commit SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedCommit(), 1),
		ExpectedEncodedLength: expectedSizeCommitSignedSSVMessage,
		IsMaxSizeForType:      false,
	}
}

func ExpectedDecidedSignedSSVMessage() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "expected decided SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedCommit(), 3),
		ExpectedEncodedLength: expectedSizeDecidedSignedSSVMessage,
		IsMaxSizeForType:      false,
	}
}

func ExpectedRoundChangeSignedSSVMessage() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "expected round change SignedSSVMessage",
		Object:                expectedSignedSSVMessageFromObject(expectedRoundChange(3), 1),
		ExpectedEncodedLength: expectedSizeRoundChangeSignedSSVMessage,
		IsMaxSizeForType:      false,
	}
}

func ExpectedProposalSignedSSVMessage() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "expected proposal SignedSSVMessage",
		Object:                expectedSignedSSVMessageWithFullDataFromObject(expectedProposal(3), 1),
		ExpectedEncodedLength: expectedSizeProposalSignedSSVMessage,
		IsMaxSizeForType:      false,
	}
}

func ExpectedPartialSignatureSignedSSVMessage() *MaxMessageTest {
	return &MaxMessageTest{
		Name:                  "expected partial signature SignedSSVMessage",
		Object:                expectedSignedSSVMessageWithFullDataFromObject(expectedPartialSignatureMessages(1), 1),
		ExpectedEncodedLength: expectedSizePartialSignatureMessagesSignedSSVMessage,
		IsMaxSizeForType:      false,
	}
}
