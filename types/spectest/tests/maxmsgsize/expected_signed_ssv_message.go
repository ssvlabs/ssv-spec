package maxmsgsize

import (
	"github.com/ssvlabs/ssv-spec/types"
	"github.com/ssvlabs/ssv-spec/types/spectest/testdoc"
)

const (
	ExpectedSizePrepareSignedSSVMessage     = 484
	ExpectedSizeCommitSignedSSVMessage      = 484
	ExpectedSizeDecidedSignedSSVMessage     = 1020
	ExpectedSizeRoundChangeSignedSSVMessage = 1948
)

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
